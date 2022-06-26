package wirepod

import (
	"fmt"
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
		- current workaround: setup specific bots with botSetup.sh
*/

func InitLeopard() {
	leopardKey := os.Getenv("LEOPARD_APIKEY")
	if leopardKey == "" {
		fmt.Println("You must set LEOPARD_APIKEY to a value.")
		os.Exit(1)
	}
	fmt.Println("Initializing Picovoice Leopard...")
	leopardSTT = leopard.Leopard{AccessKey: leopardKey}
	err := leopardSTT.Init()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Initialized Picovoice Leopard")
}

func (s *Server) ProcessIntent(req *vtt.IntentRequest) (*vtt.IntentResponse, error) {
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
	botNum = botNum + 1
	if debugLogging == true {
		fmt.Println("Bot " + strconv.Itoa(botNum) + " ESN: " + req.Device)
		fmt.Println("Bot " + strconv.Itoa(botNum) + " Session: " + req.Session)
		fmt.Println("Bot " + strconv.Itoa(botNum) + " Language: " + req.LangString)
		fmt.Println("Stream " + strconv.Itoa(botNum) + " opened.")
	}
	data := []byte{}
	data = append(data, req.FirstReq.InputAudio...)
	if len(data) > 0 {
		if data[0] == 0x4f {
			isOpus = true
			if debugLogging == true {
				fmt.Println("Bot " + strconv.Itoa(botNum) + " Stream Type: Opus")
			}
		} else {
			isOpus = false
			if debugLogging == true {
				fmt.Println("Bot " + strconv.Itoa(botNum) + " Stream Type: PCM")
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
					if voiceTimer > 0 {
						if sayStarting == true {
							if debugLogging == true {
								fmt.Printf("Starting transcription...")
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
							fmt.Printf("\rBot " + strconv.Itoa(botNum) + " Transcription: " + transcription1)
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
						time.Sleep(time.Millisecond * 150)
						processTwoData := micData
						transcription2Raw, err := leopardSTT.Process(processTwoData)
						if err != nil {
							log.Println(err)
						}
						transcription2 = strings.ToLower(transcription2Raw)
						if debugLogging == true {
							fmt.Printf("\rBot " + strconv.Itoa(botNum) + " Transcription: " + transcription2)
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
						time.Sleep(time.Millisecond * 150)
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
		if die == true {
			break
		}
		data = append(data, chunk.InputAudio...)
		micData = bytesToInt(stream, data, die, isOpus)
	}
	successMatched := processTextAll(req, transcribedText, matchListList, intentsList, isOpus)
	if successMatched == 0 {
		if debugLogging == true {
			fmt.Println("No intent was matched.")
		}
		IntentPass(req, "intent_system_noaudio", transcribedText, map[string]string{"": ""}, false)
	}
	botNum = botNum - 1
	return nil, nil
}
