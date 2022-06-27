#!/bin/bash

if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root. sudo ./start.sh"
  exit 1
fi

if [[ -d ./chipper ]]; then
   cd chipper
fi

#if [[ ! -f ./chipper ]]; then
#   if [[ -f ./go.mod ]]; then
#     echo "You need to build chipper first. This can be done with the setup.sh script."
#   else
#     echo "You must be in the chipper directory."
#   fi
#   exit 0
#fi

if [[ ! -f ./source.sh ]]; then
  echo "You need to make a source.sh file. This can be done with the setup.sh script."
  exit 0
fi

source source.sh

if [[ -f ./chipper ]]; then
  ./chipper
else
if [[ $OSTYPE == *"darwin"* ]]; then
  if [[ ! -f ./gotSys ]]; then
    go get -u golang.org/x/sys
    touch gotSys
  fi
  go run cmd/main.go
else
  /usr/local/go/bin/go run cmd/main.go
fi
fi
