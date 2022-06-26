package wirepod

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/chipper/pkg/vtt"
	opus "github.com/digital-dream-labs/opus-go/opus"
)

var NoResult string = "NoResultCommand"
var NoResultSpoken string

var botNumKG int = 0

func knowledgeAPI(spokenText string) string {
	// This is where you would make a call to an API
	if strings.Contains(spokenText, "the president of the united states") {
		return "The president of the United States is Joe Biden."
	} else if spokenText == "who is anki" || spokenText == "who was anki" || spokenText == "with mangchi" || spokenText == "is on key" || spokenText == "was on key" {
		return "Anki was the best consumer robotics company. Anki and it's people will be missed."
	} else {
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
		for doSTT == true {
			if micData != nil {
				if die == false {
					if voiceTimer > 1 {
						if sayStarting == true {
							if debugLogging == true {
								fmt.Printf("(KG) Starting transcription...")
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
							fmt.Printf("\r(KG) Bot " + strconv.Itoa(botNumKG) + " Transcription: " + transcription1)
						}
						if transcription1 != "" && transcription2 != "" && transcription1 == transcription2 {
							transcribedText = transcription1
							if debugLogging == true {
								fmt.Printf("\n")
							}
							die = true
							break
						} else if voiceTimer == 7 {
							transcribedText = transcription2
							if debugLogging == true {
								fmt.Printf("\n")
							}
							die = true
							break
						}
						time.Sleep(time.Millisecond * 100)
						processTwoData := micData
						transcription2Raw, err := leopardSTT.Process(processTwoData)
						if err != nil {
							log.Println(err)
						}
						transcription2 = strings.ToLower(transcription2Raw)
						if debugLogging == true {
							fmt.Printf("\r(KG) Bot " + strconv.Itoa(botNumKG) + " Transcription: " + transcription2)
						}
						if transcription1 != "" && transcription2 != "" && transcription1 == transcription2 {
							transcribedText = transcription1
							if debugLogging == true {
								fmt.Printf("\n")
							}
							die = true
							break
						} else if voiceTimer == 7 {
							transcribedText = transcription2
							if debugLogging == true {
								fmt.Printf("\n")
							}
							die = true
							break
						}
						time.Sleep(time.Millisecond * 200)
					}
				}
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
	NoResultSpoken = knowledgeAPI(transcribedText)
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
