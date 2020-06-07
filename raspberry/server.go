package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type ConnectedNodeMCU struct {
	IP          string
	ID          int
	LedCount    int
	BytesPerLED int
}

var connectedMCUs = make([]ConnectedNodeMCU, 0)

func main() {
	go serveWeb()

	startControllerService()
	//startUDPServer()
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
			req, err = http.NewRequest("POST", "http://"+nodeToSendTo.IP+"/on", nil)

			if err != nil {
				log.Printf("Failed to create request", err)
			}
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}
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
			parseRegistration(p, remoteaddr.IP.String())
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
			parseRegistration(p[:byteCount], ip)
		}
	}
}

type NodeMCURegistrationRequest struct {
	ID          int `json:"id"`
	LedCount    int `json:"ledCount"`
	BytesPerLED int `json:"bytesPerLed"`
}

func parseRegistration(buffer []byte, ip string) {
	log.Printf("got '%s'", buffer)
	request := NodeMCURegistrationRequest{}
	err := json.Unmarshal(buffer, &request)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("NodeMCU registered %s %v", ip, request)

		connectedMCUs = append(connectedMCUs, ConnectedNodeMCU{IP: ip, ID: request.ID, LedCount: request.LedCount, BytesPerLED: request.BytesPerLED})
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
