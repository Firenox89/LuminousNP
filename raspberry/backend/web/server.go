package web

import (
	"encoding/json"
	"log"
	"net/http"
	"utils/nodeMCU"
)

type LEDConfig struct {
	Power    bool   `json:"power"`
	UseWhite bool   `json:"useWhite"`
	Color    string `json:"color"`
	Effect   int    `json:"effect"`
}

type Node struct {
	ID string `json:"ID"`
}

type SetConfigRequest struct {
	Config LEDConfig `json:"config"`
	Nodes  []Node    `json:"nodes"`
}

func ServeWeb(connectedMCUs *[]*nodeMCU.ConnectedNode, onApplyConfig func(request SetConfigRequest)) {
	log.Printf("Start Web Server")
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

	http.HandleFunc("/getConnectedNodeMCUs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(connectedMCUs)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/lfs.img", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serve lfs.img to %s", r.RemoteAddr)
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.Handle("/", http.FileServer(http.Dir("dist")))

	log.Fatal(http.ListenAndServe(":80", nil))
}
