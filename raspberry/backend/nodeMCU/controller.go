package nodeMCU

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Controller struct {
	ConnectedMCUs []*ConnectedNode
}

func NewController() *Controller {
	controller := &Controller{make([]*ConnectedNode, 0)}
	controller.startHeartbeat()
	return controller
}

func (c *Controller) RegisterNode(ip string, id string, ledCount int, bytesPerLed int, segments []int) {
	listIndex := -1
	for index, connectedMCU := range c.ConnectedMCUs {
		if connectedMCU.ID == id {
			listIndex = index
			break
		}
	}
	node := NewConnectionNode(ip, id, ledCount, bytesPerLed, segments)
	if listIndex != -1 {
		c.ConnectedMCUs[listIndex] = node
	} else {
		log.Printf("Register Node %s ", id)
		c.ConnectedMCUs = append(c.ConnectedMCUs, node)
	}
}

func (c *Controller) startHeartbeat() {
	ticker := time.NewTicker(15 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				now := time.Now().Unix()
				for _, node := range c.ConnectedMCUs {
					if node.HeartbeatTimestamp < now-15 {
						*node.IsConnected = false
					}
				}
			}
		}
	}()
}

func (c *Controller) startUDPServer() {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("0.0.0.0"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)

		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendUDPTestResponse(ser, remoteaddr)
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

func sendUDPTestResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	var err error
	sec := time.Second / 60
	for err == nil {
		_, err = conn.WriteToUDP([]byte{
			1, 255, 0, 0, 0,
			2, 0, 255, 0, 0,
			3, 0, 0, 255, 0}, addr)
		time.Sleep(sec)
		_, err = conn.WriteToUDP([]byte{
			1, 0, 255, 0, 0,
			2, 0, 0, 255, 0,
			3, 255, 0, 0, 0}, addr)
		time.Sleep(sec)
		_, err = conn.WriteToUDP([]byte{
			1, 0, 0, 255, 0,
			2, 255, 0, 0, 0,
			3, 0, 255, 0, 0}, addr)
		time.Sleep(sec)
	}
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}
