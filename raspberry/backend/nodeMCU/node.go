package nodeMCU

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type ConnectedNode struct {
	IP                 string
	ID                 string
	LedCount           int
	BytesPerLED        int
	Segments           []int
	IsConnected        *bool
	HeartbeatTimestamp int64 `json:"-"`
}

func NewConnectionNode(
	IP string,
	ID string,
	LedCount int,
	BytesPerLED int,
	Segments []int) *ConnectedNode {
	node := &ConnectedNode{
		IP:                 IP,
		ID:                 ID,
		LedCount:           LedCount,
		BytesPerLED:        BytesPerLED,
		Segments:           Segments,
		IsConnected:        new(bool),
		HeartbeatTimestamp: time.Now().Unix(),
	}
	*node.IsConnected = true
	return node
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
