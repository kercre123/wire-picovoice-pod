package wirepod

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/digital-dream-labs/chipper/pkg/vtt"

	cheetah "github.com/Picovoice/cheetah/binding/go"
	leopard "github.com/Picovoice/leopard/binding/go"
	rhino "github.com/Picovoice/rhino/binding/go/v2"
	opus "github.com/digital-dream-labs/opus-go/opus"
)

var debugLogging bool

var botNum int = 0
var disableLiveTranscription bool = false

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

func samplesToBytes(buf []int16) []byte {
	output := make([]byte, len(buf)*2)
	for i := 0; i < len(buf); i++ {
		binary.LittleEndian.PutUint16(output[2*i:2*(i+1)], uint16(buf[i]))
	}
	return output
}

func (s *Server) ProcessIntent(req *vtt.IntentRequest) (*vtt.IntentResponse, error) {
	var voiceTimer int = 0
	var transcription1 string = ""
	var transcription2 string = ""
	var transcribedText string = ""
	var isOpus bool
	var micDataLeopard []int16
	var micDataRhino [][]int16
	var die bool = false
	var doSTT bool = true
	var sayStarting = true
	var leopardSTT leopard.Leopard
	var rhinoSTI rhino.Rhino
	var cheetahSTT cheetah.Cheetah
	var leopardFallback bool = false
	var numInRange int = 0
	var oldDataLength int = 0
	var rhinoDone bool = false
	var successMatched bool = false
	var transcribedIntent string
	var transcribedSlots map[string]string
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
	justThisBotNum := botNum
	if botNum > 1 {
		if debugLogging == true {
			fmt.Println("Multiple bots are streaming, live transcription disabled")
		}
		disableLiveTranscription = true
	} else {
		disableLiveTranscription = false
	}
	if picovoiceModeOS == "OnlyLeopard" {
		leopardFallback = true
	}
	if os.Getenv("DISABLE_LIVE_TRANSCRIPTION") == "true" {
		if debugLogging == true {
			fmt.Println("DISABLE_LIVE_TRANSCRIPTION is true, live transcription disabled")
		}
		disableLiveTranscription = true
	}
	if botNum > picovoiceInstances {
		fmt.Println("Too many bots are connected, sending error to bot " + strconv.Itoa(justThisBotNum))
		IntentPass(req, "intent_system_noaudio", "Too many bots, max is 5", map[string]string{"error": "EOF"}, true, justThisBotNum)
		botNum = botNum - 1
		return nil, nil
	} else {
		if picovoiceModeOS == "OnlyLeopard" || picovoiceModeOS == "LeopardAndRhino" {
			leopardSTT = leopardSTTArray[botNum-1]
		}
		if picovoiceModeOS == "OnlyRhino" || picovoiceModeOS == "LeopardAndRhino" {
			rhinoSTI = rhinoSTIArray[botNum-1]
		}
		if picovoiceModeOS == "OnlyCheetah" {
			cheetahSTT = cheetahSTTArray[botNum-1]
		}
	}
	if debugLogging == true {
		fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " ESN: " + req.Device)
		fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Session: " + req.Session)
		fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Language: " + req.LangString)
		fmt.Println("Stream " + strconv.Itoa(justThisBotNum) + " opened.")
	}
	data := []byte{}
	data = append(data, req.FirstReq.InputAudio...)
	if len(data) > 0 {
		if data[0] == 0x4f {
			isOpus = true
			if debugLogging == true {
				fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Stream Type: Opus")
			}
		} else {
			isOpus = false
			if debugLogging == true {
				fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Stream Type: PCM")
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
		if botNum > 1 {
			disableLiveTranscription = true
		}
		if picovoiceModeOS == "OnlyLeopard" || picovoiceModeOS == "LeopardAndRhino" {
			for doSTT == true {
				if leopardFallback == true {
					if micDataLeopard != nil && die == false && voiceTimer > 0 {
						if sayStarting == true {
							if debugLogging == true {
								fmt.Printf("Transcribing stream %d...\n", justThisBotNum)
							}
							sayStarting = false
						}
						processOneData := micDataLeopard
						transcription1Raw, err := leopardSTT.Process(processOneData)
						if err != nil {
							log.Println(err)
						}
						transcription1 = strings.ToLower(transcription1Raw)
						if debugLogging == true {
							if disableLiveTranscription == false {
								fmt.Printf("\rBot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcription1)
							}
						}
						if transcription1 != "" && transcription2 != "" && transcription1 == transcription2 {
							transcribedText = transcription1
							if debugLogging == true {
								if disableLiveTranscription == false {
									fmt.Printf("\n")
								} else {
									fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
								}
							}
							die = true
							break
						} else if voiceTimer == 7 {
							transcribedText = transcription2
							if debugLogging == true {
								if disableLiveTranscription == false {
									fmt.Printf("\n")
								} else {
									fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
								}
							}
							die = true
							break
						}
						time.Sleep(time.Millisecond * 200)
						processTwoData := micDataLeopard
						transcription2Raw, err := leopardSTT.Process(processTwoData)
						if err != nil {
							log.Println(err)
						}
						transcription2 = strings.ToLower(transcription2Raw)
						if debugLogging == true {
							if disableLiveTranscription == false {
								fmt.Printf("\rBot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcription2)
							}
						}
						if transcription1 != "" && transcription2 != "" && transcription1 == transcription2 {
							transcribedText = transcription1
							if debugLogging == true {
								if disableLiveTranscription == false {
									fmt.Printf("\n")
								} else {
									fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
								}
							}
							die = true
							break
						} else if voiceTimer == 7 {
							transcribedText = transcription2
							if debugLogging == true {
								if disableLiveTranscription == false {
									fmt.Printf("\n")
								} else {
									fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
								}
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
				IntentPass(req, "intent_system_noaudio", "EOF error", map[string]string{"error": "EOF"}, true, justThisBotNum)
				break
			}
		}
		if die == true {
			break
		}
		data = append(data, chunk.InputAudio...)
		// returns []int16, framesize unknown
		micDataLeopard = bytesToIntLeopard(stream, data, die, isOpus)
		if picovoiceModeOS == "OnlyRhino" || picovoiceModeOS == "LeopardAndRhino" {
			// returns [][]int16, 512 framesize
			micDataRhino = bytesToIntRhino(stream, data, die, isOpus)
			numInRange = 0
			for _, sample := range micDataRhino {
				if rhinoDone == false {
					if numInRange >= oldDataLength {
						isFinalized, err := rhinoSTI.Process(sample)
						if isFinalized {
							inference, err := rhinoSTI.GetInference()
							if err != nil {
								fmt.Println("Error getting inference: " + err.Error())
							}
							if inference.IsUnderstood {
								transcribedIntent = inference.Intent
								transcribedSlots = inference.Slots
								die = true
								leopardFallback = false
								rhinoDone = true
								successMatched = true
								break
							} else {
								leopardFallback = true
								rhinoDone = true
								successMatched = false
								if picovoiceModeOS == "OnlyRhino" {
									die = true
								} else {
									fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + ": " + "Rhino STI failed, falling back to Leopard STT")
								}
								break
							}
						}
						if err != nil {
							fmt.Println("Error: " + err.Error())
							break
						}
						numInRange = numInRange + 1
					} else {
						numInRange = numInRange + 1
					}
				}
			}
			oldDataLength = len(micDataRhino)
			if die == true {
				break
			}
		}
		if picovoiceModeOS == "OnlyCheetah" {
			// returns [][]int16, 512 framesize
			micDataRhino = bytesToIntRhino(stream, data, die, isOpus)
			numInRange = 0
			for _, sample := range micDataRhino {
				if numInRange >= oldDataLength {
					if sayStarting == true {
						if debugLogging == true {
							fmt.Printf("Transcribing stream %d...\n", justThisBotNum)
						}
						sayStarting = false
					}
					partialTranscript, isEndpoint, err := cheetahSTT.Process(sample)
					if partialTranscript != "" {
						transcribedText = strings.ToLower(transcribedText + strings.TrimSpace(partialTranscript) + " ")
						if debugLogging == true && disableLiveTranscription == false {
							fmt.Printf("\rBot " + strconv.Itoa(justThisBotNum) + " Transcription: " + transcribedText)
						}
					}
					if isEndpoint {
						finalTranscript, err := cheetahSTT.Flush()
						transcribedText = strings.TrimSpace(strings.ToLower(transcribedText + finalTranscript))
						if debugLogging == true {
							if disableLiveTranscription == false {
								fmt.Printf("\rBot " + strconv.Itoa(justThisBotNum) + " Final Transcription: " + transcribedText + "\n")
							} else {
								fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Final Transcription: " + transcribedText)
							}
						}
						successMatched = false
						die = true
						break
						if err != nil {
							fmt.Println("Error: " + err.Error())
						}
					}
					if err != nil {
						fmt.Println("Error: " + err.Error())
						break
					}
					numInRange = numInRange + 1
				} else {
					numInRange = numInRange + 1
				}
			}
			oldDataLength = len(micDataRhino)
			if die == true {
				break
			}
		}
	}
	if picovoiceModeOS == "OnlyLeopard" || picovoiceModeOS == "LeopardAndRhino" {
		if leopardFallback == true {
			successMatched = processTextAll(req, transcribedText, matchListList, intentsList, isOpus, justThisBotNum)
		} else {
			paramCheckerSlots(req, transcribedIntent, transcribedSlots, isOpus, justThisBotNum)
		}
	} else if picovoiceModeOS == "OnlyRhino" {
		if successMatched == true {
			paramCheckerSlots(req, transcribedIntent, transcribedSlots, isOpus, justThisBotNum)
		}
	} else if picovoiceModeOS == "OnlyCheetah" {
		successMatched = processTextAll(req, transcribedText, matchListList, intentsList, isOpus, justThisBotNum)
	}
	if successMatched == false {
		if debugLogging == true {
			fmt.Println("No intent was matched.")
		}
		IntentPass(req, "intent_system_noaudio", transcribedText, map[string]string{"": ""}, false, justThisBotNum)
	}
	botNum = botNum - 1
	if debugLogging == true {
		fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " request served.")
	}
	return nil, nil
}
