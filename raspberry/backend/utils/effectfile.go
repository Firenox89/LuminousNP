package utils

import (
	"bytes"
	"encoding/binary"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"math"
	"strings"
	"time"
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
	colorValueRamp = []string{"00", "22", "44", "66", "88", "AA", "CC", "FF"}

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

func ApplyBrightnessToColorHex(color string, brightness int) string {
	var c, _ = colorful.Hex(color)
	h, s, v := c.Hsv()
	newV := math.Min(float64(brightness)/100, v)
	return strings.TrimPrefix(colorful.Hsv(h, s, newV).Hex(), "#")
}

func GenerateFadeFromPalette(bytesPerLed int, ledCount int, palette []string, brightness int) []byte {
	var steps = 100
	var values = generateEffectHeader(steps, 500, bytesPerLed, ledCount, true, false, true)

	stepSize := float64(len(palette)) / float64(ledCount)
	for j := 0; j < steps; j++ {
		color, err := colorful.Hex("#" + ApplyBrightnessToColorHex(palette[int(float64(j)*stepSize)], brightness))
		if err != nil {
			panic(err.Error())
		}
		r, g, b := color.RGB255()
		for i := 0; i < ledCount; i++ {
			values = append(values, r)
			values = append(values, g)
			values = append(values, b)
			values = append(values, 0)
		}
	}
	return values
}

func GenerateRotationFromPalette(bytesPerLed int, ledCount int, palette []string, brightness int) []byte {
	var values = generateEffectHeader(1, defaultDelay, bytesPerLed, ledCount, true, true, false)

	stepSize := float64(len(palette)) / float64(ledCount)
	for i := 0; i < ledCount; i++ {
		color, err := colorful.Hex("#" + ApplyBrightnessToColorHex(palette[int(float64(i)*stepSize)], brightness))
		if err != nil {
			panic(err.Error())
		}
		r, g, b := color.RGB255()
		values = append(values, r)
		values = append(values, g)
		values = append(values, b)
		values = append(values, 0)
	}
	return values
}

func BuildRainbowPalette() []string {
	var hue = 0.0
	var values []string
	for j := 0; j < 360; j++ {
		hex := colorful.Hsv(hue, 1, 1).Hex()
		values = append(values, hex)
		hue += 1
		if hue > 360 {
			hue = 0
		}
	}
	return values
}

func BuildHappyPalette() []string {
	colors := colorful.FastHappyPalette(100)

	var values []string
	for j := 0; j < 100; j++ {
		hex := colors[j].Hex()
		values = append(values, hex)
	}
	return values
}

func BuildRedToGreenRamp() []string {
	var result []string
	for i := 0; i < len(colorValueRamp); i++ {
		result = append(result, "#FF"+colorValueRamp[i]+"00")
	}
	return result
}

func BuildRedToBlueRamp() []string {
	var result []string
	for i := 0; i < len(colorValueRamp); i++ {
		result = append(result, "#FF00"+colorValueRamp[i])
	}
	return result
}

func BuildBlueToGreenRamp() []string {
	var result []string
	for i := 0; i < len(colorValueRamp); i++ {
		result = append(result, "#00"+colorValueRamp[i]+"FF")
	}
	return result
}

func BuildBlueToRedRamp() []string {
	var result []string
	for i := 0; i < len(colorValueRamp); i++ {
		result = append(result, "#"+colorValueRamp[i]+"00FF")
	}
	return result
}

func BuildGreenToRedRamp() []string {
	var result []string
	for i := 0; i < len(colorValueRamp); i++ {
		result = append(result, "#"+colorValueRamp[i]+"FF00")
	}
	return result
}

func BuildGreenToBlueRamp() []string {
	var result []string
	for i := 0; i < len(colorValueRamp); i++ {
		result = append(result, "#00FF"+colorValueRamp[i])
	}
	return result
}

func OddlyInsertBlack(colors []string) []string {
	var result []string
	for i := 0; i < len(colors); i++ {
		result = append(result, colors[i], "#000000")
	}
	return result
}

func RevertLoop(colors []string) []string {
	var result []string
	for i := 0; i < len(colors); i++ {
		result = append(result, colors[i])
	}
	for i := len(colors) - 1; i > 1; i-- {
		result = append(result, colors[i-1])
	}
	return result
}

func StartZScanner(showcase NodeMappings, send func(data []byte, ip string) error) error {
	var litLeds []int
	println("Start Loop\n")
	for true {

		for i := 0; i < showcase.CountZ; i++ {
			var data = generateWARLSHeader()
			//turn off the leds from last step
			for _, lit := range litLeds {
				data = append(data, byte(lit), 0, 0, 0)
			}
			litLeds = showcase.NodesZ[i]

			for _, lit := range litLeds {
				data = append(data, byte(lit), 255, 0, 0)
			}
			err := send(data, "192.168.178.61")
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * 100)
		}

		for i := 0; i < showcase.CountX; i++ {
			var data = generateWARLSHeader()
			//turn off the leds from last step
			for _, lit := range litLeds {
				data = append(data, byte(lit), 0, 0, 0)
			}
			litLeds = showcase.NodesX[i]

			for _, lit := range litLeds {
				data = append(data, byte(lit), 255, 0, 0)
			}
			err := send(data, "192.168.178.61")
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * 100)
		}
		for i := 0; i < showcase.CountY; i++ {
			var data = generateWARLSHeader()
			//turn off the leds from last step
			for _, lit := range litLeds {
				data = append(data, byte(lit), 0, 0, 0)
			}
			litLeds = showcase.NodesY[i]

			for _, lit := range litLeds {
				data = append(data, byte(lit), 255, 0, 0)
			}
			err := send(data, "192.168.178.61")
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	return nil
}

func generateWARLSHeader() []byte {
	//1 = protocol
	//2 = time till web interfaces gets activated again
	return []byte{1, 2}
}
