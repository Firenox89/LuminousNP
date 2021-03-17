package utils

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
	CountX      int
	CountY      int
	CountZ      int
	NodeIdCount int
	NodesX      map[int][]int
	NodesY      map[int][]int
	NodesZ      map[int][]int
}

var ShowcaseType1 = NodeMappings{
	CountX:      19,
	CountY:      16,
	CountZ:      4,
	NodeIdCount: 188,
	NodesX:      buildNodesXType1(),
	NodesY:      buildNodesYType1(),
	NodesZ:      buildNodesZType1(),
}

var levelOffsetType1 = 47

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

var ShowcaseType2 = NodeMappings{
	CountX:      12,
	CountY:      10,
	CountZ:      44,
	NodeIdCount: 200,
	NodesX:      buildNodesXType2(),
	NodesY:      buildNodesYType2(),
	NodesZ:      buildNodesZType2(),
}

var levelOffsetType2 = 40

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
			nodeIds = append(nodeIds, 18 + i + levelOffsetType2*i2, 39 - i + levelOffsetType2*i2 )
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
		nodes[i] = []int{i, 199-i}
	}

	nodes[10] = appendRange([]int{}, 10, 40)

	for i := 0; i < 10; i++ {
		nodes[i+11] = []int{i+40, 199-i-10}
	}

	nodes[21] = appendRange([]int{}, 50, 80)

	for i := 0; i < 10; i++ {
		nodes[i+22] = []int{i+80, 199-i-20}
	}

	nodes[32] = appendRange([]int{}, 90, 120)
	for i := 0; i < 10; i++ {
		nodes[i+33] = []int{i+120, 199-i-30}
	}

	nodes[43] = appendRange([]int{}, 130, 160)

	return nodes
}

