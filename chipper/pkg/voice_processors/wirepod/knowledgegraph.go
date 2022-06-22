package wirepod

import (
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

func (s *Server) ProcessKnowledgeGraph(req *vtt.KnowledgeGraphRequest) (*vtt.KnowledgeGraphResponse, error) {
	var voiceTimer int = 0
	var transcription1 string = ""
	var transcription2 string = ""
	var transcription3 string = ""
	var transcription4 string = ""
	var successMatch bool = false
	var processingOne bool = false
	var processingTwo bool = false
	var processingThree bool = false
	var processingFour bool = false
	var transcribedText string
	var micData []int16
	var die bool = false
	if os.Getenv("DEBUG_LOGGING") != "true" && os.Getenv("DEBUG_LOGGING") != "false" {
		log.Println("No valid value for DEBUG_LOGGING, setting to true")
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
		log.Println("Stream " + strconv.Itoa(botNumKG) + " opened.")
	}
	data := []byte{}
	data = append(data, req.FirstReq.InputAudio...)
	stream := opus.OggStream{}
	go func() {
		time.Sleep(time.Millisecond * 500)
		for voiceTimer < 7 {
			voiceTimer = voiceTimer + 1
			time.Sleep(time.Second * 1)
		}
	}()
	go func() {
		for voiceTimer < 7 {
			if micData != nil {
				if die == false {
					if voiceTimer == 1 {
						if processingOne == false {
							processingOne = true
							processOneData := micData
							transcription1Raw, err := leopardSTT.Process(processOneData)
							if err != nil {
								log.Println(err)
							}
							transcription1 = strings.ToLower(transcription1Raw)
							log.Println("Bot " + strconv.Itoa(botNumKG) + ", Transcription 1: " + transcription1)
						}
					}
					if voiceTimer == 2 {
						if processingTwo == false {
							processingTwo = true
							processTwoData := micData
							transcription2Raw, err := leopardSTT.Process(processTwoData)
							if err != nil {
								log.Println(err)
							}
							transcription2 = strings.ToLower(transcription2Raw)
							log.Println("Bot " + strconv.Itoa(botNumKG) + ", Transcription 2: " + transcription2)
						}
					}
					if voiceTimer == 3 {
						if processingThree == false {
							processingThree = true
							processThreeData := micData
							transcription3Raw, err := leopardSTT.Process(processThreeData)
							if err != nil {
								log.Println(err)
							}
							transcription3 = strings.ToLower(transcription3Raw)
							log.Println("Bot " + strconv.Itoa(botNumKG) + ", Transcription 3: " + transcription3)
						}
					}
					if voiceTimer == 4 {
						if processingFour == false {
							processingFour = true
							processFourData := micData
							transcription4Raw, err := leopardSTT.Process(processFourData)
							if err != nil {
								log.Println(err)
							}
							transcription4 = strings.ToLower(transcription4Raw)
							log.Println("Bot " + strconv.Itoa(botNumKG) + ", Transcription 4: " + transcription4)
							successMatch = true
						}
					}
				}
			}
		}
	}()
	for {
		chunk, err := req.Stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if transcription2 != "" {
			if transcription1 == transcription2 {
				log.Println("Bot " + strconv.Itoa(botNumKG) + ", " + "Speech stopped, 2: " + transcription1)
				transcribedText = transcription1
				die = true
				break
			} else if transcription2 != "" {
				if transcription2 == transcription3 {
					log.Println("Bot " + strconv.Itoa(botNumKG) + ", " + "Speech stopped, 3: " + transcription2)
					transcribedText = transcription2
					die = true
					break
				} else if transcription3 != "" {
					if transcription3 == transcription4 {
						log.Println("Bot " + strconv.Itoa(botNumKG) + ", " + "Speech stopped, 4: " + transcription3)
						transcribedText = transcription3
						die = true
						break
					} else if transcription4 != "" {
						if transcription3 == transcription4 {
							log.Println("Bot " + strconv.Itoa(botNumKG) + ", " + "Speech stopped, 4: " + transcription4)
							transcribedText = transcription4
							die = true
							break
						} else {
							log.Println("Bot " + strconv.Itoa(botNumKG) + ", " + "Speech stopped, 4 (nm): " + transcription4)
							transcribedText = transcription4
							die = true
							break
						}
					}
				}
			}
		}
		if transcription2 == "" && transcription3 != "" {
			if transcription4 != "" {
				if transcription3 == transcription4 {
					log.Println("Speech stopped, 4: " + transcription4)
					transcribedText = transcription4
					die = true
					break
				} else {
					log.Println("Speech stopped, 4 (nm): " + transcription4)
					transcribedText = transcription4
					die = true
					break
				}
			}
		}
		if transcription3 == "" && transcription4 != "" {
			log.Println("Speech stopped, 4 (nm): " + transcription4)
			transcribedText = transcription4
			die = true
			break
		}
		if transcription4 == "" && successMatch == true {
			transcribedText = ""
			die = true
			break
		}
		data = append(data, chunk.InputAudio...)
		micData = bytesToInt(stream, data, die)
	}
	NoResultSpoken = "This is a placeholder! You said: " + transcribedText
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
