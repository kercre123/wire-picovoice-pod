package wirepod

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"encoding/json"

	"github.com/pkg/errors"
	leopard "github.com/Picovoice/leopard/binding/go"
	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/chipper/pkg/vtt"
	opus "github.com/digital-dream-labs/opus-go/opus"
	"github.com/soundhound/houndify-sdk-go"
)

var hclient houndify.Client
var houndEnable bool = true

var disableLiveTranscriptionKG bool = false

func ParseSpokenResponse(serverResponseJSON string) (string, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(serverResponseJSON), &result)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("failed to decode json")
	}
	if !strings.EqualFold(result["Status"].(string), "OK") {
		return "", errors.New(result["ErrorMessage"].(string))
	}
	if result["NumToReturn"].(float64) < 1 {
		return "", errors.New("no results to return")
	}
	return result["AllResults"].([]interface{})[0].(map[string]interface{})["SpokenResponse"].(string), nil
}

func InitHoundify() {
	if os.Getenv("HOUNDIFY_ENABLED") == "true" {
		if os.Getenv("HOUNDIFY_CLIENT_ID") == "" {
			fmt.Println("Houndify Client ID not provided.")
			houndEnable = false
		}
		if os.Getenv("HOUNDIFY_CLIENT_KEY") == "" {
			fmt.Println("Houndify Client Key not provided.")
			houndEnable = false
		}
		if houndEnable == true {
			hclient = houndify.Client{
				ClientID:  os.Getenv("HOUNDIFY_CLIENT_ID"),
				ClientKey: os.Getenv("HOUNDIFY_CLIENT_KEY"),
			}
			fmt.Println("Houndify initialized!")
		}
	} else {
		houndEnable = false
	}
}

var NoResult string = "NoResultCommand"
var NoResultSpoken string

var botNumKG int = 0

func knowledgeAPI(sessionID string, spokenText string) string {
	if houndEnable == true {
		hReq := houndify.TextRequest{
			Query:             spokenText,
			UserID:            "victor",
			RequestID:         sessionID,
			RequestInfoFields: make(map[string]interface{}),
		}
		serverResponse, err := hclient.TextSearch(hReq)
		if err != nil {
			fmt.Println(err)
		}
		robotWords, err := ParseSpokenResponse(serverResponse)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Houndify Response: " + robotWords)
		return robotWords
	} else {
		fmt.Println("Houndify is not enabled, using placeholder.")
		return "This is a placeholder. You said " + spokenText
	}
}

func (s *Server) ProcessKnowledgeGraph(req *vtt.KnowledgeGraphRequest) (*vtt.KnowledgeGraphResponse, error) {
	var voiceTimer int = 0
	var transcription1 string = ""
	var transcription2 string = ""
	var transcribedText string
	var isOpus bool
	var micData []int16
	var die bool = false
	var doSTT bool = true
	var sayStarting = true
	var leopardSTT leopard.Leopard
	if os.Getenv("DEBUG_LOGGING") != "true" && os.Getenv("DEBUG_LOGGING") != "false" {
		fmt.Println("No valid value for DEBUG_LOGGING, setting to true")
		debugLogging = true
	} else {
		if os.Getenv("DEBUG_LOGGING") == "true" {
			debugLogging = true
		} else {
			debugLogging = false
		}
	}
	botNumKG = botNumKG + 1
	justThisBotNum := botNumKG
	if botNum > 1 {
		disableLiveTranscription = true
	} else {
		disableLiveTranscription = false
	}
	if justThisBotNum == 1 {
		leopardSTT = leopardSTT1
	} else if justThisBotNum == 2 {
		leopardSTT = leopardSTT2
	} else if justThisBotNum == 3 {
		leopardSTT = leopardSTT3
	} else if justThisBotNum == 4 {
		leopardSTT = leopardSTT4
	} else if justThisBotNum == 5 {
		leopardSTT = leopardSTT5
	} else {
		fmt.Println("Too many bots are connected, sending error to bot " + strconv.Itoa(justThisBotNum))
		NoResultSpoken = knowledgeAPI(req.Session, transcribedText)
		kg := pb.KnowledgeGraphResponse{
			Session:     req.Session,
			DeviceId:    req.Device,
			CommandType: NoResult,
			SpokenText:  NoResultSpoken,
		}
		botNumKG = botNumKG - 1
		if err := req.Stream.Send(&kg); err != nil {
			return nil, err
		}
		return &vtt.KnowledgeGraphResponse{
			Intent: &kg,
		}, nil
		return nil, nil
	}
	if debugLogging == true {
		fmt.Println("(KG) Bot " + strconv.Itoa(botNumKG) + " ESN: " + req.Device)
		fmt.Println("(KG) Bot " + strconv.Itoa(botNumKG) + " Session: " + req.Session)
		fmt.Println("(KG) Bot " + strconv.Itoa(botNumKG) + " Language: " + req.LangString)
		fmt.Println("(KG) Stream " + strconv.Itoa(botNumKG) + " opened.")
	}
	data := []byte{}
	data = append(data, req.FirstReq.InputAudio...)
	if len(data) > 0 {
		if data[0] == 0x4f {
			isOpus = true
			if debugLogging == true {
				fmt.Println("(KG) Bot " + strconv.Itoa(botNumKG) + " Stream Type: Opus")
			}
		} else {
			isOpus = false
			if debugLogging == true {
				fmt.Println("(KG) Bot " + strconv.Itoa(botNumKG) + " Stream Type: PCM")
			}
		}
	}
	stream := opus.OggStream{}
	go func() {
		if isOpus == true {
			time.Sleep(time.Millisecond * 500)
		} else {
			time.Sleep(time.Millisecond * 1100)
		}
		for voiceTimer < 7 {
			voiceTimer = voiceTimer + 1
			time.Sleep(time.Millisecond * 750)
		}
	}()
	go func() {
		if botNumKG > 1 {
			fmt.Println("(KG) Multiple bots are streaming, live transcription disabled")
			disableLiveTranscriptionKG = true
		}
		for doSTT == true {
			if micData != nil && die == false && voiceTimer > 0 {
				if sayStarting == true {
					if debugLogging == true {
						fmt.Printf("(KG) Transcribing stream %d...\n", justThisBotNum)
					}
					sayStarting = false
				}
				processOneData := micData
				transcription1Raw, err := leopardSTT.Process(processOneData)
				if err != nil {
					log.Println(err)
				}
				transcription1 = strings.ToLower(transcription1Raw)
				if debugLogging == true {
					if disableLiveTranscriptionKG == false {
						fmt.Printf("\r(KG) Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcription1)
					}
				}
				if transcription1 != "" && transcription2 != "" && transcription1 == transcription2 {
					transcribedText = transcription1
					if debugLogging == true {
						if disableLiveTranscriptionKG == false {
							fmt.Printf("\n")
						} else {
							fmt.Println("\r(KG) Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
						}
					}
					die = true
					break
				} else if voiceTimer == 7 {
					transcribedText = transcription2
					if debugLogging == true {
						if disableLiveTranscriptionKG == false {
							fmt.Printf("\n")
						} else {
							fmt.Println("\r(KG) Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
						}
					}
					die = true
					break
				}
				time.Sleep(time.Millisecond * 150)
				processTwoData := micData
				transcription2Raw, err := leopardSTT.Process(processTwoData)
				if err != nil {
					log.Println(err)
				}
				transcription2 = strings.ToLower(transcription2Raw)
				if debugLogging == true {
					if disableLiveTranscriptionKG == false {
						fmt.Printf("\r(KG) Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcription2)
					}
				}
				if transcription1 != "" && transcription2 != "" && transcription1 == transcription2 {
					transcribedText = transcription1
					if debugLogging == true {
						if disableLiveTranscriptionKG == false {
							fmt.Printf("\n")
						} else {
							fmt.Println("\r(KG) Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
						}
					}
					die = true
					break
				} else if voiceTimer == 7 {
					transcribedText = transcription2
					if debugLogging == true {
						if disableLiveTranscriptionKG == false {
							fmt.Printf("\n")
						} else {
							fmt.Println("\r(KG) Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
						}
					}
					die = true
					break
				}
				time.Sleep(time.Millisecond * 150)
			}
		}
	}()
	for {
		chunk, err := req.Stream.Recv()
		if err != nil {
			if err == io.EOF {
				transcribedText = ""
				break
			}
		}

		data = append(data, chunk.InputAudio...)
		micData = bytesToInt(stream, data, die, isOpus)
		if die == true {
			break
		}
	}
	NoResultSpoken = knowledgeAPI(req.Session, transcribedText)
	kg := pb.KnowledgeGraphResponse{
		Session:     req.Session,
		DeviceId:    req.Device,
		CommandType: NoResult,
		SpokenText:  NoResultSpoken,
	}
	botNumKG = botNumKG - 1
	if err := req.Stream.Send(&kg); err != nil {
		return nil, err
	}
	return &vtt.KnowledgeGraphResponse{
		Intent: &kg,
	}, nil

}
