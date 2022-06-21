package wirepod

import (
	"encoding/binary"
	"io"
	"log"
	"os"

	opus "github.com/digital-dream-labs/opus-go/opus"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

var processingOne bool = false
var processingTwo bool = false
var processingThree bool = false
var processingFour bool = false

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

func bytesToInt(stream opus.OggStream, data []byte, die bool) []int16 {
	if die == true {
		return nil
	}
	n, err := stream.Decode(data)
	if err != nil {
		log.Println(err)
	}
	return bytesToSamples(n)
}
