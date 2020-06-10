package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

const defaultDelay = 250

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
	testFlag1         int32 = 1 << 2
	testFlag2         int32 = 1 << 3
	testFlag3         int32 = 1 << 4
)

func generateEffectHeader(delay int16, bytesPerLed int16, ledCount int16, repeat bool, shiftCircular bool) []byte {
	var flags int32
	if repeat {
		flags = flags | repeatFlag
	}
	if shiftCircular {
		flags = flags | shiftCircularFlag
	}
	header := EffectFormatHeader{1, delay, bytesPerLed, ledCount, flags}
	log.Printf("effect header created %v", header)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, header)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}

func GenerateColorSwitchEffect(bytesPerLed int16, ledCount int16) []byte {
	var values = generateEffectHeader(defaultDelay, bytesPerLed, ledCount, true, false)

	for j := 0; j < 10; j++ {
		for i := 0; i < int(ledCount); i++ {
			values = append(values, 255)
			values = append(values, 0)
			values = append(values, 0)
			values = append(values, 0)
		}
		for i := 0; i < int(ledCount); i++ {
			values = append(values, 0)
			values = append(values, 255)
			values = append(values, 0)
			values = append(values, 0)
		}
		for i := 0; i < int(ledCount); i++ {
			values = append(values, 0)
			values = append(values, 0)
			values = append(values, 255)
			values = append(values, 0)
		}
		for i := 0; i < int(ledCount); i++ {
			values = append(values, 0)
			values = append(values, 0)
			values = append(values, 0)
			values = append(values, 255)
		}
	}
	return values
}

func GenerateRainbowFade(bytesPerLed int16, ledCount int16) []byte {
	var values = generateEffectHeader(defaultDelay, bytesPerLed, ledCount, true, false)

	hsl := HSL{S: 1, L: 0.5}
	for j := 0; j < 50; j++ {
		rgb := hsl.ToRGB()
		for i := 0; i < int(ledCount); i++ {
			values = append(values, byte(rgb.R*255))
			values = append(values, byte(rgb.G*255))
			values = append(values, byte(rgb.B*255))
			values = append(values, 0)
		}
		hsl.H = hsl.H + 0.002
	}
	return values
}

func GenerateRunningRainbow(bytesPerLed int16, ledCount int16) []byte {
	var values = generateEffectHeader(defaultDelay, bytesPerLed, ledCount, true, true)

	hsl := HSL{S: 1, L: 0.5}
	hueStepSize := 1.0 / float64(ledCount)
	for i := 0; i < int(ledCount); i++ {
		rgb := hsl.ToRGB()
		values = append(values, byte(rgb.R*255))
		values = append(values, byte(rgb.G*255))
		values = append(values, byte(rgb.B*255))
		values = append(values, 0)
		hsl.H = hsl.H + hueStepSize
	}
	return values
}
