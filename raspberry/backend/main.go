package main

import (
	"log"
	"utils/nodeMCU"
	"utils/utils"
	"utils/web"
)

func main() {
	nodeMCUController := nodeMCU.NewController()
	go nodeMCUController.StartControllerService()

	web.ServeWeb(&nodeMCUController.ConnectedMCUs, func(request web.SetConfigRequest) {
		SendNodeConfig(request, nodeMCUController)
	})

	//startUDPServer()
}

func SendNodeConfig(request web.SetConfigRequest, controller *nodeMCU.Controller) {
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
	log.Printf("Config: Power %v, Color %s, Effect %d", request.Config.Power, request.Config.Color, request.Config.Effect)
	for _, node := range nodes {
		go sendConfig(request, node)
	}
}

func sendConfig(request web.SetConfigRequest, node *nodeMCU.ConnectedNode) {
	var err error
	if !request.Config.Power {
		err = node.PowerOff()
	} else {
		switch request.Config.Effect {
		case 0:
			err = node.PowerOn()
			break
		case 1:
			err = node.ColorFill(request.Config.Color)
			break
		case 2:
			effectData := utils.GenerateColorSwitchEffect(int16(node.BytesPerLED), int16(node.LedCount))
			node.SendEffectData(effectData)
			break
		case 3:
			effectData := utils.GenerateRainbowFade(int16(node.BytesPerLED), int16(node.LedCount))
			node.SendEffectData(effectData)
			break
		case 4:
			effectData := utils.GenerateRunningRainbow(int16(node.BytesPerLED), int16(node.LedCount))
			node.SendEffectData(effectData)
			break
		}
	}
	if err != nil {
		log.Printf("Failed set node %s, err %v", node.ID, err)
	}
}
