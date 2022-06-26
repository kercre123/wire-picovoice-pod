package wirepod

import (
	"fmt"
	"strings"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/chipper/pkg/vtt"
)

func IntentPass(req *vtt.IntentRequest, intentThing string, speechText string, intentParams map[string]string, isParam bool) (*vtt.IntentResponse, error) {
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
		fmt.Println("Intent Sent: " + intentThing)
		if isParam == true {
			fmt.Println("Parameters Sent:", intentParams)
		} else {
			fmt.Println("No Parameters Sent")
		}
	}
	return r, nil
}

func processTextAll(req *vtt.IntentRequest, voiceText string, listOfLists [][]string, intentList []string, isOpus bool) int {
	var matched int = 0
	var intentNum int = 0
	var successMatched int = 0
	for _, b := range listOfLists {
		for _, c := range b {
			if strings.Contains(voiceText, c) {
				if isOpus == true {
					paramChecker(req, intentList[intentNum], voiceText)
				} else {
					prehistoricParamChecker(req, intentList[intentNum], voiceText)
				}
				successMatched = 1
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
