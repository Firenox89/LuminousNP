package nodeMCU

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ConnectedNode struct {
	IP          string
	ID          string
	LedCount    int
	BytesPerLED int
	isConnected bool
	Connection  net.Conn `json:"-"`
}

func NewConnectionNode(
	IP string,
	ID string,
	LedCount int,
	BytesPerLED int,
	Connection net.Conn) *ConnectedNode {
	node := &ConnectedNode{
		IP:          IP,
		ID:          ID,
		LedCount:    LedCount,
		BytesPerLED: BytesPerLED,
		isConnected: true,
		Connection:  Connection,
	}
	node.startHeartbeat()
	return node
}

func (n *ConnectedNode) startHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				_, err := n.Connection.Write([]byte("Ping"))
				if err != nil {
					log.Printf("ping failed for " + n.ID)
					done <- true
					ticker.Stop()
					n.isConnected = false
				}
				buffer := make([]byte, 64)
				err = n.Connection.SetReadDeadline(time.Now().Add(1 * time.Second))
				_, err = n.Connection.Read(buffer)
				if err != nil {
					log.Printf("ping failed for " + n.ID)
					done <- true
					ticker.Stop()
					n.isConnected = false
				}
			}
		}
	}()
}

func (n *ConnectedNode) SendEffectData(effectData []byte) {
	log.Printf("Try to send effect file, size %d", len(effectData))

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(len(effectData)))
	if err != nil {
		log.Printf("Failed to write buffer", err)
	}
	_, err = n.Connection.Write(buf.Bytes())

	_, err = n.Connection.Write(effectData)
	if err != nil {
		log.Printf("Failed to send effect file", err)
	}
}

func (n *ConnectedNode) StartEffect(effectId int) error {
	req, err := http.NewRequest("POST", "http://"+n.IP+"/startEffect?id="+strconv.Itoa(effectId), nil)
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

func (n *ConnectedNode) PowerOn() error {
	req, err := http.NewRequest("POST", "http://"+n.IP+"/on", nil)
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
