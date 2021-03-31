package utils

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
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

func StartZScannerSynced(
	ctx context.Context,
	palette []string,
	state ColorState,
	onColorStateUpdate func(),
) {
	initColor, err := colorful.Hex(palette[0])
	black := colorful.Color{}
	for true {
		for i := 0; i < state.stateSizeX; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if err != nil {
				panic(err)
			}
			if i == 0 {
				state.fillX(state.stateSizeX-1, black)
			} else {
				state.fillX(i-1, black)
			}
			state.fillX(i, initColor)
			onColorStateUpdate()
			time.Sleep(time.Millisecond * defaultDelay)
		}
	}
}

func StartRainbowSynced(
	ctx context.Context,
	palette []string,
	state ColorState,
	onColorStateUpdate func(),
) {
	var hue float64 = 0
	for true {
		color := colorful.Hsv(hue, 1, 1)
		select {
		case <-ctx.Done():
			return
		default:
		}
		state.fill(color)
		onColorStateUpdate()
		time.Sleep(time.Millisecond * defaultDelay)
		hue += 1
		if hue > 360 {
			hue = 0
		}
	}
}

func StartZFlowSynced(
	ctx context.Context,
	palette []string,
	state ColorState,
	onColorStateUpdate func(),
) {
	var paletteLen = float64(len(palette))
	var paletteInc = paletteLen / float64(state.stateSizeZ)
	var pos = 0.0

	//at default speed it takes 10 seconds to rotate the palette pos
	var posInc = paletteLen / 500

	fmt.Printf("posinc %f", posInc)

	for true {
		pos += posInc
		select {
		case <-ctx.Done():
			return
		default:
		}
		for i := 0; i < state.stateSizeZ; i++ {
			pos += paletteInc
			var index = int(math.Round(pos))
			if float64(index) >= paletteLen {
				index = 0
				pos -= paletteLen
			}
			color, err := colorful.Hex(palette[index])
			if err != nil {
				panic("could not parse " + palette[index])
			}
			state.fillZ(i, color)
		}
		onColorStateUpdate()
		time.Sleep(time.Millisecond * defaultDelay)
	}
}

func StartRainbowXFlowSynced(
	ctx context.Context,
	palette []string,
	state ColorState,
	onColorStateUpdate func(),
) {
	var hue float64 = 0
	var hueInc float64 = 360.0 / float64(state.stateSizeX)

	for true {
		hue += 3
		select {
		case <-ctx.Done():
			return
		default:
		}
		for i := 0; i < state.stateSizeX; i++ {
			hue += hueInc
			if hue > 360 {
				hue -= 360
			}

			color := colorful.Hsv(hue, 1, 1)
			state.fillX(i, color)
		}
		onColorStateUpdate()
		time.Sleep(time.Millisecond * defaultDelay)
	}
}

func StartRainbowYFlowSynced(
	ctx context.Context,
	palette []string,
	state ColorState,
	onColorStateUpdate func(),
) {
	var hue float64 = 0
	var hueInc float64 = 360.0 / float64(state.stateSizeY)

	for true {
		hue += 3
		select {
		case <-ctx.Done():
			return
		default:
		}
		for i := 0; i < state.stateSizeY; i++ {
			hue += hueInc
			if hue > 360 {
				hue -= 360
			}

			color := colorful.Hsv(hue, 1, 1)
			state.fillY(i, color)
		}
		onColorStateUpdate()
		time.Sleep(time.Millisecond * defaultDelay)
	}
}

func generateWARLSHeader() []byte {
	//1 = protocol
	//2 = time till web interfaces gets activated again
	return []byte{1, 2}
}
