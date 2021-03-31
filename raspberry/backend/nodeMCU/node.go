package nodeMCU

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"net"
	"net/http"
	"strconv"
	"utils/utils"
)

type ConnectedNode struct {
	IP         string
	ID         string
	Type       utils.ShowCaseType
	Effects    []string
	Palettes   []string
	Brightness int
	On         bool
}

func NewConnectionNode(
	IP string,
	ID string,
	Type utils.ShowCaseType,
	Effects []string,
	Palettes []string,
	Brightness int,
	On bool,
) *ConnectedNode {
	return &ConnectedNode{
		IP:         IP,
		ID:         ID,
		Type:       Type,
		Effects:    Effects,
		Palettes:   Palettes,
		Brightness: Brightness,
		On:         On,
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

func (n *ConnectedNode) SendDatagram(data []byte) error {
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

func (n *ConnectedNode) mapColorStateAndSend(state utils.ColorState) {
	nodeData := utils.MapColorStateToNodes(n.Type, state)
	err := n.SendDatagram(generateDRGBPackage(nodeData))
	if err != nil {
		panic(err)
	}
}

func generateWARLSPackage(nodeData []colorful.Color) []byte {
	//1 = protocol
	//2 = time till web interfaces gets activated again
	data := []byte{1, 2}
	for index, color := range nodeData{
		r,g,b := color.RGB255()
		data = append(data, byte(index), r, g, b)
	}
	return data
}

func generateDRGBPackage(nodeData []colorful.Color) []byte {
	//2 = protocol
	//2 = time till web interfaces gets activated again
	data := []byte{2, 2}
	for _, color := range nodeData{
		r,g,b := color.RGB255()
		data = append(data, r, g, b)
	}
	return data
}
