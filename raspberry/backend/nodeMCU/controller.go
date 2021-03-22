package nodeMCU

import (
	"context"
	"log"
	"strings"
	"utils/utils"
)

type ColorPalette struct {
	Name   string
	Colors []string
}

type Effect struct {
	Name    string
	handler func(ctx context.Context, palette ColorPalette, caseType ShowCaseType, send func(data []byte) error)
}

var colorPaletteList = []ColorPalette{
	{"Red2Green", utils.OddlyInsertBlack(utils.BuildRedToGreenRamp())},
	{"Red2Blue", utils.OddlyInsertBlack(utils.BuildRedToBlueRamp())},
	{"Green2Red", utils.OddlyInsertBlack(utils.BuildGreenToRedRamp())},
	{"Green2Blue", utils.OddlyInsertBlack(utils.BuildGreenToBlueRamp())},
	{"Blue2Red", utils.OddlyInsertBlack(utils.BuildBlueToRedRamp())},
	{"Blue2Green", utils.OddlyInsertBlack(utils.BuildBlueToGreenRamp())},

	{"Red2Green", utils.RevertLoop(utils.BuildRedToGreenRamp())},
	{"Red2Blue", utils.RevertLoop(utils.BuildRedToBlueRamp())},
	{"Green2Red", utils.RevertLoop(utils.BuildGreenToRedRamp())},
	{"Green2Blue", utils.RevertLoop(utils.BuildGreenToBlueRamp())},
	{"Blue2Red", utils.RevertLoop(utils.BuildBlueToRedRamp())},
	{"Blue2Green", utils.RevertLoop(utils.BuildBlueToGreenRamp())},

	{"Rainbow", utils.BuildRainbowPalette()},
	{"Happy", utils.BuildHappyPalette()},
}

var effects = []Effect{
	{Name: "Scanner", handler: func(
		ctx context.Context,
		palette ColorPalette,
		caseType ShowCaseType,
		send func(data []byte) error,
	) {
		utils.StartScanner(ctx, palette.Colors, GetNodeMappingForType(caseType), send)
	}},
}

type Controller struct {
	ConnectedMCUs  []*ConnectedNode
	currentContext context.Context
	cancel         context.CancelFunc
	currentPalette ColorPalette
	currentEffect  Effect
}

func NewController() *Controller {
	controller := &Controller{}
	controller.RefreshNodes()
	controller.currentContext, controller.cancel = context.WithCancel(context.Background())
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
			node.State.Bri,
			node.State.On,
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

func (c *Controller) StopEffects() {
	c.cancel()
}

func (c *Controller) SetEffectId(id int) {
	c.currentEffect = effects[id]

	c.restartEffects()
}

func (c *Controller) SetPaletteId(id int) {
	c.currentPalette = colorPaletteList[id]

	c.restartEffects()
}

func (c *Controller) restartEffects() {
	c.cancel()
	c.currentContext, c.cancel = context.WithCancel(context.Background())
	for _, node := range c.ConnectedMCUs {
		go c.currentEffect.handler(c.currentContext, c.currentPalette, node.Type, node.SendWARLSDatagram)
	}
}

func (c *Controller) GetPaletteNames() []string {
	var names []string
	for i := 0; i < len(colorPaletteList); i++ {
		names = append(names, colorPaletteList[i].Name)
	}
	return names
}

func (c *Controller) GetEffectNames() []string {
	var names []string
	for i := 0; i < len(effects); i++ {
		names = append(names, effects[i].Name)
	}
	return names
}
