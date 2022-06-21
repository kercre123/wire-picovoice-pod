package wirepod

import (
	"encoding/binary"
	"io"
	"log"
	"os"

	leopard "github.com/Picovoice/leopard/binding/go"
	opus "github.com/digital-dream-labs/opus-go/opus"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

var leopardSTT leopard.Leopard

var processingOne bool = false
var processingTwo bool = false
var processingThree bool = false
var processingFour bool = false

func InitLeopard() {
	log.Println("Initializing Leopard")
	leopardSTT = leopard.Leopard{AccessKey: "access key"}
	err := leopardSTT.Init()
	if err != nil {
		log.Println(err)
	}
	defer leopardSTT.Delete()
	log.Println("Initialized Leopard")
}

func pcmToWav(pcmFile string, wavFile string) {
	in, err := os.Open(pcmFile)
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create(wavFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	e := wav.NewEncoder(out, 16000, 16, 1, 1)
	audioBuf, err := newAudioIntBuffer(in)
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Write(audioBuf); err != nil {
		log.Fatal(err)
	}
	if err := e.Close(); err != nil {
		log.Fatal(err)
	}
}

func newAudioIntBuffer(r io.Reader) (*audio.IntBuffer, error) {
	buf := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  16000,
		},
	}
	for {
		var sample int16
		err := binary.Read(r, binary.LittleEndian, &sample)
		switch {
		case err == io.EOF:
			return &buf, nil
		case err != nil:
			return nil, err
		}
		buf.Data = append(buf.Data, int(sample))
	}
}

func bytesToSamples(buf []byte) []int16 {
	samples := make([]int16, len(buf)/2)
	for i := 0; i < len(buf)/2; i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(buf[i*2:]))
	}
	return samples
}

func bytesToInt(stream opus.OggStream, data []byte, numBot int, voiceTimer int, die bool) {
	if die == true {
		return
	}
	n, err := stream.Decode(data)
	if voiceTimer == 1 {
		log.Println("Starting transcription")
		if processingOne == false {
			processingOne = true
			processOneData := bytesToSamples(n)
			transcription1, err := leopardSTT.Process(processOneData)
			if err != nil {
				log.Println(err)
			}
			log.Println("1: " + transcription1)
		}
	}
	if voiceTimer == 2 {
		if processingTwo == false {
			processingTwo = true
			processTwoData := bytesToSamples(n)
			transcription2, err := leopardSTT.Process(processTwoData)
			if err != nil {
				log.Println(err)
			}
			log.Println("2: " + transcription2)
		}
	}
	if voiceTimer == 3 {
		if processingThree == false {
			processingThree = true
			processThreeData := bytesToSamples(n)
			transcription3, err := leopardSTT.Process(processThreeData)
			if err != nil {
				log.Println(err)
			}
			log.Println("3: " + transcription3)
		}
	}
	if voiceTimer == 4 {
		if processingFour == false {
			processingFour = true
			processFourData := bytesToSamples(n)
			transcription4, err := leopardSTT.Process(processFourData)
			if err != nil {
				log.Println(err)
			}
			log.Println("4: " + transcription4)
		}
	}
	if err != nil {
		log.Println(err)
	}
}
