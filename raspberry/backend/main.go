package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"utils/nodeMCU"
	"utils/utils"
	"utils/web"
)

var (
	effectStore      = make(map[string][]byte)
	effectStoreMutex = sync.RWMutex{}
)

var effectList = []web.Effect{
	{0, "Just White", false, false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		return node.PowerOn(config.Brightness)
	}},
	{1, "Single Color", true, false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		return node.ColorFill(utils.ApplyBrightnessToColorHex(config.Color, config.Brightness))
	}},
	{2, "Fade", false, true, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData := utils.GenerateFadeFromPalette(node.BytesPerLED, node.LedCount, colorPaletteList[config.ColorPaletteId].Colors, config.Brightness)
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
	{3, "Rotation", false, true, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData := utils.GenerateRotationFromPalette(node.BytesPerLED, node.LedCount, colorPaletteList[config.ColorPaletteId].Colors, config.Brightness)
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
	{4, "Restart Node", false, false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		return node.Restart()
	}},
	{5, "Restart Pi", false, false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		if err := exec.Command("sudo", "reboot").Run(); err != nil {
			fmt.Println("Failed to initiate shutdown:", err)
		}
		return nil
	}},
}

var colorPaletteList = []web.ColorPalette{
	{0, utils.OddlyInsertBlack(utils.BuildRedToGreenRamp())},
	{1, utils.OddlyInsertBlack(utils.BuildRedToBlueRamp())},
	{2, utils.OddlyInsertBlack(utils.BuildGreenToRedRamp())},
	{3, utils.OddlyInsertBlack(utils.BuildGreenToBlueRamp())},
	{4, utils.OddlyInsertBlack(utils.BuildBlueToRedRamp())},
	{5, utils.OddlyInsertBlack(utils.BuildBlueToGreenRamp())},

	{6, utils.RevertLoop(utils.BuildRedToGreenRamp())},
	{7, utils.RevertLoop(utils.BuildRedToBlueRamp())},
	{8, utils.RevertLoop(utils.BuildGreenToRedRamp())},
	{9, utils.RevertLoop(utils.BuildGreenToBlueRamp())},
	{10, utils.RevertLoop(utils.BuildBlueToRedRamp())},
	{11, utils.RevertLoop(utils.BuildBlueToGreenRamp())},

	{12, utils.BuildRainbowPalette()},
	{13, utils.BuildHappyPalette()},
}

func main() {
	nodeMCUController := nodeMCU.NewController()

	web.ServeWeb(
		&effectList,
		&colorPaletteList,
		&nodeMCUController.ConnectedMCUs,
		func(request web.SetConfigRequest) {
			processNodesConfig(request, nodeMCUController)
		},
		func(request web.RegistrationRequest, ip string) {
			nodeMCUController.RegisterNode(ip, request.ID, request.LedCount, request.BytesPerLED, request.Segments)
		},
		func(effectId string) []byte {
			effectStoreMutex.RLock()
			data := effectStore[effectId]
			effectStoreMutex.RUnlock()
			return data
		})

	//startUDPServer()
}

func processNodesConfig(request web.SetConfigRequest, controller *nodeMCU.Controller) {
	log.Printf("Process config request...")
	var nodes []*nodeMCU.ConnectedNode
	for _, requestedNode := range request.Nodes {
		log.Printf("Find node %s", requestedNode.ID)
		node := controller.GetNodeForID(requestedNode.ID)
		if node != nil {
			nodes = append(nodes, node)
		} else {
			log.Printf("warning: no node found for %s", requestedNode.ID)
		}
	}
	log.Printf("Config: Power %v, Color %s, Effect %d, palette %d", request.Config.Power, request.Config.Color, request.Config.Effect, request.Config.ColorPaletteId)
	for _, node := range nodes {
		go sendConfig(request, node)
	}
}

func sendConfig(request web.SetConfigRequest, node *nodeMCU.ConnectedNode) {
	var err error
	if !request.Config.Power {
		err = node.PowerOff()
	} else {
		err = effectList[request.Config.Effect].Handler(node, request.Config)
	}
	if err != nil {
		log.Printf("Failed set node %s, err %v", node.ID, err)
	}
}
