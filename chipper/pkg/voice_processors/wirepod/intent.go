package wirepod

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/digital-dream-labs/chipper/pkg/vtt"

	leopard "github.com/Picovoice/leopard/binding/go"
	opus "github.com/digital-dream-labs/opus-go/opus"
)

var leopardSTT leopard.Leopard

var debugLogging bool

var botNum int = 0

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
	TODO:
	1. Implement jdocs. These are files which are stored on the bot which contain the bot's
	default location, unit settings, etc. Helpful for weather.
		- workaround, ask the user for settings during setup.sh
	3. Overall take shell out of the picture (https://github.com/asticode/go-asticoqui)
	4. Maybe find a way to detect silence in the audio for better end handling.
		- probably unnecessary
*/

func InitLeopard() {
	leopardKey := os.Getenv("LEOPARD_APIKEY")
	if leopardKey == "" {
		log.Println("You must set LEOPARD_APIKEY to a value.")
		os.Exit(1)
	}
	log.Println("Initializing Leopard")
	leopardSTT = leopard.Leopard{AccessKey: leopardKey}
	err := leopardSTT.Init()
	if err != nil {
		log.Println(err)
	}
	log.Println("Initialized Leopard")
}

func (s *Server) ProcessIntent(req *vtt.IntentRequest) (*vtt.IntentResponse, error) {
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
	var isOpus bool
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
	botNum = botNum + 1
	log.Println("Bot " + strconv.Itoa(botNum) + " ESN: " + req.Device)
	log.Println("Bot " + strconv.Itoa(botNum) + " Session: " + req.Session)
	log.Println("Bot " + strconv.Itoa(botNum) + " Language: " + req.LangString)
	if debugLogging == true {
		log.Println("Stream " + strconv.Itoa(botNum) + " opened.")
	}
	data := []byte{}
	data = append(data, req.FirstReq.InputAudio...)
	if len(data) > 0 {
		if data[0] == 0x4f {
			isOpus = true
			log.Println("Bot " + strconv.Itoa(botNum) + " Stream Type: Opus")
		} else {
			isOpus = false
			log.Println("Bot " + strconv.Itoa(botNum) + " Stream Type: PCM")
		}
	}
	stream := opus.OggStream{}
	go func() {
		if isOpus == true {
			time.Sleep(time.Millisecond * 300)
		} else {
			time.Sleep(time.Millisecond * 1100)
		}
		for voiceTimer < 7 {
			voiceTimer = voiceTimer + 1
			time.Sleep(time.Millisecond * 700)
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
							log.Println("Bot " + strconv.Itoa(botNum) + ", Transcription 1: " + transcription1)
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
							log.Println("Bot " + strconv.Itoa(botNum) + ", Transcription 2: " + transcription2)
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
							log.Println("Bot " + strconv.Itoa(botNum) + ", Transcription 3: " + transcription3)
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
							log.Println("Bot " + strconv.Itoa(botNum) + ", Transcription 4: " + transcription4)
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
				IntentPass(req, "intent_system_noaudio", "EOF error", map[string]string{"error": "EOF"}, true)
				break
			}
		}
		if transcription2 != "" {
			if transcription1 == transcription2 {
				log.Println("Bot " + strconv.Itoa(botNum) + ", " + "Speech stopped, 2: " + transcription1)
				transcribedText = transcription1
				die = true
				break
			} else if transcription2 != "" {
				if transcription2 == transcription3 {
					log.Println("Bot " + strconv.Itoa(botNum) + ", " + "Speech stopped, 3: " + transcription2)
					transcribedText = transcription2
					die = true
					break
				} else if transcription3 != "" {
					if transcription3 == transcription4 {
						log.Println("Bot " + strconv.Itoa(botNum) + ", " + "Speech stopped, 4: " + transcription3)
						transcribedText = transcription3
						die = true
						break
					} else if transcription4 != "" {
						if transcription3 == transcription4 {
							log.Println("Bot " + strconv.Itoa(botNum) + ", " + "Speech stopped, 4: " + transcription4)
							transcribedText = transcription4
							die = true
							break
						} else {
							log.Println("Bot " + strconv.Itoa(botNum) + ", " + "Speech stopped, 4 (nm): " + transcription4)
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
		micData = bytesToInt(stream, data, die, isOpus)
	}
	successMatched := processTextAll(req, transcribedText, matchListList, intentsList, isOpus)
	if successMatched == 0 {
		if debugLogging == true {
			log.Println("No intent was matched.")
		}
		IntentPass(req, "intent_system_noaudio", transcribedText, map[string]string{"": ""}, false)
	}
	botNum = botNum - 1
	return nil, nil
}
