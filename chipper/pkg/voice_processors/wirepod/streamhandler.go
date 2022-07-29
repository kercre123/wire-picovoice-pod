package wirepod

import (
	"encoding/binary"
	"fmt"

	opus "github.com/digital-dream-labs/opus-go/opus"
)

func bytesToSamples(buf []byte) []int16 {
	samples := make([]int16, len(buf)/2)
	for i := 0; i < len(buf)/2; i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(buf[i*2:]))
	}
	return samples
}

func bytesToIntLeopard(stream opus.OggStream, data []byte, die bool, isOpus bool) []int16 {
	// detect if data is pcm or opus
	if die {
		return nil
	}
	if isOpus {
		// opus
		n, err := stream.Decode(data)
		if err != nil {
			fmt.Println(err)
		}
		return bytesToSamples(n)
	} else {
		// pcm
		return bytesToSamples(data)
	}
}

func bytesToIntRhino(stream opus.OggStream, data []byte, die bool, isOpus bool) [][]int16 {
	// detect if data is pcm or opus
	if die {
		return nil
	}
	if isOpus {
		// opus
		n, err := stream.Decode(data)
		if err != nil {
			fmt.Println(err)
		}
		nint := bytesToSamples(n)
		// divide nint into chunks of 512 samples
		chunks := make([][]int16, len(nint)/512)
		for i := 0; i < len(nint)/512; i++ {
			chunks[i] = nint[i*512 : (i+1)*512]
		}
		return chunks
	} else {
		// pcm
		nint := bytesToSamples(data)
		// divide nint into chunks of 512 samples
		chunks := make([][]int16, len(nint)/512)
		for i := 0; i < len(nint)/512; i++ {
			chunks[i] = nint[i*512 : (i+1)*512]
		}
		return chunks
	}
}
