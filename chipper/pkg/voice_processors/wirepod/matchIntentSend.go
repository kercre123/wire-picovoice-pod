package wirepod

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/chipper/pkg/vtt"
)

type intentsStruct []struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Utterances  []string `json:"utterances"`
	Intent      string   `json:"intent"`
	Params      struct {
		ParamName  string `json:"paramname"`
		ParamValue string `json:"paramvalue"`
	} `json:"params"`
	Exec string `json:"exec"`
}

func IntentPass(req *vtt.IntentRequest, intentThing string, speechText string, intentParams map[string]string, isParam bool, justThisBotNum int) (*vtt.IntentResponse, error) {
	intent := pb.IntentResponse{
		IsFinal: true,
		IntentResult: &pb.IntentResult{
			QueryText:  speechText,
			Action:     intentThing,
			Parameters: intentParams,
		},
	}
	if err := req.Stream.Send(&intent); err != nil {
		return nil, err
	}
	r := &vtt.IntentResponse{
		Intent: &intent,
	}
	if debugLogging == true {
		fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " Intent Sent: " + intentThing)
		if isParam == true {
			fmt.Println("Bot "+strconv.Itoa(justThisBotNum)+" Parameters Sent:", intentParams)
		} else {
			fmt.Println("No Parameters Sent")
		}
	}
	return r, nil
}

func customIntentHandler(req *vtt.IntentRequest, voiceText string, intentList []string, isOpus bool, justThisBotNum int) bool {
	var successMatched bool = false
	if _, err := os.Stat("./customIntents.json"); err == nil {
		var customIntentJSON intentsStruct
		customIntentJSONFile, err := os.ReadFile("./customIntents.json")
		json.Unmarshal(customIntentJSONFile, &customIntentJSON)
		for _, c := range customIntentJSON {
			for _, v := range c.Utterances {
				if strings.Contains(voiceText, v) {
					if debugLogging == true {
						fmt.Println("Custom Intent Matched: " + c.Name + " - " + c.Description + " - " + c.Intent)
					}
					var intentParams map[string]string
					var isParam bool = false
					if c.Params.ParamValue != "" {
						if debugLogging == true {
							fmt.Println("Custom Intent Parameter: " + c.Params.ParamName + " - " + c.Params.ParamValue)
						}
						intentParams = map[string]string{c.Params.ParamName: c.Params.ParamValue}
						isParam = true
					}
					customIntentExec := exec.Command("/bin/bash", c.Exec)
					customOut, err := customIntentExec.Output()
					if err != nil {
						fmt.Println(err)
					}
					if debugLogging == true {
						fmt.Println("Custom Intent Exec Output: " + string(customOut))
					}
					IntentPass(req, c.Intent, voiceText, intentParams, isParam, justThisBotNum)
					successMatched = true
					break
				}
				if successMatched == true {
					break
				}
			}
		}
		if err != nil {
			fmt.Println(err)
		}

	}
	return successMatched
}

func processTextAll(req *vtt.IntentRequest, voiceText string, listOfLists [][]string, intentList []string, isOpus bool, justThisBotNum int) bool {
	var matched int = 0
	var intentNum int = 0
	var successMatched bool = false
	customIntentMatched := customIntentHandler(req, voiceText, intentList, isOpus, justThisBotNum)
	if customIntentMatched == false {
		for _, b := range listOfLists {
			for _, c := range b {
				if strings.Contains(voiceText, c) {
					if isOpus == true {
						paramChecker(req, intentList[intentNum], voiceText, justThisBotNum)
					} else {
						prehistoricParamChecker(req, intentList[intentNum], voiceText, justThisBotNum)
					}
					successMatched = true
					matched = 1
					break
				}
			}
			if matched == 1 {
				matched = 0
				break
			}
			intentNum = intentNum + 1
		}
	} else {
		successMatched = true
	}
	return successMatched
}
