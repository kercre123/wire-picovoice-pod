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

func InitPicovoice() {
	var picovoiceKey string
	picovoiceKeyOS := os.Getenv("PICOVOICE_APIKEY")
	leopardKeyOS := os.Getenv("LEOPARD_APIKEY")
	if picovoiceInstancesOS == "" {
		fmt.Println("PICOVOICE_INSTANCES is not set, using default value of 5")
		picovoiceInstances = 5
	} else {
		picovoiceInstancesToInt, err := strconv.Atoi(picovoiceInstancesOS)
		picovoiceInstances = picovoiceInstancesToInt
		if err != nil {
			fmt.Println("PICOVOICE_INSTANCES is not a valid integer, using default value of 5")
			picovoiceInstances = 5
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
		picovoiceModeOS = "OnlyLeopard"
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
	fmt.Println("Initializing Picovoice Instances...")
	for i := 0; i < picovoiceInstances; i++ {
		fmt.Println("Initializing Picovoice Instance " + strconv.Itoa(i))
		if picovoiceModeOS == "OnlyLeopard" || picovoiceModeOS == "LeopardAndRhino" {
			leopardSTTArray = append(leopardSTTArray, leopard.Leopard{AccessKey: picovoiceKey})
			leopardSTTArray[i].Init()
		}
		if picovoiceModeOS == "OnlyRhino" || picovoiceModeOS == "LeopardAndRhino" {
			if strings.Contains(os.Getenv("GOARCH"), "arm") {
				rhinoSTIArray = append(rhinoSTIArray, rhino.NewRhino(picovoiceKey, "./armintents.rhn"))
			} else {
				rhinoSTIArray = append(rhinoSTIArray, rhino.NewRhino(picovoiceKey, "./amd64intents.rhn"))
			}
			rhinoSTIArray[i].Init()
		}
	}
}
