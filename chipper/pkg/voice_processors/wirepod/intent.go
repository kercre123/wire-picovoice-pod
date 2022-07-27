package wirepod

import (
	"fmt"
	"strconv"

	"github.com/digital-dream-labs/chipper/pkg/vtt"
)

func (s *Server) ProcessIntent(req *vtt.IntentRequest) (*vtt.IntentResponse, error) {
	var successMatched bool
	transcribedText, transcribedSlots, isRhino, justThisBotNum, isOpus, err := sttHandler(req, false)
	if err != nil {
		IntentPass(req, "intent_system_noaudio", "voice processing error", map[string]string{"error": err.Error()}, true, justThisBotNum)
		return nil, nil
	}
	if isRhino == true {
		successMatched = true
		paramCheckerSlots(req, transcribedText, transcribedSlots, isOpus, justThisBotNum)
	} else {
		successMatched = processTextAll(req, transcribedText, matchListList, intentsList, isOpus, justThisBotNum)
	}
	if successMatched == false {
		if debugLogging == true {
			fmt.Println("No intent was matched.")
		}
		IntentPass(req, "intent_system_noaudio", transcribedText, map[string]string{"": ""}, false, justThisBotNum)
	}
	if debugLogging == true {
		fmt.Println("Bot " + strconv.Itoa(justThisBotNum) + " request served.")
	}
	return nil, nil
}
