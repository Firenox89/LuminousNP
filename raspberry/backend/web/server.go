package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"utils/nodeMCU"
)

type Effect struct {
	ID         int
	Name       string
	NeedsColor bool
	Handler    func(node *nodeMCU.ConnectedNode, config LEDConfig) error `json:"-"`
}

type LEDConfig struct {
	Power    bool   `json:"power"`
	UseWhite bool   `json:"useWhite"`
	Color    string `json:"color"`
	Effect   int    `json:"effect"`
}

type RegistrationRequest struct {
	ID          string `json:"id"`
	LedCount    int    `json:"ledCount"`
	BytesPerLED int    `json:"bytesPerLed"`
	Segments    []int  `json:"segments"`
}

type Node struct {
	ID string `json:"ID"`
}

type SetConfigRequest struct {
	Config LEDConfig `json:"config"`
	Nodes  []Node    `json:"nodes"`
}

func ServeWeb(
	effectList *[]Effect,
	connectedMCUs *[]*nodeMCU.ConnectedNode,
	onApplyConfig func(request SetConfigRequest),
	onNodeRegisterRequest func(request RegistrationRequest, ip string),
	effectDataGetter func(nodeID string) []byte) {
	log.Printf("Start Web Server")
	setupWebAPI(effectList, connectedMCUs, onApplyConfig)

	setupNodeAPI(effectDataGetter, onNodeRegisterRequest)

	http.Handle("/", http.FileServer(http.Dir("dist")))

	log.Fatal(http.ListenAndServe(":80", nil))
}

func setupNodeAPI(effectDataGetter func(nodeID string) []byte, onNodeRegisterRequest func(request RegistrationRequest, ip string)) {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		decoder := json.NewDecoder(r.Body)
		var request RegistrationRequest
		err := decoder.Decode(&request)
		if err != nil {
			w.WriteHeader(500)
		} else {
			onNodeRegisterRequest(request, ip)
		}
	})

	http.HandleFunc("/lfs.img", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serve lfs.img to %s", r.RemoteAddr)
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/effectFile", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serve effect file to %s", r.RemoteAddr)
		nodeIDs, ok := r.URL.Query()["id"]
		if !ok {
			log.Printf("Effect parameter not found %v", r.URL.Query())
		} else {
			nodeID := nodeIDs[0]

			effectData := effectDataGetter(nodeID)
			w.Header().Add("Content-Length", strconv.Itoa(len(effectData)))

			_, err := w.Write(effectData)
			if err != nil {
				log.Printf("error on sending effect file %v", err)
			}
		}
	})
}

func setupWebAPI(effectList *[]Effect, connectedMCUs *[]*nodeMCU.ConnectedNode, onApplyConfig func(request SetConfigRequest)) {
	http.HandleFunc("/setConfig", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		decoder := json.NewDecoder(r.Body)
		var config SetConfigRequest
		err := decoder.Decode(&config)
		if err != nil {
			w.WriteHeader(500)
		} else {
			onApplyConfig(config)
		}
	})

	http.HandleFunc("/getEffectList", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(*effectList)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/getConnectedNodeMCUs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(*connectedMCUs)
		if err != nil {
			log.Fatal(err)
		}
	})
}
