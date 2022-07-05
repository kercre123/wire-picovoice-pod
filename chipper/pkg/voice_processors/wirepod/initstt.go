package wirepod

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	leopard "github.com/Picovoice/leopard/binding/go"
	rhino "github.com/Picovoice/rhino/binding/go/v2"
)

var leopardSTTArray []leopard.Leopard
var rhinoSTIArray []rhino.Rhino
var picovoiceInstancesOS string = os.Getenv("PICOVOICE_INSTANCES")
var picovoiceModeOS string = os.Getenv("PICOVOICE_MODE")
var picovoiceInstances int
var rhinoSensitivity float32 = 0.05
var rhinoEndpointDurationSec float32 = 0.55

func InitPicovoice() {
	var picovoiceKey string
	picovoiceKeyOS := os.Getenv("PICOVOICE_APIKEY")
	leopardKeyOS := os.Getenv("LEOPARD_APIKEY")
	if picovoiceInstancesOS == "" {
		picovoiceInstances = 3
	} else {
		picovoiceInstancesToInt, err := strconv.Atoi(picovoiceInstancesOS)
		picovoiceInstances = picovoiceInstancesToInt
		if err != nil {
			fmt.Println("PICOVOICE_INSTANCES is not a valid integer, using default value of 3")
			picovoiceInstances = 3
		}
		fmt.Println("Initializing " + strconv.Itoa(picovoiceInstances) + " Picovoice Instances...")
	}
	if picovoiceKeyOS == "" {
		if leopardKeyOS == "" {
			fmt.Println("You must set PICOVOICE_APIKEY to a value.")
			os.Exit(1)
		} else {
			fmt.Println("PICOVOICE_APIKEY is not set, using LEOPARD_APIKEY")
			picovoiceKey = leopardKeyOS
		}
	} else {
		picovoiceKey = picovoiceKeyOS
	}
	if picovoiceModeOS == "" {
		picovoiceModeOS = "LeopardAndRhino"
	} else {
		if picovoiceModeOS != "OnlyLeopard" && picovoiceModeOS != "OnlyRhino" && picovoiceModeOS != "LeopardAndRhino" && picovoiceModeOS != "OlderPi" {
			fmt.Println("PICOVOICE_MODE is not set to a valid value, using default value of OnlyLeopard")
			picovoiceModeOS = "OnlyLeopard"
		}
	}
	if picovoiceModeOS == "OlderPi" {
		picovoiceInstances = 1
		picovoiceModeOS = "OnlyRhino"
	}
	fmt.Println("Picovoice Mode: " + picovoiceModeOS)
	fmt.Println("Initializing " + strconv.Itoa(picovoiceInstances) + " Picovoice Instances...")
	for i := 0; i < picovoiceInstances; i++ {
		fmt.Println("Initializing Picovoice Instance " + strconv.Itoa(i))
		if picovoiceModeOS == "OnlyLeopard" || picovoiceModeOS == "LeopardAndRhino" {
			leopardSTTArray = append(leopardSTTArray, leopard.Leopard{AccessKey: picovoiceKey})
			if i == picovoiceInstances-1 {
				leopardSTTArray[i].Init()
			} else {
				go leopardSTTArray[i].Init()
			}
		}
		if picovoiceModeOS == "OnlyRhino" {
			if strings.Contains(os.Getenv("GOARCH"), "arm") && strings.Contains(os.Getenv("GOOS"), "linux") {
				rhinoSTIArray = append(rhinoSTIArray, rhino.Rhino{AccessKey: picovoiceKey, ContextPath: "./rhn/piintents.rhn", Sensitivity: rhinoSensitivity, EndpointDurationSec: rhinoEndpointDurationSec})
			} else if strings.Contains(os.Getenv("GOOS"), "darwin") {
				rhinoSTIArray = append(rhinoSTIArray, rhino.Rhino{AccessKey: picovoiceKey, ContextPath: "./rhn/darwinintentsnoweather.rhn", Sensitivity: rhinoSensitivity, EndpointDurationSec: rhinoEndpointDurationSec})
			} else {
				rhinoSTIArray = append(rhinoSTIArray, rhino.Rhino{AccessKey: picovoiceKey, ContextPath: "./rhn/amd64intents.rhn", Sensitivity: rhinoSensitivity, EndpointDurationSec: rhinoEndpointDurationSec})
			}
			if i == picovoiceInstances-1 {
				rhinoSTIArray[i].Init()
			} else {
				go rhinoSTIArray[i].Init()
			}
		}
		if picovoiceModeOS == "LeopardAndRhino" {
			if strings.Contains(os.Getenv("GOARCH"), "arm") && strings.Contains(os.Getenv("GOOS"), "linux") {
				rhinoSTIArray = append(rhinoSTIArray, rhino.Rhino{AccessKey: picovoiceKey, ContextPath: "./rhn/piintentsnoweather.rhn", Sensitivity: rhinoSensitivity, EndpointDurationSec: rhinoEndpointDurationSec})
			} else if strings.Contains(os.Getenv("GOOS"), "darwin") {
				rhinoSTIArray = append(rhinoSTIArray, rhino.Rhino{AccessKey: picovoiceKey, ContextPath: "./rhn/darwinintentsnoweather.rhn", Sensitivity: rhinoSensitivity, EndpointDurationSec: rhinoEndpointDurationSec})
			} else {
				rhinoSTIArray = append(rhinoSTIArray, rhino.Rhino{AccessKey: picovoiceKey, ContextPath: "./rhn/amd64intentsnoweather.rhn", Sensitivity: rhinoSensitivity, EndpointDurationSec: rhinoEndpointDurationSec})
			}
			if i == picovoiceInstances-1 {
				rhinoSTIArray[i].Init()
			} else {
				go rhinoSTIArray[i].Init()
			}
		}
	}
	fmt.Println("Picovoice Initialized!")
}
