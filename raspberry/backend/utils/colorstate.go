package utils

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
)

type ColorState struct {
	stateSizeX int
	stateSizeY int
	stateSizeZ int
	state []colorful.Color
}

func NewColorState(stateSizeX int, stateSizeY int, stateSizeZ int) ColorState {
	var cs = ColorState {
		stateSizeX: stateSizeX,
		stateSizeY: stateSizeY,
		stateSizeZ: stateSizeZ,
	}

	var size = cs.stateSizeX * cs.stateSizeY * cs.stateSizeZ
	fmt.Printf("colorstate size %d\n", size)
	for i := 0; i < size; i++ {
		cs.state = append(cs.state, colorful.Color{})
	}

	return cs
}

func (s ColorState) fill(color colorful.Color)  {
	for i := 0; i < len(s.state); i++ {
		s.state[i].B = color.B
		s.state[i].G = color.G
		s.state[i].R = color.R
	}
}

func (s ColorState) fillX(x int, color colorful.Color)  {
	for z := 0; z < s.stateSizeZ; z++ {
		for y := 0; y < s.stateSizeY; y++ {
			s.set(x, y, z, color)
		}
	}
}

func (s ColorState) fillY(y int, color colorful.Color)  {
	for z := 0; z < s.stateSizeZ; z++ {
		for x := 0; x < s.stateSizeX; x++ {
			s.set(x, y, z, color)
		}
	}
}

func (s ColorState) fillZ(z int, color colorful.Color)  {
	var zSize = s.stateSizeY * s.stateSizeX
	var start = zSize * z
	var end = zSize * (z+1)
	for i := start; i < end; i++ {
		s.state[i].B = color.B
		s.state[i].G = color.G
		s.state[i].R = color.R
	}
}

func (s ColorState) set(x int, y int, z int, color colorful.Color) {
	s.checkBounds(x, y, z)
	var index = s.xyzToIndex(x,y,z)
	s.state[index].B = color.B
	s.state[index].G = color.G
	s.state[index].R = color.R
}

func (s ColorState) get(x int, y int, z int) colorful.Color {
	s.checkBounds(x, y, z)
	var index = s.xyzToIndex(x,y,z)
	return s.state[index]
}

func (s ColorState) checkBounds(x int, y int, z int) {
	if x >= s.stateSizeX {
		panic(fmt.Sprintf("x > maxX x: %d maxX: %d", x, s.stateSizeX))
	}
	if y >= s.stateSizeY {
		panic(fmt.Sprintf("y > maxY y: %d maxY: %d", y, s.stateSizeY))
	}
	if z >= s.stateSizeZ {
		panic(fmt.Sprintf("z > maxZ z: %d maxZ: %d", z, s.stateSizeZ))
	}
}

func (s ColorState) xyzToIndex(x int, y int, z int) int {
	var xOffset = x
	var yOffset = y*s.stateSizeX
	var zOffset = z*s.stateSizeX*s.stateSizeY
	return xOffset + yOffset + zOffset
}