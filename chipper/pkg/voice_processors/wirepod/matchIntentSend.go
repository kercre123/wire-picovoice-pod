package wirepod

import (
	"fmt"
	"strconv"
	"strings"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/chipper/pkg/vtt"
)

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

func processTextAll(req *vtt.IntentRequest, voiceText string, listOfLists [][]string, intentList []string, isOpus bool, justThisBotNum int) bool {
	var matched int = 0
	var intentNum int = 0
	var successMatched bool = false
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
	return successMatched
}
