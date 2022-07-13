package wirepod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case r.URL.Path == "/api/add_custom_intent":
		name := r.FormValue("name")
		description := r.FormValue("description")
		utterances := r.FormValue("utterances")
		intent := r.FormValue("intent")
		paramName := r.FormValue("paramname")
		paramValue := r.FormValue("paramvalue")
		exec := r.FormValue("exec")
		if name == "" || description == "" || utterances == "" || intent == "" {
			fmt.Fprintf(w, "missing required field (name, description, utterances, and intent are required)")
			return
		}
		if _, err := os.Stat("./customIntents.json"); err == nil {
			fmt.Println("Found customIntents.json")
			var customIntentJSON intentsStruct
			customIntentJSONFile, err := os.ReadFile("./customIntents.json")
			json.Unmarshal(customIntentJSONFile, &customIntentJSON)
			fmt.Println("Number of custom intents (current): " + strconv.Itoa(len(customIntentJSON)))
			customIntentJSON = append(customIntentJSON, struct {
				Name        string   `json:"name"`
				Description string   `json:"description"`
				Utterances  []string `json:"utterances"`
				Intent      string   `json:"intent"`
				Params      struct {
					ParamName  string `json:"paramname"`
					ParamValue string `json:"paramvalue"`
				} `json:"params"`
				Exec string `json:"exec"`
			}{Name: name, Description: description, Utterances: strings.Split(utterances, ","), Intent: intent, Params: struct {
				ParamName  string `json:"paramname"`
				ParamValue string `json:"paramvalue"`
			}{ParamName: paramName, ParamValue: paramValue}, Exec: exec})
			var logParam string
			var logExec string
			if paramName == "" && paramValue == "" {
				logParam = "No Parameters"
			} else {
				logParam = "Parameter Name: " + paramName + ", Parameter Value: " + paramValue
			}
			if exec == "" {
				logExec = "No program to execute"
			} else {
				logExec = "Program to execute: " + exec
			}
			fmt.Println("New custom intent added. Name: " + name + ", Description: " + description + ", Utterances: " + utterances + ", " + logParam + ", " + logExec)
			fmt.Println("Number of custom intents (newfile): " + strconv.Itoa(len(customIntentJSON)))
			customIntentJSONFile, err = json.Marshal(customIntentJSON)
			if err != nil {
				fmt.Println(err)
			}
			err = ioutil.WriteFile("./customIntents.json", customIntentJSONFile, 0644)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Creating customIntents.json")
			customIntentJSONFile, err := json.Marshal([]struct {
				Name        string   `json:"name"`
				Description string   `json:"description"`
				Utterances  []string `json:"utterances"`
				Intent      string   `json:"intent"`
				Params      struct {
					ParamName  string `json:"paramname"`
					ParamValue string `json:"paramvalue"`
				} `json:"params"`
				Exec string `json:"exec"`
			}{{Name: name, Description: description, Utterances: strings.Split(utterances, ","), Intent: intent, Params: struct {
				ParamName  string `json:"paramname"`
				ParamValue string `json:"paramvalue"`
			}{ParamName: paramName, ParamValue: paramValue}, Exec: exec}})
			if err != nil {
				fmt.Println(err)
			}
			var logParam string
			var logExec string
			if paramName == "" && paramValue == "" {
				logParam = "No Parameters"
			} else {
				logParam = "Parameter Name: " + paramName + ", Parameter Value: " + paramValue
			}
			if exec == "" {
				logExec = "No program to execute"
			} else {
				logExec = "Program to execute: " + exec
			}
			fmt.Println("New custom intent added. Name: " + name + ", Description: " + description + ", Utterances: " + utterances + ", " + logParam + ", " + logExec)
			err = ioutil.WriteFile("./customIntents.json", customIntentJSONFile, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Fprintf(w, "intent added successfully")
		return
	case r.URL.Path == "/api/edit_custom_intent":
		number := r.FormValue("number")
		name := r.FormValue("name")
		description := r.FormValue("description")
		utterances := r.FormValue("utterances")
		intent := r.FormValue("intent")
		paramName := r.FormValue("paramname")
		paramValue := r.FormValue("paramvalue")
		exec := r.FormValue("exec")
		if number == "" {
			fmt.Fprintf(w, "err: a number is required")
			return
		}
		if name == "" && description == "" && utterances == "" && intent == "" && paramName == "" && paramValue == "" && exec == "" {
			fmt.Fprintf(w, "err: an entry must be edited")
			return
		}
		if _, err := os.Stat("./customIntents.json"); err == nil {
			// do nothing
		} else {
			fmt.Fprintf(w, "err: you must create an intent first")
			return
		}
		var customIntentJSON intentsStruct
		customIntentJSONFile, err := os.ReadFile("./customIntents.json")
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(customIntentJSONFile, &customIntentJSON)
		newNumbera, err := strconv.Atoi(number)
		newNumber := newNumbera - 1
		if newNumber > len(customIntentJSON) {
			fmt.Fprintf(w, "err: there are only "+strconv.Itoa(len(customIntentJSON))+" intents")
			return
		}
		fmt.Println(customIntentJSON[newNumber].Name + " custom intent is being edited")
		if name != "" {
			customIntentJSON[newNumber].Name = name
			fmt.Println("Name changed to " + name)
		}
		if description != "" {
			customIntentJSON[newNumber].Description = description
			fmt.Println("Description changed to " + description)
		}
		if utterances != "" {
			customIntentJSON[newNumber].Utterances = strings.Split(utterances, ",")
			fmt.Println("Utterances changed to " + utterances)
		}
		if intent != "" {
			customIntentJSON[newNumber].Intent = intent
			fmt.Println("Intent changed to " + intent)
		}
		if paramName != "" {
			customIntentJSON[newNumber].Params.ParamName = paramName
			fmt.Println("Parameter name changed to " + paramName)
		}
		if paramValue != "" {
			customIntentJSON[newNumber].Params.ParamValue = paramValue
			fmt.Println("Parameter value changed to " + paramValue)
		}
		if exec != "" {
			customIntentJSON[newNumber].Exec = exec
			fmt.Println("Program to execute changed to " + exec)
		}
		newCustomIntentJSONFile, err := json.Marshal(customIntentJSON)
		err = ioutil.WriteFile("./customIntents.json", newCustomIntentJSONFile, 0644)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, "intent edited successfully")
		return
	case r.URL.Path == "/api/get_custom_intents_json":
		if _, err := os.Stat("./customIntents.json"); err == nil {
			// do nothing
		} else {
			fmt.Fprintf(w, "err: you must create an intent first")
			return
		}
		customIntentJSONFile, err := ioutil.ReadFile("./customIntents.json")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(customIntentJSONFile))
		return
	case r.URL.Path == "/api/remove_custom_intent":
		number := r.FormValue("number")
		if number == "" {
			fmt.Fprintf(w, "err: a number is required")
			return
		}
		if _, err := os.Stat("./customIntents.json"); err == nil {
			// do nothing
		} else {
			fmt.Fprintf(w, "err: you must create an intent first")
			return
		}
		var customIntentJSON intentsStruct
		customIntentJSONFile, err := os.ReadFile("./customIntents.json")
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(customIntentJSONFile, &customIntentJSON)
		newNumbera, err := strconv.Atoi(number)
		newNumber := newNumbera - 1
		if newNumber > len(customIntentJSON) {
			fmt.Fprintf(w, "err: there are only "+strconv.Itoa(len(customIntentJSON))+" intents")
			return
		}
		fmt.Println(customIntentJSON[newNumber].Name + " custom intent is being removed")
		customIntentJSON = append(customIntentJSON[:newNumber], customIntentJSON[newNumber+1:]...)
		newCustomIntentJSONFile, err := json.Marshal(customIntentJSON)
		err = ioutil.WriteFile("./customIntents.json", newCustomIntentJSONFile, 0644)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, "intent removed successfully")
		return
	case r.URL.Path == "/api/add_bot":
		botESN := r.FormValue("esn")
		botLocation := r.FormValue("location")
		botUnits := r.FormValue("units")
		botFirmwarePrefix := r.FormValue("firmwareprefix")
		var is_early_opus bool
		var use_play_specific bool
		if botESN == "" || botLocation == "" || botUnits == "" || botFirmwarePrefix == "" {
			fmt.Fprintf(w, "err: all fields are required")
			return
		}
		firmwareSplit := strings.Split(botFirmwarePrefix, ".")
		if len(firmwareSplit) != 2 {
			fmt.Fprintf(w, "err: firmware prefix must be in the format: 1.5")
			return
		}
		if botUnits != "F" && botUnits != "C" {
			fmt.Fprintf(w, "err: units must be either F or C")
			return
		}
		firmware1, err := strconv.Atoi(firmwareSplit[0])
		firmware2, err := strconv.Atoi(firmwareSplit[1])
		if err != nil {
			fmt.Fprintf(w, "err: firmware prefix must be in the format: 1.5")
			return
		}
		if firmware1 >= 1 && firmware2 < 6 {
			is_early_opus = false
			use_play_specific = true
		} else if firmware1 >= 1 && firmware2 >= 6 {
			is_early_opus = false
			use_play_specific = false
		} else if firmware1 == 0 {
			is_early_opus = true
			use_play_specific = true
		} else {
			fmt.Fprintf(w, "err: firmware prefix must be in the format: 1.5")
			return
		}
		type botConfigStruct []struct {
			Esn             string `json:"esn"`
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}
		var botConfig botConfigStruct
		if _, err := os.Stat("./botConfig.json"); err == nil {
			// read botConfig.json and append to it with the form information
			botConfigFile, err := ioutil.ReadFile("./botConfig.json")
			if err != nil {
				fmt.Println(err)
			}
			json.Unmarshal(botConfigFile, &botConfig)
			botConfig = append(botConfig, struct {
				Esn             string `json:"esn"`
				Location        string `json:"location"`
				Units           string `json:"units"`
				UsePlaySpecific bool   `json:"use_play_specific"`
				IsEarlyOpus     bool   `json:"is_early_opus"`
			}{Esn: botESN, Location: botLocation, Units: botUnits, UsePlaySpecific: use_play_specific, IsEarlyOpus: is_early_opus})
			newBotConfigJSONFile, err := json.Marshal(botConfig)
			err = ioutil.WriteFile("./botConfig.json", newBotConfigJSONFile, 0644)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			botConfig = append(botConfig, struct {
				Esn             string `json:"esn"`
				Location        string `json:"location"`
				Units           string `json:"units"`
				UsePlaySpecific bool   `json:"use_play_specific"`
				IsEarlyOpus     bool   `json:"is_early_opus"`
			}{Esn: botESN, Location: botLocation, Units: botUnits, UsePlaySpecific: use_play_specific, IsEarlyOpus: is_early_opus})
			newBotConfigJSONFile, err := json.Marshal(botConfig)
			err = ioutil.WriteFile("./botConfig.json", newBotConfigJSONFile, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Fprintf(w, "bot added successfully")
		return
	case r.URL.Path == "/api/remove_bot":
		number := r.FormValue("number")
		if _, err := os.Stat("./botConfig.json"); err == nil {
			// do nothing
		} else {
			fmt.Fprintf(w, "err: you must create a bot first")
			return
		}
		type botConfigStruct []struct {
			Esn             string `json:"esn"`
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}
		var botConfigJSON botConfigStruct
		botConfigJSONFile, err := os.ReadFile("./botConfig.json")
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(botConfigJSONFile, &botConfigJSON)
		newNumbera, err := strconv.Atoi(number)
		newNumber := newNumbera - 1
		if newNumber > len(botConfigJSON) {
			fmt.Fprintf(w, "err: there are only "+strconv.Itoa(len(botConfigJSON))+" bots")
			return
		}
		fmt.Println(botConfigJSON[newNumber].Esn + " bot is being removed")
		botConfigJSON = append(botConfigJSON[:newNumber], botConfigJSON[newNumber+1:]...)
		newBotConfigJSONFile, err := json.Marshal(botConfigJSON)
		err = ioutil.WriteFile("./botConfig.json", newBotConfigJSONFile, 0644)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, "bot removed successfully")
		return
	case r.URL.Path == "/api/edit_bot":
		number := r.FormValue("number")
		botESN := r.FormValue("esn")
		botLocation := r.FormValue("location")
		botUnits := r.FormValue("units")
		botFirmwarePrefix := r.FormValue("firmwareprefix")
		if botESN == "" || botLocation == "" || botUnits == "" || botFirmwarePrefix == "" {
			fmt.Fprintf(w, "err: all fields are required")
			return
		}
		firmwareSplit := strings.Split(botFirmwarePrefix, ".")
		if len(firmwareSplit) != 2 {
			fmt.Fprintf(w, "err: firmware prefix must be in the format: 1.5")
			return
		}
		if botUnits != "F" && botUnits != "C" {
			fmt.Fprintf(w, "err: units must be either F or C")
			return
		}
		var is_early_opus bool
		var use_play_specific bool
		firmware1, err := strconv.Atoi(firmwareSplit[0])
		firmware2, err := strconv.Atoi(firmwareSplit[1])
		if err != nil {
			fmt.Fprintf(w, "err: firmware prefix must be in the format: 1.5")
			return
		}
		if firmware1 >= 1 && firmware2 < 6 {
			is_early_opus = false
			use_play_specific = true
		} else if firmware1 >= 1 && firmware2 >= 6 {
			is_early_opus = false
			use_play_specific = false
		} else if firmware1 == 0 {
			is_early_opus = true
			use_play_specific = true
		} else {
			fmt.Fprintf(w, "err: firmware prefix must be in the format: 1.5")
			return
		}
		type botConfigStruct []struct {
			Esn             string `json:"esn"`
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}
		var botConfig botConfigStruct
		if _, err := os.Stat("./botConfig.json"); err == nil {
			// read botConfig.json and append to it with the form information
			botConfigFile, err := ioutil.ReadFile("./botConfig.json")
			if err != nil {
				fmt.Println(err)
			}
			json.Unmarshal(botConfigFile, &botConfig)
			newNumbera, err := strconv.Atoi(number)
			newNumber := newNumbera - 1
			botConfig[newNumber].Esn = botESN
			botConfig[newNumber].Location = botLocation
			botConfig[newNumber].Units = botUnits
			botConfig[newNumber].UsePlaySpecific = use_play_specific
			botConfig[newNumber].IsEarlyOpus = is_early_opus
			newBotConfigJSONFile, err := json.Marshal(botConfig)
			err = ioutil.WriteFile("./botConfig.json", newBotConfigJSONFile, 0644)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Fprintln(w, "err: you must create a bot first")
			return
		}
		fmt.Fprintf(w, "bot edited successfully")
		return
	case r.URL.Path == "/api/get_bot_json":
		if _, err := os.Stat("./botConfig.json"); err == nil {
			// do nothing
		} else {
			fmt.Fprintf(w, "err: you must add a bot first")
			return
		}
		botConfigJSONFile, err := ioutil.ReadFile("./botConfig.json")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(botConfigJSONFile))
		return
	}
}

func StartWebServer() {
	http.HandleFunc("/api/", apiHandler)
	fileServer := http.FileServer(http.Dir("./webroot"))
	http.Handle("/", fileServer)

	fmt.Printf("Starting server at port 8080 (http://localhost:8080)\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
