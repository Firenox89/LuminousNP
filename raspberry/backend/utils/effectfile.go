package utils

import (
	"bytes"
	"encoding/binary"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"math/rand"
)

const defaultDelay = 50

type EffectFormatHeader struct {
	SchemaVersion int16
	DelayPerFrame int16
	BytesPerLED   int16
	LedCount      int16
	Flags         int32
}

var (
	repeatFlag        int32 = 1 << 0
	shiftCircularFlag int32 = 1 << 1
)

func generateEffectHeader(delay int16, bytesPerLed int, ledCount int, repeat bool, shiftCircular bool) []byte {
	var flags int32
	if repeat {
		flags = flags | repeatFlag
	}
	if shiftCircular {
		flags = flags | shiftCircularFlag
	}
	header := EffectFormatHeader{1, delay, int16(bytesPerLed), int16(ledCount), flags}
	log.Printf("effect header created %v", header)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, header)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}

func GenerateColorFadeEffect(bytesPerLed int, ledCount int, hexColor string) ([]byte, error) {
	var values = generateEffectHeader(defaultDelay*2, bytesPerLed, ledCount, true, false)
	var c, err = colorful.Hex("#"+hexColor)
	if err != nil {
		return nil, err
	}
	var steps = 32
	var inc = 0.8 / float64(steps)
	var h,s,v = c.Hsv()
	s = 1
	v = 0.2

	for j := 0; j < steps; j++ {
		r, g, b := colorful.Hsv(h, s, v).RGB255()
		for i := 0; i < ledCount; i++ {
			values = append(values, r)
			values = append(values, g)
			values = append(values, b)
			values = append(values, 0)
		}
		v += inc
	}
	for j := 0; j < steps; j++ {
		r, g, b := colorful.Hsv(h, s, v).RGB255()
		for i := 0; i < ledCount; i++ {
			values = append(values, r)
			values = append(values, g)
			values = append(values, b)
			values = append(values, 0)
		}
		v -= inc
	}
	return values, nil
}

func GenerateRainbowFade(bytesPerLed int, ledCount int) []byte {
	var values = generateEffectHeader(defaultDelay, bytesPerLed, ledCount, true, false)

	var steps = 200
	var inc = 360.0 / float64(steps)
	var hue = rand.Float64() * 360
	for j := 0; j < steps; j++ {
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		for i := 0; i < ledCount; i++ {
			values = append(values, r)
			values = append(values, g)
			values = append(values, b)
			values = append(values, 0)
		}
		hue += inc
		if hue > 360 {
			hue = 0
		}
	}
	return values
}

func GenerateRunningRainbow(bytesPerLed int, ledCount int) []byte {
	var values = generateEffectHeader(defaultDelay, bytesPerLed, ledCount, true, true)

	hueStepSize := 360.0 / float64(ledCount)
	var hue = rand.Float64() * 360
	for i := 0; i < ledCount; i++ {
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		values = append(values, r)
		values = append(values, g)
		values = append(values, b)
		values = append(values, 0)
		hue += hueStepSize
		if hue > 360 {
			hue = 0
		}
	}
	return values
}
