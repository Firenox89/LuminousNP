package utils

import (
	"github.com/lucasb-eyer/go-colorful"
	"math"
)

type ShowCaseType int

const (
	ShowcaseType1 ShowCaseType = 1 << iota
	ShowcaseType2
)

type XYZPos struct {
	x int
	y int
	z int
}
/*
    7 Y-axis
   /
  /
0+-----------> X-axis
 |
 |
 |
 V Z-axis
*/
type NodeMappings struct {
	Type        ShowCaseType
	CountX      int
	CountY      int
	CountZ      int
	NodeIdCount int
	NodeGridMapping []XYZPos
}

var type1 = NodeMappings{
	Type:        ShowcaseType1,
	CountX:      19,
	CountY:      16,
	CountZ:      4,
	NodeIdCount: 188,
	NodeGridMapping: buildNodeGridMapType1(),
}

var type2 = NodeMappings{
	Type:        ShowcaseType2,
	CountX:      12,
	CountY:      10,
	CountZ:      44,
	NodeIdCount: 200,
	NodeGridMapping: buildNodeGridMapType2(),
}

var levelOffsetType1 = 47
var levelOffsetType2 = 40

func GetMaxXYZ() (int, int, int) {
	return type1.CountX, type1.CountY, type2.CountZ
}

func MapColorStateToNodes(t ShowCaseType, cs ColorState) []colorful.Color {
	switch t {
	case ShowcaseType1:
		return mapType1(cs)
	case ShowcaseType2:
		return mapType2(cs)
	}
	panic("Unknown type")
}

func buildNodeGridMapType1() []XYZPos {
	var xyzMap = make([]XYZPos, 188)

	var nodesZ = buildNodesZType1()
	var nodesY = buildNodesYType1()
	var nodesX = buildNodesXType1()

	for iz := 0; iz < 4; iz++ {
		var zPos = 43 - iz * 11
		var zNodes = nodesZ[iz]
		for yPos := 0; yPos < 16; yPos++ {
			var yNodes = nodesY[yPos]
			for xPos := 0; xPos < 19; xPos++ {
				var xNodes = nodesX[xPos]
				for i := 0; i <len(xNodes); i++ {
					var nodeId = xNodes[i]
					if contains(yNodes, nodeId) && contains(zNodes, nodeId) {
						xyzMap[nodeId] = XYZPos{x: xPos, y: yPos, z: zPos}
					}
				}
			}
		}
	}

	return xyzMap
}

func buildNodeGridMapType2() []XYZPos {
	var xyzMap = make([]XYZPos, 200)

	var nodesZ = buildNodesZType2()
	var nodesY = buildNodesYType2()
	var nodesX = buildNodesXType2()

	var yMultiplier = 16.0 / 10.0
	var xMultiplier = 19.0 / 12.0

	for zPos := 0; zPos < 44; zPos++ {
		var zNodes = nodesZ[zPos]
		for yi := 0; yi < 10; yi++ {
			var yPos = int(math.Round(float64(yi) * yMultiplier))
			var yNodes = nodesY[yi]
			for xi := 0; xi < 19; xi++ {
				var xPos = int(math.Round(float64(xi) * xMultiplier))
				var xNodes = nodesX[xi]
				for i := 0; i <len(xNodes); i++ {
					var nodeId = xNodes[i]
					if contains(yNodes, nodeId) && contains(zNodes, nodeId) {
						xyzMap[nodeId] = XYZPos{x: xPos, y: yPos, z: zPos}
					}
				}
			}
		}
	}

	return xyzMap
}

func contains(nodes []int, id int) bool {
	for i := 0; i < len(nodes); i++ {
		if nodes[i] == id {
			return true
		}
	}
	return false
}

func mapType1(cs ColorState) []colorful.Color {
	var colors = make([]colorful.Color, len(type1.NodeGridMapping))
	for i, pos := range type1.NodeGridMapping {
		colors[i] = cs.get(pos.x, pos.y, pos.z)
	}

	return colors
}

func mapType2(cs ColorState) []colorful.Color {
	var colors = make([]colorful.Color, len(type2.NodeGridMapping))
	for i, pos := range type2.NodeGridMapping {
		colors[i] = cs.get(pos.x, pos.y, pos.z)
	}
	return colors
}

func appendRange(ids []int, start int, end int) []int {
	for i := start; i < end; i++ {
		ids = append(ids, i)
	}
	return ids
}

func buildNodesXType1() map[int][]int {
	var nodes = make(map[int][]int, 19)

	var nodeIds []int

	for i := 0; i < 4; i++ {
		if i%2 == 0 {
			nodeIds = appendRange(nodeIds, levelOffsetType1*i+32, levelOffsetType1*i+47)
		} else {
			nodeIds = appendRange(nodeIds, levelOffsetType1*i, levelOffsetType1*i+15)
		}
	}
	nodes[0] = nodeIds

	for i := 1; i < 18; i++ {
		nodes[i] = []int{32 - i, 14 + i + levelOffsetType1, 32 - i + levelOffsetType1*2, 14 - + i + levelOffsetType1*3}
	}

	nodeIds = []int{}
	for i := 0; i < 4; i++ {
		if i%2 == 1 {
			nodeIds = appendRange(nodeIds, levelOffsetType1*i+32, levelOffsetType1*i+47)
		} else {
			nodeIds = appendRange(nodeIds, levelOffsetType1*i, levelOffsetType1*i+15)
		}
	}
	nodes[18] = nodeIds

	return nodes
}

func buildNodesYType1() map[int][]int {
	var nodes = make(map[int][]int, 16)

	var nodeIds []int

	for i := 0; i < 15; i++ {
		nodeIds = []int{}
		for j := 0; j < 4; j++ {
			nodeIds = append(nodeIds, i+levelOffsetType1*j, levelOffsetType1-1-i+levelOffsetType1*j)
		}
		nodes[i] = nodeIds
	}

	nodeIds = []int{}
	for i := 0; i < 4; i++ {
		nodeIds = appendRange(nodeIds, levelOffsetType1*i+15, levelOffsetType1*i+32)
	}
	nodes[15] = nodeIds

	return nodes
}

func buildNodesZType1() map[int][]int {
	var nodes = make(map[int][]int, 4)

	for i := 0; i < 4; i++ {
		nodes[i] = appendRange([]int{}, levelOffsetType1*i, levelOffsetType1*(i+1))
	}

	return nodes
}

func buildNodesXType2() map[int][]int {
	var nodes = make(map[int][]int, 12)

	//add the left bars
	nodes[0] = appendRange([]int{}, 160, 200)

	var nodeIds []int

	for i := 0; i < 4; i++ {
		nodeIds = appendRange(nodeIds, levelOffsetType2*i+18, levelOffsetType2*i+25)
	}
	nodes[1] = nodeIds

	for j := 0; j < 8; j++ {
		nodeIds = []int{}
		for i := 0; i < 4; i++ {
			nodeIds = append(nodeIds, levelOffsetType2*i+17-j, levelOffsetType2*i+25+j)
		}
		nodes[j+2] = nodeIds
	}

	nodeIds = []int{}
	for i := 0; i < 4; i++ {
		nodeIds = appendRange(nodeIds, levelOffsetType2*i+33, levelOffsetType2*i+40)
	}
	nodes[10] = nodeIds

	//add the right bars
	nodeIds = []int{}
	for i := 0; i < 4; i++ {
		nodeIds = appendRange(nodeIds, levelOffsetType2*i, levelOffsetType2*i+10)
	}
	nodes[11] = nodeIds

	return nodes
}

func buildNodesYType2() map[int][]int {
	var nodes = make(map[int][]int, 10)

	var nodeIds []int

	nodeIds = appendRange(nodeIds, 0, 10)
	nodeIds = appendRange(nodeIds, 40, 50)
	nodeIds = appendRange(nodeIds, 80, 90)
	nodeIds = appendRange(nodeIds, 120, 130)
	nodeIds = appendRange(nodeIds, 160, 200)

	nodes[0] = nodeIds

	nodeIds = []int{}

	for i := 0; i < 4; i++ {
		nodeIds = appendRange(nodeIds, 10+i*levelOffsetType2, 18+i*levelOffsetType2)
	}

	nodes[1] = nodeIds

	for i := 0; i < 7; i++ {
		nodeIds = []int{}
		for i2 := 0; i2 < 4; i2++ {
			nodeIds = append(nodeIds, 18+i+levelOffsetType2*i2, 39-i+levelOffsetType2*i2)
		}
		nodes[i+2] = nodeIds
	}

	nodeIds = []int{}

	for i := 0; i < 4; i++ {
		nodeIds = appendRange(nodeIds, 25+i*levelOffsetType2, 33+i*levelOffsetType2)
	}

	nodes[9] = nodeIds
	return nodes
}

func buildNodesZType2() map[int][]int {
	var nodes = make(map[int][]int, 44)

	for i := 0; i < 10; i++ {
		nodes[i] = []int{i, 199 - i}
	}

	nodes[10] = appendRange([]int{}, 10, 40)

	for i := 0; i < 10; i++ {
		nodes[i+11] = []int{i + 40, 199 - i - 10}
	}

	nodes[21] = appendRange([]int{}, 50, 80)

	for i := 0; i < 10; i++ {
		nodes[i+22] = []int{i + 80, 199 - i - 20}
	}

	nodes[32] = appendRange([]int{}, 90, 120)
	for i := 0; i < 10; i++ {
		nodes[i+33] = []int{i + 120, 199 - i - 30}
	}

	nodes[43] = appendRange([]int{}, 130, 160)

	return nodes
}
