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
	FrameCount    int16
	DelayPerFrame int16
	BytesPerLED   int16
	LedCount      int16
	Flags         int32
}

var (
	repeatFlag        int32 = 1 << 0
	shiftCircularFlag int32 = 1 << 1
	interpolateFlag   int32 = 1 << 2
	fillFlag          int32 = 1 << 3 //TODO implement
)

func generateEffectHeader(frameCount int, delay int16, bytesPerLed int, ledCount int, repeat bool, shiftCircular bool, interpolate bool) []byte {
	var flags int32
	if repeat {
		flags = flags | repeatFlag
	}
	if shiftCircular {
		flags = flags | shiftCircularFlag
	}
	if interpolate {
		flags = flags | interpolateFlag
	}
	header := EffectFormatHeader{int16(frameCount), delay, int16(bytesPerLed), int16(ledCount), flags}
	log.Printf("effect header created %v", header)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, header)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}

func min(r uint8, g uint8, b uint8) uint8 {
	x := r
	if x > g {
		x = g
	}
	if x > b {
		x = b
	}
	return x
}

func GenerateColorFadeEffect(bytesPerLed int, ledCount int, hexColor string) ([]byte, error) {
	var steps = 32
	frameCount := steps * 2
	var values = generateEffectHeader(frameCount, defaultDelay*3, bytesPerLed, ledCount, true, false, true)

	var c, err = colorful.Hex("#" + hexColor)

	if err != nil {
		return nil, err
	}
	var inc = 0.8 / float64(steps)

	var h, s, v = c.Hsv()
	s = 1
	v = 0.2

	for j := 0; j < steps; j++ {
		r, g, b := colorful.Hsv(h, s, v).RGB255()
		w := min(r, g, b)
		for i := 0; i < ledCount; i++ {
			values = append(values, g)
			values = append(values, r)
			values = append(values, b)
			values = append(values, w)
		}
		v += inc
	}
	for j := 0; j < steps; j++ {
		r, g, b := colorful.Hsv(h, s, v).RGB255()
		w := min(r, g, b)
		for i := 0; i < ledCount; i++ {
			values = append(values, g)
			values = append(values, r)
			values = append(values, b)
			values = append(values, w)
		}
		v -= inc
	}
	return values, nil
}

func GenerateRainbowFade(bytesPerLed int, ledCount int) []byte {
	var steps = 200
	var values = generateEffectHeader(steps, defaultDelay, bytesPerLed, ledCount, true, false, false)

	var inc = 360.0 / float64(steps)
	var hue = rand.Float64() * 360
	for j := 0; j < steps; j++ {
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		for i := 0; i < ledCount; i++ {
			values = append(values, g)
			values = append(values, r)
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

func GenerateInterpolateTest(bytesPerLed int, ledCount int) []byte {
	var frames = 4
	var values = generateEffectHeader(frames, 1500, bytesPerLed, ledCount, true, false, true)

	for i := 0; i < ledCount; i++ {
		values = append(values, 255)
		values = append(values, 0)
		values = append(values, 0)
		values = append(values, 0)
	}
	for i := 0; i < ledCount; i++ {
		values = append(values, 0)
		values = append(values, 255)
		values = append(values, 0)
		values = append(values, 0)
	}
	for i := 0; i < ledCount; i++ {
		values = append(values, 0)
		values = append(values, 0)
		values = append(values, 255)
		values = append(values, 0)
	}
	for i := 0; i < ledCount; i++ {
		values = append(values, 0)
		values = append(values, 0)
		values = append(values, 0)
		values = append(values, 255)
	}
	return values
}

func GenerateWarmColorFade(bytesPerLed int, ledCount int) []byte {
	var steps = 20
	var values = generateEffectHeader(steps, 1000, bytesPerLed, ledCount, true, false, true)

	colors := colorful.FastWarmPalette(steps)
	for j := 0; j < steps; j++ {
		r, g, b := colors[j].RGB255()
		for i := 0; i < ledCount; i++ {
			values = append(values, g)
			values = append(values, r)
			values = append(values, b)
			values = append(values, 0)
		}
	}
	return values
}

func GenerateRunningRainbow(bytesPerLed int, ledCount int) []byte {
	var values = generateEffectHeader(1, defaultDelay, bytesPerLed, ledCount, true, true, false)

	hueStepSize := 360.0 / float64(ledCount)
	var hue = rand.Float64() * 360
	for i := 0; i < ledCount; i++ {
		r, g, b := colorful.Hsv(hue, 1, 1).RGB255()
		values = append(values, g)
		values = append(values, r)
		values = append(values, b)
		values = append(values, 0)
		hue += hueStepSize
		if hue > 360 {
			hue = 0
		}
	}
	return values
}

func GenerateRunningWarmColors(bytesPerLed int, ledCount int) ([]byte, error) {
	var values = generateEffectHeader(1, defaultDelay, bytesPerLed, ledCount, true, true, false)

	colors := colorful.FastWarmPalette(ledCount/3 + 1)

	for i := 0; i < ledCount; i++ {
		r, g, b := colors[i/3].RGB255()
		w := min(r, g, b)
		values = append(values, g)
		values = append(values, r)
		values = append(values, b)
		values = append(values, w)
	}
	return values, nil
}

func GenerateRunningHappyColors(bytesPerLed int, ledCount int) ([]byte, error) {
	var values = generateEffectHeader(1, defaultDelay, bytesPerLed, ledCount, true, true, false)

	colors := colorful.FastHappyPalette(ledCount/3 + 1)

	for i := 0; i < ledCount; i++ {
		r, g, b := colors[i/3].RGB255()
		values = append(values, g)
		values = append(values, r)
		values = append(values, b)
		values = append(values, 0)
	}
	return values, nil
}
