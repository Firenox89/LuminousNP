package nodeMCU

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"utils/utils"
)

type ColorPalette struct {
	Name   string
	Colors []string
}

type Effect struct {
	Name    string
	handler func(
		ctx context.Context,
		palette ColorPalette,
		colorState utils.ColorState,
		onChange func(),
	)
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
		colorState utils.ColorState,
		onChange func(),
	) {
		utils.StartZScannerSynced(ctx, palette.Colors, colorState, onChange)
	}},
	{Name: "Rainbow", handler: func(
		ctx context.Context,
		palette ColorPalette,
		colorState utils.ColorState,
		onChange func(),
	) {
		utils.StartRainbowSynced(ctx, palette.Colors, colorState, onChange)
	}},
	{Name: "ZFlow", handler: func(
		ctx context.Context,
		palette ColorPalette,
		colorState utils.ColorState,
		onChange func(),
	) {
		utils.StartZFlowSynced(ctx, palette.Colors, colorState, onChange)
	}},
	{Name: "RainbowXFlow", handler: func(
		ctx context.Context,
		palette ColorPalette,
		colorState utils.ColorState,
		onChange func(),
	) {
		utils.StartRainbowXFlowSynced(ctx, palette.Colors, colorState, onChange)
	}},
	{Name: "RainbowYFlow", handler: func(
		ctx context.Context,
		palette ColorPalette,
		colorState utils.ColorState,
		onChange func(),
	) {
		utils.StartRainbowYFlowSynced(ctx, palette.Colors, colorState, onChange)
	}},
}

func (c *Controller) getStateUpdater() (utils.ColorState, func()) {
	stateSizeX, stateSizeY, stateSizeZ := utils.GetMaxXYZ()
	colorState := utils.NewColorState(stateSizeX, stateSizeY, stateSizeZ)
	onChange := func() {
		c.updateCounter++
		for _, node := range c.ConnectedMCUs {
			go node.mapColorStateAndSend(colorState)
		}
	}
	return colorState, onChange
}

type Controller struct {
	ConnectedMCUs  []*ConnectedNode
	currentContext context.Context
	cancel         context.CancelFunc
	currentPalette ColorPalette
	currentEffect  Effect
	updateCounter int
}

func NewController() *Controller {
	controller := &Controller{}
	controller.RefreshNodes()
	controller.currentContext, controller.cancel = context.WithCancel(context.Background())
	controller.currentPalette = colorPaletteList[0]
	controller.currentEffect = effects[0]

	//go controller.updatesPerSecondCounter()

	return controller
}

func (c *Controller) updatesPerSecondCounter(){
	for {
		time.Sleep(time.Second)
		var ups = c.updateCounter
		c.updateCounter = 0
		log.Printf("ups %d\n", ups)
	}
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

func getTypeFromLEDCount(count int) utils.ShowCaseType {
	if count == 200 {
		return utils.ShowcaseType2
	} else if count == 188 {
		return utils.ShowcaseType1
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
	fmt.Printf("Restart Effect %s palette %s", c.currentEffect.Name, c.currentPalette.Name)
	c.cancel()
	c.currentContext, c.cancel = context.WithCancel(context.Background())
	colorState, onChange := c.getStateUpdater()
	go c.currentEffect.handler(c.currentContext, c.currentPalette, colorState, onChange)
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
