package nodeMCU

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type RegistrationRequest struct {
	ID          string `json:"id"`
	LedCount    int    `json:"ledCount"`
	BytesPerLED int    `json:"bytesPerLed"`
}

type Controller struct {
	ConnectedMCUs []*ConnectedNode
}

func NewController() *Controller {
	return &Controller{make([]*ConnectedNode, 0)}
}

func (c *Controller) StartControllerService() {
	log.Printf("Listen on port 4488")
	p := make([]byte, 2048)
	addr := net.TCPAddr{
		Port: 4488,
		IP:   net.ParseIP("0.0.0.0"),
	}
	ser, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		conn, err := ser.AcceptTCP()
		ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
		log.Printf("Got connection %s", ip)
		byteCount, err := conn.Read(p)

		log.Printf("Read %d bytes, from %s %s", byteCount, ip, p)
		if err != nil {
			log.Printf("Some error  %v", err)
			continue
		} else {
			c.parseRegistration(p[:byteCount], ip, conn)
		}
	}
}

func (c *Controller) parseRegistration(buffer []byte, ip string, conn *net.TCPConn) {
	log.Printf("got '%s'", buffer)
	request := RegistrationRequest{}
	err := json.Unmarshal(buffer, &request)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("NodeMCU registered %s %v", ip, request)

		listIndex := -1
		for index, connectedMCU := range c.ConnectedMCUs {
			if connectedMCU.ID == request.ID {
				listIndex = index
				break
			}
		}
		node := NewConnectionNode(ip, request.ID, request.LedCount, request.BytesPerLED, conn)
		if listIndex != -1 {
			c.ConnectedMCUs[listIndex] = node
		} else {
			c.ConnectedMCUs = append(c.ConnectedMCUs, node)
		}
	}
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
		} else {
			c.parseRegistration(p, remoteaddr.IP.String(), nil)
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
