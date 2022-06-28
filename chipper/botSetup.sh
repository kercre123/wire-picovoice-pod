#!/bin/bash

# This script will be used to make individual bots known to the server.
# This should hopefully only be temporary while jdocs get implemented.
# Usage: ./botSetup.sh "<esn>" "<vicos-version>" "<location>" "<units>"
# <esn>: example: "0060059b" (quotes included), found in CCIS (double press button on charger, lift the lift up then down). This is used as a unique identifier for your bot.
# <firmware-prefix>: example: "1.8". This can also be found in CCIS. Only include the first two numbers. If you see 1.6.0.3331, put "1.6", etc. Some intents works differently on older versions of VicOS.
# <location>: example: "Des Moines, Iowa". This will be used for the weather command. The Weather API does a good job figuring out the location from a messy input, so you can put "des moines ia" and it should work.
# <units>: example: "F". Has to be F or C. F for fahrenheit, C for celcius. This will be used for the weather command.

if [[ -d chipper ]]; then
    cd chipper
fi

if [[ ! -f ./start.sh ]]; then
    echo "This script must be run in the ./chipper directory."
    exit 1
fi

if [[ ! -n $4 ]]; then
    echo 'Usage: ./botSetup.sh <esn> <firmware-prefix> "<location>" <units>'
    echo 'Example: ./botSetup.sh 0060059b 1.8 "Des Moines, Iowa" F'
    exit 1
fi

ESN=$(echo $1 | tr '[:upper:]' '[:lower:]')
#IFS='.' read -ra VICV <<< ${2}
VICV1="$(cut -d'.' -f1 <<<"$2")"
VICV2="$(cut -d'.' -f2 <<<"$2")"
LOCATION=${3}
UNITS=$(echo $4 | tr '[:lower:]' '[:upper:]')

echo
echo "ESN: ${ESN}"
echo "Firmware prefix: ${VICV1}.${VICV2}"
echo "Location: ${LOCATION}"
echo "Units: ${UNITS}"
echo

if [[ ${UNITS} != "F" && ${UNITS} != "C" ]]; then
    echo "ERROR: Units must be F or C. Exiting."
    echo "Usage: ./botSetup.sh <esn> <firmware-prefix> <location> <units>"
    exit 1
fi

if [[ ! -n ${VICV1} ]]; then
    echo "ERROR: VicOS Version must have at least two numbers. Exiting."
    echo "Example: 1.8"
    echo "Usage: ./botSetup.sh <esn> <firmware-prefix> <location> <units>"
fi

IS_EARLY_OPUS="false"

if [[ ${VICV1} == 1 ]]; then
    if [[ ${VICV2} > 5 ]]; then
        USE_PLAY_SPECIFIC="false"
    elif [[ ${VICV2} < 6 ]]; then
        USE_PLAY_SPECIFIC="true"
    fi
elif [[ ${VICV1} == 0 ]]; then
    USE_PLAY_SPECIFIC="true"
    IS_EARLY_OPUS="true"
elif [[ ${VICV1} == 2 ]]; then
    USE_PLAY_SPECIFIC="false"
else
    USE_PLAY_SPECIFIC="false"
fi

mkdir -p botConfigs
echo -n "{" > botConfigs/${ESN}.json
echo -n '"location": ' >> botConfigs/${ESN}.json
echo -n '"' >> botConfigs/${ESN}.json
echo -n ${LOCATION} >> botConfigs/${ESN}.json
echo -n '", ' >> botConfigs/${ESN}.json
echo -n '"units": ' >> botConfigs/${ESN}.json
echo -n '"' >> botConfigs/${ESN}.json
echo -n ${UNITS} >> botConfigs/${ESN}.json
echo -n '", ' >> botConfigs/${ESN}.json
echo -n '"use_play_specific": ' >> botConfigs/${ESN}.json
echo -n "${USE_PLAY_SPECIFIC}" >> botConfigs/${ESN}.json
echo -n ', ' >> botConfigs/${ESN}.json
echo -n '"is_early_opus": ' >> botConfigs/${ESN}.json
echo -n "${IS_EARLY_OPUS}" >> botConfigs/${ESN}.json
echo "}" >> botConfigs/${ESN}.json

echo "Done! ${ESN} is configured."
