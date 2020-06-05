package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	go serveWeb()

	startControllerService()
	//startUDPServer()
}

type SetConfigRequest struct {
	Config LEDConfig `json:"config"`
}

type LEDConfig struct {
	Power    bool   `json:"power"`
	UseWhite bool   `json:"useWhite"`
	Color    string `json:"color"`
	Effect   int    `json:"effect"`
}

func serveWeb() {
	http.HandleFunc("/setConfig", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got %v", r.Body)

		decoder := json.NewDecoder(r.Body)
		var config SetConfigRequest
		err := decoder.Decode(&config)
		if err != nil {
			panic(err)
		}
		log.Printf("got %v", config)
		log.Printf("got %s", config.Config.Color)
	})

	http.Handle("/", http.FileServer(http.Dir("raspberry/web/dist")))

	log.Fatal(http.ListenAndServe(":8081", nil))
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
