package nodeMCU

import (
	"log"
	"strings"
)

type Controller struct {
	ConnectedMCUs []*ConnectedNode
}

func NewController() *Controller {
	controller := &Controller{}
	controller.RefreshNodes()
	return controller
}

func (c *Controller) RefreshNodes() {
	var wledNodes = ScanNetwork()
	log.Printf("%d nodes found.\n", len(wledNodes))
	var nodes []*ConnectedNode
	for _, node := range wledNodes {
		var lastIPFragment = node.IP[strings.LastIndex(node.IP, ".")+1:]
		var newNode = NewConnectionNode(
			node.IP,
			lastIPFragment,
			getTypeFromLEDCount(node.Info.Leds.Count),
			node.Effects,
			node.Palettes,
		)
		nodes = append(nodes, newNode)
	}
	c.ConnectedMCUs = nodes
}

func getTypeFromLEDCount(count int) ShowCaseType {
	if count == 200 {
		return ShowcaseType2
	} else if count == 188 {
		return ShowcaseType1
	} else {
		log.Fatalf("Unknown showcase type, led count %d", count)
		return 0
	}
}

func (c *Controller) GetNodeForID(id string) *ConnectedNode {
	for _, connectedNode := range c.ConnectedMCUs {
		if connectedNode.ID == id {
			return connectedNode
		}
	}
	return nil
}

