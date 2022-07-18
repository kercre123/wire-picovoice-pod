package wirepod

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/chipper/pkg/vtt"
	"github.com/pkg/errors"
	"github.com/soundhound/houndify-sdk-go"
)

var hclient houndify.Client
var houndEnable bool = true

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
	return result["AllResults"].([]interface{})[0].(map[string]interface{})["SpokenResponseLong"].(string), nil
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
			hclient.EnableConversationState()
			fmt.Println("Houndify initialized!")
		}
	} else {
		houndEnable = false
	}
}

var NoResult string = "NoResultCommand"
var NoResultSpoken string

func knowledgeAPI(spokenText string, req *vtt.KnowledgeGraphRequest) string {
	if houndEnable == true {
		hReq := houndify.TextRequest{
			Query:             spokenText,
			UserID:            req.Device,
			RequestID:         req.Session,
			RequestInfoFields: make(map[string]interface{}),
		}
		if debugLogging == true {
			fmt.Println("Making request to Houndify...")
		}
		serverResponse, err := hclient.TextSearch(hReq)
		if err != nil {
			fmt.Println(err)
		}
		robotWords, err := ParseSpokenResponse(serverResponse)
		if err != nil {
			fmt.Println(err)
		}
		if debugLogging == true {
			fmt.Println("Houndify Response: " + robotWords)
		}
		return robotWords
	} else {
		fmt.Println("Houndify is not enabled, using placeholder.")
		return "This is a placeholder. You said " + spokenText
	}
}

func (s *Server) ProcessKnowledgeGraph(req *vtt.KnowledgeGraphRequest) (*vtt.KnowledgeGraphResponse, error) {
	if picovoiceModeOS == "OnlyRhino" {
		fmt.Println("Knowledge Graph does not work in Rhino mode.")
		NoResultSpoken = "Knowledge Graph does not work in Rhino mode."
		kg := pb.KnowledgeGraphResponse{
			Session:     req.Session,
			DeviceId:    req.Device,
			CommandType: NoResult,
			SpokenText:  NoResultSpoken,
		}
		botNum = botNum - 1
		if err := req.Stream.Send(&kg); err != nil {
			return nil, err
		}
		return &vtt.KnowledgeGraphResponse{
			Intent: &kg,
		}, nil
		return nil, nil
	}
	transcribedText, transcribedSlots, isRhino, justThisBotNum, isOpus, err := sttHandler(req, true)
	if transcribedSlots == nil && isRhino == false && isOpus == false {
		// don't do anything
	}
	if err != nil {
		fmt.Println(err)
		NoResultSpoken = err.Error()
		kg := pb.KnowledgeGraphResponse{
			Session:     req.Session,
			DeviceId:    req.Device,
			CommandType: NoResult,
			SpokenText:  NoResultSpoken,
		}
		if err := req.Stream.Send(&kg); err != nil {
			return nil, err
		}
		return &vtt.KnowledgeGraphResponse{
			Intent: &kg,
		}, nil
		return nil, nil
	}
	NoResultSpoken = knowledgeAPI(transcribedText, req)
	kg := pb.KnowledgeGraphResponse{
		Session:     req.Session,
		DeviceId:    req.Device,
		CommandType: NoResult,
		SpokenText:  NoResultSpoken,
	}
	if debugLogging == true {
		fmt.Println("(KG) Bot " + strconv.Itoa(justThisBotNum) + " request served.")
	}
	if err := req.Stream.Send(&kg); err != nil {
		return nil, err
	}
	return &vtt.KnowledgeGraphResponse{
		Intent: &kg,
	}, nil

}
