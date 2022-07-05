package wirepod

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/digital-dream-labs/chipper/pkg/vtt"
)

func paramCheckerSlots(req *vtt.IntentRequest, intent string, slots map[string]string, isOpus bool, justThisBotNum int) {
	var intentParam string
	var intentParamValue string
	var newIntent string
	var isParam bool
	var intentParams map[string]string
	var botLocation string = "San Francisco"
	var botUnits string = "F"
	var botPlaySpecific bool = false
	var botIsEarlyOpus bool = false
	if _, err := os.Stat("./botConfigs/" + req.Device + ".json"); err == nil {
		if debugLogging == true {
			fmt.Println("Found bot config file for " + req.Device)
		}
		botConfigByte, err := os.ReadFile("./botConfigs/" + req.Device + ".json")
		if err != nil {
			fmt.Println(err)
		}
		type botConfigJSON struct {
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}
		var botConfJSON botConfigJSON
		json.Unmarshal(botConfigByte, &botConfJSON)
		botLocation = botConfJSON.Location
		botUnits = botConfJSON.Units
		botPlaySpecific = botConfJSON.UsePlaySpecific
		botIsEarlyOpus = botConfJSON.IsEarlyOpus
	}
	if botPlaySpecific == true {
		if strings.Contains(intent, "intent_play_blackjack") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "blackjack"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_fistbump") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "fist_bump"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_rollcube") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "roll_cube"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_popawheelie") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "pop_a_wheelie"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_pickupcube") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "pick_up_cube"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_keepaway") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "keep_away"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else {
			newIntent = intent
			intentParam = ""
			intentParamValue = ""
			isParam = false
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	if strings.Contains(intent, "volume") {
		if slots["volume"] != "" {
			newIntent = "intent_imperative_volumelevel_extend"
			isParam = true
			intentParam = "volume_level"
			if strings.Contains(slots["volume"], "medium low") {
				intentParamValue = "VOLUME_2"
			} else if strings.Contains(slots["volume"], "low") {
				intentParamValue = "VOLUME_1"
			} else if strings.Contains(slots["volume"], "medium high") {
				intentParamValue = "VOLUME_4"
			} else if strings.Contains(slots["volume"], "high") {
				intentParamValue = "VOLUME_5"
			} else if strings.Contains(slots["volume"], "medium") {
				intentParamValue = "VOLUME_3"
			} else {
				intentParamValue = "VOLUME_1"
			}
		} else {
			isParam = false
			intentParam = ""
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "eyecolor") {
		isParam = true
		newIntent = "intent_imperative_eyecolor_specific_extend"
		intentParam = "eye_color"
		if strings.Contains(slots["eye_color"], "purple") {
			intentParamValue = "COLOR_PURPLE"
		} else if strings.Contains(slots["eye_color"], "blue") || strings.Contains(slots["eye_color"], "sapphire") {
			intentParamValue = "COLOR_BLUE"
		} else if strings.Contains(slots["eye_color"], "yellow") {
			intentParamValue = "COLOR_YELLOW"
		} else if strings.Contains(slots["eye_color"], "teal") || strings.Contains(slots["eye_color"], "tell") {
			intentParamValue = "COLOR_TEAL"
		} else if strings.Contains(slots["eye_color"], "green") {
			intentParamValue = "COLOR_GREEN"
		} else if strings.Contains(slots["eye_color"], "orange") {
			intentParamValue = "COLOR_ORANGE"
		} else {
			newIntent = intent
			intentParamValue = ""
			isParam = false
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "_selfie") {
		newIntent = "intent_photo_take_extend"
		intentParam = "entity_photo_selfie"
		intentParamValue = "photo_selfie"
		isParam = true
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "_noselfie") {
		newIntent = "intent_photo_take_extend"
		intentParam = "entity_photo_selfie"
		intentParamValue = ""
		isParam = true
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "settimer") {
		isParam = true
		newIntent = intent
		slotNum := slots["num"]
		slotUnit := slots["unit"]
		timerSecs, err := strconv.Atoi(slotNum)
		if err != nil {
			fmt.Println(err)
		}
		if slotNum != "" && slotUnit != "" {
			if strings.Contains(slotUnit, "minute") {
				timerSecs = timerSecs * 60
			} else if strings.Contains(slotUnit, "hour") {
				timerSecs = timerSecs * 60 * 60
			}
		}
		if debugLogging == true {
			fmt.Println("Seconds parsed from speech: " + strconv.Itoa(timerSecs))
		}
		intentParam = "timer_duration"
		intentParamValue = strconv.Itoa(timerSecs)
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "global_stop_extend") {
		isParam = true
		newIntent = intent
		intentParam = "what_to_stop"
		intentParamValue = "timer"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_knowledgegraph_prompt") {
		isParam = false
		newIntent = "intent_knowledge_promptquestion"
		intentParam = ""
		intentParamValue = ""
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_weather_extend") {
		isParam = true
		newIntent = intent
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit := weatherParser("no", botLocation, botUnits)
		intentParams = map[string]string{"condition": condition, "is_forecast": is_forecast, "local_datetime": local_datetime, "speakable_location_string": speakable_location_string, "temperature": temperature, "temperature_unit": temperature_unit}
	} else {
		if intentParam == "" {
			newIntent = intent
			intentParam = ""
			intentParamValue = ""
			isParam = false
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	if isOpus == false || botIsEarlyOpus == true || botPlaySpecific == true {
		if strings.Contains(intent, "intent_play_blackjack") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "blackjack"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_fistbump") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "fist_bump"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_rollcube") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "roll_cube"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_praise") {
			isParam = false
			newIntent = "intent_imperative_affirmative"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_abuse") {
			isParam = false
			newIntent = "intent_imperative_negative"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	IntentPass(req, newIntent, intent, intentParams, isParam, justThisBotNum)
}

func paramChecker(req *vtt.IntentRequest, intent string, speechText string, justThisBotNum int) {
	var intentParam string
	var intentParamValue string
	var newIntent string
	var isParam bool
	var intentParams map[string]string
	var botLocation string = "San Francisco"
	var botUnits string = "F"
	var botPlaySpecific bool = false
	var botIsEarlyOpus bool = false
	if _, err := os.Stat("./botConfigs/" + req.Device + ".json"); err == nil {
		if debugLogging == true {
			fmt.Println("Found bot config file for " + req.Device)
		}
		botConfigByte, err := os.ReadFile("./botConfigs/" + req.Device + ".json")
		if err != nil {
			fmt.Println(err)
		}
		type botConfigJSON struct {
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}
		var botConfJSON botConfigJSON
		json.Unmarshal(botConfigByte, &botConfJSON)
		botLocation = botConfJSON.Location
		botUnits = botConfJSON.Units
		botPlaySpecific = botConfJSON.UsePlaySpecific
		botIsEarlyOpus = botConfJSON.IsEarlyOpus
	}
	if botPlaySpecific == true {
		if strings.Contains(intent, "intent_play_blackjack") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "blackjack"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_fistbump") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "fist_bump"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_rollcube") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "roll_cube"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_popawheelie") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "pop_a_wheelie"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_pickupcube") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "pick_up_cube"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_keepaway") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "keep_away"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else {
			newIntent = intent
			intentParam = ""
			intentParamValue = ""
			isParam = false
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	if strings.Contains(intent, "intent_photo_take_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, "me") || strings.Contains(speechText, "self") {
			intentParam = "entity_photo_selfie"
			intentParamValue = "photo_selfie"
		} else {
			intentParam = "entity_photo_selfie"
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_eyecolor") {
		isParam = true
		newIntent = "intent_imperative_eyecolor_specific_extend"
		intentParam = "eye_color"
		if strings.Contains(speechText, "purple") {
			intentParamValue = "COLOR_PURPLE"
		} else if strings.Contains(speechText, "blue") || strings.Contains(speechText, "sapphire") {
			intentParamValue = "COLOR_BLUE"
		} else if strings.Contains(speechText, "yellow") {
			intentParamValue = "COLOR_YELLOW"
		} else if strings.Contains(speechText, "teal") || strings.Contains(speechText, "tell") {
			intentParamValue = "COLOR_TEAL"
		} else if strings.Contains(speechText, "green") {
			intentParamValue = "COLOR_GREEN"
		} else if strings.Contains(speechText, "orange") {
			intentParamValue = "COLOR_ORANGE"
		} else {
			newIntent = intent
			intentParamValue = ""
			isParam = false
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_weather_extend") {
		isParam = true
		newIntent = intent
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit := weatherParser(speechText, botLocation, botUnits)
		intentParams = map[string]string{"condition": condition, "is_forecast": is_forecast, "local_datetime": local_datetime, "speakable_location_string": speakable_location_string, "temperature": temperature, "temperature_unit": temperature_unit}
	} else if strings.Contains(intent, "intent_imperative_volumelevel_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, "medium lo") || strings.Contains(speechText, "media lo") || strings.Contains(speechText, "medium bo") || strings.Contains(speechText, "media bo") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_2"
		} else if strings.Contains(speechText, "low") || strings.Contains(speechText, "quiet") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else if strings.Contains(speechText, "medium high") || strings.Contains(speechText, "media high") || strings.Contains(speechText, "medium hide") || strings.Contains(speechText, "media hide") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_4"
		} else if strings.Contains(speechText, "medium") || strings.Contains(speechText, "normal") || strings.Contains(speechText, "regular") || strings.Contains(speechText, "media") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_3"
		} else if strings.Contains(speechText, "high") || strings.Contains(speechText, "loud") || strings.Contains(speechText, "hide") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_5"
		} else if strings.Contains(speechText, "mute") || strings.Contains(speechText, "nothing") || strings.Contains(speechText, "silent") || strings.Contains(speechText, "off") || strings.Contains(speechText, "zero") || strings.Contains(speechText, "meet") {
			// there is no VOLUME_0 :(
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_names_username_extend") {
		var username string
		var nameSplitter string
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, " is ") {
			nameSplitter = " is "
		} else if strings.Contains(speechText, "'s") {
			nameSplitter = "'s"
		} else if strings.Contains(speechText, "names") {
			nameSplitter = "names"
		}
		if strings.Contains(speechText, " is ") || strings.Contains(speechText, "'s") || strings.Contains(speechText, "names") {
			splitPhrase := strings.SplitAfter(speechText, nameSplitter)
			username = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				username = username + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			if debugLogging == true {
				fmt.Println("Name parsed from speech: " + "`" + username + "`")
			}
			intentParam = "username"
			intentParamValue = username
			intentParams = map[string]string{intentParam: intentParamValue}
		} else {
			if debugLogging == true {
				fmt.Println("No name parsed from speech")
			}
			intentParam = "username"
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	} else if strings.Contains(intent, "intent_clock_settimer_extend") {
		isParam = true
		newIntent = intent
		timerSecs := words2num(speechText)
		if debugLogging == true {
			fmt.Println("Seconds parsed from speech: " + timerSecs)
		}
		intentParam = "timer_duration"
		intentParamValue = timerSecs
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_global_stop_extend") {
		isParam = true
		newIntent = intent
		intentParam = "what_to_stop"
		intentParamValue = "timer"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_playmessage_extend") {
		var given_name string
		isParam = true
		newIntent = intent
		intentParam = "given_name"
		if strings.Contains(speechText, " for ") {
			splitPhrase := strings.SplitAfter(speechText, " for ")
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_recordmessage_extend") {
		var given_name string
		isParam = true
		newIntent = intent
		intentParam = "given_name"
		if strings.Contains(speechText, " for ") {
			splitPhrase := strings.SplitAfter(speechText, " for ")
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else {
		if intentParam == "" {
			newIntent = intent
			intentParam = ""
			intentParamValue = ""
			isParam = false
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	if botIsEarlyOpus == true {
		if strings.Contains(intent, "intent_imperative_praise") {
			isParam = false
			newIntent = "intent_imperative_affirmative"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_abuse") {
			isParam = false
			newIntent = "intent_imperative_negative"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_love") {
			isParam = false
			newIntent = "intent_greeting_hello"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	IntentPass(req, newIntent, speechText, intentParams, isParam, justThisBotNum)
}

func prehistoricParamChecker(req *vtt.IntentRequest, intent string, speechText string, justThisBotNum int) {
	// intent.go detects if the stream uses opus or PCM.
	// If the stream is PCM, it is likely a bot with 0.10.
	// This accounts for the newer 0.10.1### builds.
	var intentParam string
	var intentParamValue string
	var newIntent string
	var isParam bool
	var intentParams map[string]string
	var botLocation string = "San Francisco"
	var botUnits string = "F"
	if _, err := os.Stat("./botConfigs/" + req.Device + ".json"); err == nil {
		if debugLogging == true {
			fmt.Println("Found bot config file for " + req.Device)
		}
		botConfigByte, err := os.ReadFile("./botConfigs/" + req.Device + ".json")
		if err != nil {
			fmt.Println(err)
		}
		type botConfigJSON struct {
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
		}
		var botConfJSON botConfigJSON
		json.Unmarshal(botConfigByte, &botConfJSON)
		botLocation = botConfJSON.Location
		botUnits = botConfJSON.Units
	}
	if strings.Contains(intent, "intent_photo_take_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, "me") || strings.Contains(speechText, "self") {
			intentParam = "entity_photo_selfie"
			intentParamValue = "photo_selfie"
		} else {
			intentParam = "entity_photo_selfie"
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_eyecolor") {
		// leaving stuff like this in case someone wants to add features like this to older software
		isParam = true
		newIntent = "intent_imperative_eyecolor_specific_extend"
		intentParam = "eye_color"
		if strings.Contains(speechText, "purple") {
			intentParamValue = "COLOR_PURPLE"
		} else if strings.Contains(speechText, "blue") || strings.Contains(speechText, "sapphire") {
			intentParamValue = "COLOR_BLUE"
		} else if strings.Contains(speechText, "yellow") {
			intentParamValue = "COLOR_YELLOW"
		} else if strings.Contains(speechText, "teal") || strings.Contains(speechText, "tell") {
			intentParamValue = "COLOR_TEAL"
		} else if strings.Contains(speechText, "green") {
			intentParamValue = "COLOR_GREEN"
		} else if strings.Contains(speechText, "orange") {
			intentParamValue = "COLOR_ORANGE"
		} else {
			newIntent = intent
			intentParamValue = ""
			isParam = false
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_weather_extend") {
		isParam = true
		newIntent = intent
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit := weatherParser(speechText, botLocation, botUnits)
		intentParams = map[string]string{"condition": condition, "is_forecast": is_forecast, "local_datetime": local_datetime, "speakable_location_string": speakable_location_string, "temperature": temperature, "temperature_unit": temperature_unit}
	} else if strings.Contains(intent, "intent_imperative_volumelevel_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, "medium lo") || strings.Contains(speechText, "media lo") || strings.Contains(speechText, "medium bo") || strings.Contains(speechText, "media bo") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_2"
		} else if strings.Contains(speechText, "low") || strings.Contains(speechText, "quiet") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else if strings.Contains(speechText, "medium high") || strings.Contains(speechText, "media high") || strings.Contains(speechText, "medium hide") || strings.Contains(speechText, "media hide") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_4"
		} else if strings.Contains(speechText, "medium") || strings.Contains(speechText, "normal") || strings.Contains(speechText, "regular") || strings.Contains(speechText, "media") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_3"
		} else if strings.Contains(speechText, "high") || strings.Contains(speechText, "loud") || strings.Contains(speechText, "hide") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_5"
		} else if strings.Contains(speechText, "mute") || strings.Contains(speechText, "nothing") || strings.Contains(speechText, "silent") || strings.Contains(speechText, "off") || strings.Contains(speechText, "zero") || strings.Contains(speechText, "meet") {
			// there is no VOLUME_0 :(
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_names_username_extend") {
		var username string
		var nameSplitter string
		isParam = true
		newIntent = "intent_names_username"
		if strings.Contains(speechText, "is") {
			nameSplitter = "is"
		} else if strings.Contains(speechText, "'s") {
			nameSplitter = "'s"
		} else if strings.Contains(speechText, "names") {
			nameSplitter = "names"
		}
		if strings.Contains(speechText, "is") || strings.Contains(speechText, "'s") || strings.Contains(speechText, "names") {
			splitPhrase := strings.SplitAfter(speechText, nameSplitter)
			username = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				username = username + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			if debugLogging == true {
				fmt.Println("Name parsed from speech: " + "`" + username + "`")
			}
			intentParam = "username"
			intentParamValue = username
			intentParams = map[string]string{intentParam: intentParamValue}
		} else {
			if debugLogging == true {
				fmt.Println("No name parsed from speech")
			}
			intentParam = "username"
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	} else if strings.Contains(intent, "intent_clock_settimer_extend") {
		isParam = true
		newIntent = "intent_clock_settimer"
		timerSecs := words2num(speechText)
		if debugLogging == true {
			fmt.Println("Seconds parsed from speech: " + timerSecs)
		}
		intentParam = "timer_duration"
		intentParamValue = timerSecs
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_global_stop_extend") {
		isParam = true
		newIntent = "intent_global_stop"
		intentParam = "what_to_stop"
		intentParamValue = "timer"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_playmessage_extend") {
		var given_name string
		isParam = true
		newIntent = "intent_message_playmessage"
		intentParam = "given_name"
		if strings.Contains(speechText, " for ") {
			splitPhrase := strings.SplitAfter(speechText, " for ")
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_recordmessage_extend") {
		var given_name string
		isParam = true
		newIntent = "intent_message_recordmessage"
		intentParam = "given_name"
		if strings.Contains(speechText, " for ") {
			splitPhrase := strings.SplitAfter(speechText, " for ")
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_blackjack") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "blackjack"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_fistbump") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "fist_bump"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_rollcube") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "roll_cube"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_praise") {
		isParam = false
		newIntent = "intent_imperative_affirmative"
		intentParam = ""
		intentParamValue = ""
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_abuse") {
		isParam = false
		newIntent = "intent_imperative_negative"
		intentParam = ""
		intentParamValue = ""
		intentParams = map[string]string{intentParam: intentParamValue}
	} else {
		newIntent = intent
		intentParam = ""
		intentParamValue = ""
		isParam = false
		intentParams = map[string]string{intentParam: intentParamValue}
	}
	IntentPass(req, newIntent, speechText, intentParams, isParam, justThisBotNum)
}
