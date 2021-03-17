package web

import (
	"encoding/json"
	"log"
	"net/http"
	"utils/nodeMCU"
)

type Effect struct {
	ID                int
	Name              string
	NeedsColor        bool
	NeedsColorPalette bool
	Handler           func(node *nodeMCU.ConnectedNode, config LEDConfig) error `json:"-"`
}

type ColorPalette struct {
	ID     int
	Colors []string
}

type LEDConfig struct {
	Power          bool   `json:"power"`
	UseWhite       bool   `json:"useWhite"`
	Color          string `json:"color"`
	ColorPaletteId int    `json:"colorPaletteId"`
	Effect         int    `json:"effect"`
	Brightness     int    `json:"brightness"`
}

type RegistrationRequest struct {
	ID          string `json:"id"`
	LedCount    int    `json:"ledCount"`
	BytesPerLED int    `json:"bytesPerLed"`
	Segments    []int  `json:"segments"`
}

type Node struct {
	ID             string `json:"ID"`
	ActiveSegments []int  `json:"segments"`
}

type SetConfigRequest struct {
	Config LEDConfig `json:"config"`
	Nodes  []Node    `json:"nodes"`
}

func ServeWeb(effectList *[]Effect, colorPaletteList *[]ColorPalette, connectedMCUs *nodeMCU.Controller, onApplyConfig func(request SetConfigRequest)) {
	log.Printf("Start Web Server")
	setupWebAPI(effectList, colorPaletteList, connectedMCUs, onApplyConfig)

	http.Handle("/", http.FileServer(http.Dir("dist")))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
  		log.Printf("Could not use Port 80, try 1234. Reason: " + err.Error())
		log.Fatal(http.ListenAndServe(":1234", nil))
	}
}

func setupWebAPI(effectList *[]Effect, colorPaletteList *[]ColorPalette, connectedMCUs *nodeMCU.Controller, onApplyConfig func(request SetConfigRequest)) {
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

	http.HandleFunc("/getColorPaletteList", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(*colorPaletteList)
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
