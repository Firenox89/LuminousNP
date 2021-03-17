package main

import (
	"log"
	"utils/nodeMCU"
	"utils/utils"
	"utils/web"
)

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
		&[]web.Effect{},
		&colorPaletteList,
		nodeMCUController,
		func(request web.SetConfigRequest) {
			processNodesConfig(request, nodeMCUController)
		})
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
		//err = effectList[request.Config.Effect].Handler(node, request.Config)
	}
	if err != nil {
		log.Printf("Failed set node %s, err %v", node.ID, err)
	}
}
