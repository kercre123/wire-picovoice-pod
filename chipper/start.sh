#!/bin/bash

if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root. sudo ./start.sh"
  exit 1
fi

if [[ -d ./chipper ]]; then
   cd chipper
fi

UNAME=$(uname -a)

if [[ "${UNAME}" == *"x86_64"* ]]; then
   export GOARCH="amd64"
   echo "amd64 architecture confirmed."
elif [[ "${UNAME}" == *"aarch64"* ]]; then
   export GOARCH="arm64"
   echo "aarch64 architecture confirmed."
elif [[ "${UNAME}" == *"armv7l"* ]]; then
   export GOARCH="arm"
   echo "armv7l architecture confirmed."
else
   echo "Your CPU architecture not supported. This script currently supports x86_64, aarch64, and armv7l."
   exit 1
fi

if [[ $OSTYPE == *"darwin"* ]]; then
  export GOOS="darwin"
else
  export GOOS="linux"
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
