package nodeMCU

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"net/http"
)

type ConnectedNode struct {
	IP          string
	ID          string
	LedCount    int
	BytesPerLED int
	Connection  net.Conn `json:"-"`
}

func (n *ConnectedNode) SendEffectData(effectData []byte) {
	log.Printf("Try to send effect file, size %d", len(effectData))

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(len(effectData)))
	if err != nil {
		log.Printf("Failed to write buffer", err)
	}
	log.Printf("buffer %d%d%d%d", buf.Bytes()[0], buf.Bytes()[1], buf.Bytes()[2], buf.Bytes()[3])
	_, err = n.Connection.Write(buf.Bytes())

	_, err = n.Connection.Write(effectData)
	if err != nil {
		log.Printf("Failed to send effect file", err)
	}
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
	defer resp.Body.Close()

	return err
}
