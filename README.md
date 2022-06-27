# wire-picovoice-pod

This repo contains a custom Vector escape pod made from [chipper](https://github.com/digital-dream-labs/chipper) and [vector-cloud](https://github.com/digital-dream-labs/vector-cloud).

This repo is a copy of [wire-pod](https://github.com/kercre123/wire-pod) but instead of using Coqui STT, it uses Picovoice Leopard. Leopard is faster, more accurate, and supports more hardware than Coqui, but it is not a totally local solution (processing is done locally, but it uploads usage to a server). The Coqui STT version will still be developed alongside this.

## Program descriptions

`chipper` - Chipper is a program used on Digital Dream Lab's servers which takes in a Vector's voice stream, puts it into a speech-to-text processor, and spits out an intent. This is also likely used on the official escape pod. This repo contains an older tree of chipper which does not have the "intent graph" feature (it caused an error upon every new stream), and it now has a working voice processor.

`vector-cloud` - Vector-cloud is the program which runs on Vector himself which uploads the mic stream to a chipper instance. This repo has an older tree of vector-cloud which also does not have the "intent graph" feature and has been modified to allow for a custom CA cert.

## Configuring, installing, running

NOTE: This only works with OSKR-unlocked, Dev-unlocked, or Whiskey robots running VicOS version 1.4 and above.

### Linux

(This currently only works on Arch or Debian-based Linux)

```
cd ~
git clone https://github.com/kercre123/wire-picovoice-pod.git
cd wire-picovoice-pod
sudo ./setup.sh

# You should be able to just press enter for all of the settings
```

Now install the files created by the script onto the bot:

`sudo ./setup.sh scp <vectorip> <path/to/key>`

Example:

`sudo ./setup.sh scp 192.168.1.150 /home/wire/id_rsa_Vector-R2D2`

If you are on my custom software (WireOS), you do not have to provide an SSH key,

Example:

`sudo ./setup.sh scp 192.168.1.150`

The bot should now be configured to communicate with your server. You do not need to restart the bot to start using voice commands with the new server environment, but you will need to restart him at some point for weather commands to be reliable.

To start chipper, run:

```
cd chipper
sudo ./start.sh
```

### Windows

1. Install WSL (Windows Subsystem for Linux)
	- Open Powershell
	- Run `wsl --install`
	- Reboot the system
	- Run `wsl --install -d Ubuntu-20.04`
	- Open up Ubuntu 20.04 in start menu and configure it like it says.
2. Find IP address
	- Open Powershell
	- Run `ipconfig`
	- Find your computer's IPv4 address and note it somewhere. It usually starts with 10.0. or 192.168.
3. Install wire-pod
	- Follow the Linux instructions from above
	- Enter the IP you got from `ipconfig` earlier instead of the one provided by setup.sh
	- Use the default port and do not enter a different one
4. Setup firewall rules
	- Open Powershell
	- Run `Set-ExecutionPolicy`
	- When it asks, enter `Bypass`
	- Download [this file](https://wire.my.to/wsl-firewall.ps1)
	- Go to your Downloads folder in File Explorer and Right Click -> Run as administrator


After all of that, try a voice command.

## Configure specific bots

`./chipper/botSetup.sh` is there to help configure specific bots. This is not required for operation, and is there for only if you have multiple users using your instance of wire-picovoice-pod and you would like to use different locations/units for weather depending on the bot. It also helps chipper know if a bot is on an older version of VicOS so it can account for that. This is only a stand-in until jdocs and stuff get implemented.

```
Usage: ./chipper/botSetup.sh <esn> <firmware-prefix> "<location>" <units>
Example: ./chipper/botSetup.sh 0060059b 1.8 "Des Moines, Iowa" F
```

## 0.10-era bots

0.10 and below use raw PCM streams rather than the modern Opus streams. This has support for those streams and no special configuration server-side is required for it. However: you will need to get a domain which is the same length as `chipper-dev.api.anki.com`, make sure to run this on a port that is 3 characters long (like the default 443), add true TLS certificates (can be done in ./chipper/source.sh. make sure to include the chain), and run these commands (SSHed into the bot):

```
cd /anki/bin
systemctl stop vic-cloud
cp vic-cloud orig-vic-cloud
sed -i "s/chipper-dev.api.anki.com:443/<domain>:<port>/g" vic-cloud
systemctl start vic-cloud
```

## Status

OS Support:

- Arch
- Debian/Ubuntu/other APT distros
- Windows (WSL only)

Architecture Support:

- amd64/x86_64
- arm64/aarch64

Things wire-picovoice-pod has worked on:

- Raspberry Pi 4B+ 4GB RAM with Raspberry Pi OS
	- Recommended platform, very fast
	- 64-bit only (I think)
- Raspberry Pi 4B+ 4GB RAM with Manjaro 22.04
- Nintendo Switch with L4T Ubuntu
- Desktop with Ryzen 5 3600, 16 GB RAM with Ubuntu 22.04
- Laptop with mobile i7
- Late 2009 iMac with Core 2 Duo
- Android Devices
	- Pixel 4, Note 4, Razer Phone, Oculus Quest 2, OnePlus 7 Pro, Moto G6, Pixel 2
	- If you run into an error when trying to execute start.sh, please open an issue. This is a Picovoice Leopard issue and can be solved by editing the leopard module.
	- [Termux](https://github.com/termux/termux-app) proot-distro: Use Ubuntu, make sure to use a port above 1024 and not the default 443.
	- Linux Deploy: Works stock, just make sure to choose the arch that matches your device in settings. Also use a bigger image size, at least 3 GB.

General notes:

- If you get this error when running chipper, you are using a port that is being taken up by a program already: `panic: runtime error: invalid memory address or nil pointer dereference`
	- Run `./setup.sh` with the 5th and 6th option to change the port, you will need to push files to the bot again.
- If you want to disable logging from the voice processor, edit `./chipper/source.sh` and change `DEBUG_LOGGING` to `false`

Current implemented actions:

- Good robot
- Bad robot
- Change your eye color
- Change your eye color to <color>
	- blue, purple, teal, green, yellow
- How old are you
- Start exploring ("deploring" works better)
- Go home (or "go to your charger")
- Go to sleep
- Good morning
- Good night
- What time is it
- Goodbye
- Happy new year
- Happy holidays
- Hello
- Sign in alexa
- Sign out alexa
- I love you
- Move forward
- Turn left
- Turn right
- Roll your cube
- Pop a wheelie
- Fistbump
- Blackjack (say yes/no instead of hit/stand)
- Yes (affirmative)
- No (negative)
- What's my name
- Take a photo
- Take a photo of me
- What's the weather
	- Requires API setup
	- weatherapi.com is implemented, use the 5th option in `./setup.sh` to set it up
	- To set a default location, use the `botSetup.sh` script in the `./chipper` directory
- What's the weather in <location>
	- Requires API setup
	- weatherapi.com is implemented, use the 5th option in `./setup.sh` to set it up
- Im sorry
- Back up
- Come here
- Volume down
- Be quiet
- Volume up
- Look at me
- Set the volume to <volume>
	- High, medium high, medium, medium low, low
- Shut up
- My name is <name>
- I have a question
	- Requires API setup
	- Houndify is implemented, use the 5th option in `./setup.sh` to set it up
- Set a timer for <time> seconds
- Set a timer for <time> minutes
- Check the timer
- Stop the timer
- Dance
- Pick up the cube
- Fetch the cube
- Find the cube
- Do a trick
- Record a message for <name>
	- Enable `Messaging` feature in Vector's webViz Features tab
- Play a message for <name>
	- Enable `Messaging` feature in Vector's webViz Features tab
- Play keepaway
	- This may only be a feature in 1.5 and lower

## Credits

- [Digital Dream Labs](https://github.com/digital-dream-labs) for saving Vector and for open sourcing chipper which made this possible
- [dietb](https://github.com/dietb) for rewriting chipper and giving tips
- [GitHub Copilot](https://copilot.github.com/) for being awesome