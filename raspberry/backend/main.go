package main

import (
	"log"
	"sync"
	"utils/nodeMCU"
	"utils/utils"
	"utils/web"
)

var (
	effectStore = make(map[string][] byte)
    effectStoreMutex = sync.RWMutex{}
)


var effectList = []web.Effect{
	{0, "Just White", false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		return node.PowerOn()
	}},
	{1, "Fill Color", true, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		return node.ColorFill(config.Color)
	}},
	{2, "FadeInOut", true, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData, err := utils.GenerateColorFadeEffect(node.BytesPerLED, node.LedCount, config.Color)
		if err != nil {
			return err
		}
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
	{3, "RainbowFade", false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData := utils.GenerateRainbowFade(node.BytesPerLED, node.LedCount)
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
	{4, "Rainbow Rotation", false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData := utils.GenerateRunningRainbow(node.BytesPerLED, node.LedCount)
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
	{5, "Warm Rotation", false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData, err := utils.GenerateRunningWarmColors(node.BytesPerLED, node.LedCount)
		if err != nil {
			return err
		}
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
	{6, "Happy Rotation", false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData, err := utils.GenerateRunningHappyColors(node.BytesPerLED, node.LedCount)
		if err != nil {
			return err
		}
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
	{7, "Warm Fade", false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData := utils.GenerateWarmColorFade(node.BytesPerLED, node.LedCount)
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
	{8, "Interpolate Test", false, func(node *nodeMCU.ConnectedNode, config web.LEDConfig) error {
		effectData := utils.GenerateInterpolateTest(node.BytesPerLED, node.LedCount)
		effectStoreMutex.Lock()
		effectStore[node.ID] = effectData
		effectStoreMutex.Unlock()
		return node.StartEffect()
	}},
}

func main() {
	nodeMCUController := nodeMCU.NewController()

	web.ServeWeb(
		&effectList,
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
		err = effectList[request.Config.Effect].Handler(node, request.Config)
	}
	if err != nil {
		log.Printf("Failed set node %s, err %v", node.ID, err)
	}
}
