package nodeMCU

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
)

type ShowCaseType int

const (
	ShowcaseType1 ShowCaseType = 1 << iota
	ShowcaseType2
)

type ConnectedNode struct {
	IP       string
	ID       string
	Type     ShowCaseType
	Effects  []string
	Palettes []string
}

func NewConnectionNode(
	IP string,
	ID string,
	Type ShowCaseType,
	Effects []string,
	Palettes []string,
) *ConnectedNode {
	return &ConnectedNode{
		IP:       IP,
		ID:       ID,
		Type:     Type,
		Effects:  Effects,
		Palettes: Palettes,
	}
}

func (n *ConnectedNode) StartEffect() error {
	req, err := http.NewRequest("POST", "http://"+n.IP+"/startEffect", nil)
	if err == nil {
		return sendRequest(req)
	}
	return err
}

func (n *ConnectedNode) Restart() error {
	req, err := http.NewRequest("POST", "http://"+n.IP+"/restart", nil)
	if err == nil {
		return sendRequest(req)
	}
	return err
}

func (n *ConnectedNode) PowerOff() error {
	req, err := http.NewRequest("POST", "http://"+n.IP+"/off", nil)
	if err == nil {
		return sendRequest(req)
	}
	return err
}

func (n *ConnectedNode) PowerOn(brightness int) error {
	req, err := http.NewRequest("POST", "http://"+n.IP+"/on?brightness="+strconv.Itoa(brightness), nil)
	if err == nil {
		return sendRequest(req)
	}
	return err
}

func (n *ConnectedNode) ColorFill(color string) error {
	req, err := http.NewRequest("POST", "http://"+n.IP+"/fill?level=0&color="+color, nil)
	if err == nil {
		return sendRequest(req)
	}
	return err
}

func sendRequest(req *http.Request) error {
	log.Printf("Send request %s", req.URL)
	client := &http.Client{}
	resp, err := client.Do(req)

	if resp != nil {
		log.Printf("Request status code " + resp.Status)
	}

	return err
}

func (n *ConnectedNode) SendWARLSDatagram(data []byte) error {
	addr := &net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("0.0.0.0"),
	}
	ser, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Listen error %v\n", err)
		return err
	}
	targetAddr, err := net.ResolveUDPAddr("udp", n.IP+":21324")
	if err != nil {
		fmt.Printf("Resolve error %v\n", err)
		return err
	}
	_, err = ser.WriteToUDP(data, targetAddr)
	if err != nil {
		fmt.Printf("Send error %v\n", err)
	}
	return err
}
