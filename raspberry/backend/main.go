package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"utils/nodeMCU"
)

func main() {
	nodeMCUController := nodeMCU.NewController()

	log.Printf("Start Web Server")

	http.HandleFunc("/setEffect", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err == nil {
			nodeMCUController.SetEffectId(id)
		} else {
			fmt.Printf("unable to parse effect id %s\n", r.URL.Query().Get("id"))
		}
	})

	http.HandleFunc("/setPalette", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err == nil {
			nodeMCUController.SetPaletteId(id)
		} else {
			fmt.Printf("unable to parse palette id %s\n", r.URL.Query().Get("id"))
		}
	})

	http.HandleFunc("/getEffectList", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(nodeMCUController.GetEffectNames())
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/getColorPaletteList", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(nodeMCUController.GetPaletteNames())
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/getConnectedNodeMCUs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(nodeMCUController.ConnectedMCUs)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/stopEffects", func(w http.ResponseWriter, r *http.Request) {
		nodeMCUController.StopEffects()
	})
	http.Handle("/", http.FileServer(http.Dir("dist")))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Printf("Could not use Port 80, try 1234. Reason: " + err.Error())
		log.Fatal(http.ListenAndServe(":1234", nil))
	}
}
