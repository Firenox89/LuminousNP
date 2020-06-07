package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ConnectedNodeMCU struct {
	IP          string
	ID          int
	LedCount    int
	BytesPerLED int
	Connection  net.Conn `json:"-"`
}

var connectedMCUs = make([]ConnectedNodeMCU, 0)

func main() {
	generateEffectTestFile()
	go serveWeb()

	startControllerService()
	//startUDPServer()
}

func generateEffectTestFile() {
	effectData := generateColorSwitchEffect(250, 4, 188)
	err := ioutil.WriteFile("current.effect", effectData, 0644)
	if err != nil {
		// handle error
	}
}

func generateEffectHeader(delay int16, bytesPerLed int16, ledCount int16) []byte {
	effect := EffectFormatHeader{1, delay, bytesPerLed, ledCount, 0}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, effect)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}

func generateColorSwitchEffect(delay int16, bytesPerLed int16, ledCount int16) []byte {
	var values = generateEffectHeader(delay, bytesPerLed, ledCount)

	for j := 0; j < 10; j++ {
		for i := 0; i < int(ledCount); i++ {
			values = append(values, 255)
			values = append(values, 0)
			values = append(values, 0)
			values = append(values, 0)
		}
		for i := 0; i < int(ledCount); i++ {
			values = append(values, 0)
			values = append(values, 255)
			values = append(values, 0)
			values = append(values, 0)
		}
		for i := 0; i < int(ledCount); i++ {
			values = append(values, 0)
			values = append(values, 0)
			values = append(values, 255)
			values = append(values, 0)
		}
		for i := 0; i < int(ledCount); i++ {
			values = append(values, 0)
			values = append(values, 0)
			values = append(values, 0)
			values = append(values, 255)
		}
	}
	return values
}

type EffectFormatHeader struct {
	SchemaVersion int16
	DelayPerFrame int16
	BytesPerLED   int16
	LedCount      int16
	Flags         int32
}

type SetConfigRequest struct {
	Config LEDConfig `json:"config"`
	Nodes  []Node    `json:"nodes"`
}

type LEDConfig struct {
	Power    bool   `json:"power"`
	UseWhite bool   `json:"useWhite"`
	Color    string `json:"color"`
	Effect   int    `json:"effect"`
}

type Node struct {
	ID int `json:"ID"`
}

func serveWeb() {
	http.HandleFunc("/setConfig", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}
		log.Printf("got %v", r.Body)

		decoder := json.NewDecoder(r.Body)
		var config SetConfigRequest
		err := decoder.Decode(&config)
		if err != nil {
			w.WriteHeader(500)
		} else {
			log.Printf("got %v", config)
			log.Printf("got %s", config.Config.Color)
			sendNodeConfig(config)
		}
	})

	http.HandleFunc("/getConnectedNodeMCUs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(connectedMCUs)
		if err != nil {
			log.Fatal(err)
		}
	})
	http.Handle("/", http.FileServer(http.Dir("raspberry/web/dist")))

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func sendNodeConfig(config SetConfigRequest) {
	var nodesToSendTo []ConnectedNodeMCU
	for _, nodeInConfig := range config.Nodes {
		for _, connectedNode := range connectedMCUs {
			if connectedNode.ID == nodeInConfig.ID {
				nodesToSendTo = append(nodesToSendTo, connectedNode)
				break
			}
		}
	}
	for _, nodeToSendTo := range nodesToSendTo {
		var req *http.Request
		var err error
		if !config.Config.Power {
			req, err = http.NewRequest("POST", "http://"+nodeToSendTo.IP+"/off", nil)

			if err != nil {
				log.Printf("Failed to create request", err)
			}
		} else {
			switch config.Config.Effect {
			case 0:
				req, err = http.NewRequest("POST", "http://"+nodeToSendTo.IP+"/on", nil)
				break
			case 1:
				body := url.Values{}
				body.Add("level", "0")
				body.Add("color", config.Config.Color)
				req, err = http.NewRequest("POST", "http://"+nodeToSendTo.IP+"/fill?level=0&color="+config.Config.Color, nil)
				break
			case 2:
				effectData := generateColorSwitchEffect(250, int16(nodeToSendTo.BytesPerLED), int16(nodeToSendTo.LedCount))
				sendEffectData(nodeToSendTo.Connection, effectData)
				req, err = http.NewRequest("POST", "http://"+nodeToSendTo.IP+"/playEffect", nil)
				break
			}

			if err != nil {
				log.Printf("Failed to create request", err)
			}
		}
		go sendRequest(err, req)
	}
}

func sendEffectData(conn net.Conn, effectData []byte) {
	log.Printf("Try to send effect file, size %d", len(effectData))

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(len(effectData)))
	if err != nil {
		log.Printf("Failed to write buffer", err)
	}
	log.Printf("buffer %d%d%d%d", buf.Bytes()[0],buf.Bytes()[1],buf.Bytes()[2],buf.Bytes()[3])
	_, err = conn.Write(buf.Bytes())

	_, err = conn.Write(effectData)
	if err != nil {
		log.Printf("Failed to send effect file", err)
	}
}

func sendRequest(err error, req *http.Request) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed %v", req, err)
	}
	resp.Body.Close()
}

func startUDPServer() {
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
			parseRegistration(p, remoteaddr.IP.String(), nil)
		}
		go sendUDPTestResponse(ser, remoteaddr)
	}
}

func startControllerService() {
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
			parseRegistration(p[:byteCount], ip, conn)
		}
	}
}

type NodeMCURegistrationRequest struct {
	ID          int `json:"id"`
	LedCount    int `json:"ledCount"`
	BytesPerLED int `json:"bytesPerLed"`
}

func parseRegistration(buffer []byte, ip string, conn *net.TCPConn) {
	log.Printf("got '%s'", buffer)
	request := NodeMCURegistrationRequest{}
	err := json.Unmarshal(buffer, &request)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("NodeMCU registered %s %v", ip, request)

		connectedMCUs = append(connectedMCUs, ConnectedNodeMCU{IP: ip, ID: request.ID, LedCount: request.LedCount, BytesPerLED: request.BytesPerLED, Connection: conn})
	}
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
