package nodeMCU

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WLEDNode struct {
	IP string
	State struct {
		On         bool `json:"on"`
		Bri        int  `json:"bri"`
		Transition int  `json:"transition"`
		Ps         int  `json:"ps"`
		Pss        int  `json:"pss"`
		Pl         int  `json:"pl"`
		Ccnf       struct {
			Min  int `json:"min"`
			Max  int `json:"max"`
			Time int `json:"time"`
		} `json:"ccnf"`
		Nl struct {
			On   bool `json:"on"`
			Dur  int  `json:"dur"`
			Fade bool `json:"fade"`
			Mode int  `json:"mode"`
			Tbri int  `json:"tbri"`
		} `json:"nl"`
		Udpn struct {
			Send bool `json:"send"`
			Recv bool `json:"recv"`
		} `json:"udpn"`
		Lor     int `json:"lor"`
		Mainseg int `json:"mainseg"`
		Seg     []struct {
			ID    int     `json:"id"`
			Start int     `json:"start"`
			Stop  int     `json:"stop"`
			Len   int     `json:"len"`
			Grp   int     `json:"grp"`
			Spc   int     `json:"spc"`
			On    bool    `json:"on"`
			Bri   int     `json:"bri"`
			Col   [][]int `json:"col"`
			Fx    int     `json:"fx"`
			Sx    int     `json:"sx"`
			Ix    int     `json:"ix"`
			Pal   int     `json:"pal"`
			Sel   bool    `json:"sel"`
			Rev   bool    `json:"rev"`
			Mi    bool    `json:"mi"`
		} `json:"seg"`
	} `json:"state"`
	Info struct {
		Ver  string `json:"ver"`
		Vid  int    `json:"vid"`
		Leds struct {
			Count   int   `json:"count"`
			Rgbw    bool  `json:"rgbw"`
			Wv      bool  `json:"wv"`
			Pin     []int `json:"pin"`
			Pwr     int   `json:"pwr"`
			Maxpwr  int   `json:"maxpwr"`
			Maxseg  int   `json:"maxseg"`
			Seglock bool  `json:"seglock"`
		} `json:"leds"`
		Str      bool   `json:"str"`
		Name     string `json:"name"`
		Udpport  int    `json:"udpport"`
		Live     bool   `json:"live"`
		Lm       string `json:"lm"`
		Lip      string `json:"lip"`
		Ws       int    `json:"ws"`
		Fxcount  int    `json:"fxcount"`
		Palcount int    `json:"palcount"`
		Wifi     struct {
			Bssid   string `json:"bssid"`
			Rssi    int    `json:"rssi"`
			Signal  int    `json:"signal"`
			Channel int    `json:"channel"`
		} `json:"wifi"`
		Arch     string `json:"arch"`
		Core     string `json:"core"`
		Lwip     int    `json:"lwip"`
		Freeheap int    `json:"freeheap"`
		Uptime   int    `json:"uptime"`
		Opt      int    `json:"opt"`
		Brand    string `json:"brand"`
		Product  string `json:"product"`
		Mac      string `json:"mac"`
	} `json:"info"`
	Effects  []string `json:"effects"`
	Palettes []string `json:"palettes"`
}

func ScanNetwork() []WLEDNode {
	var localAddress = getLocalNetworkAddress().IP.String()
	var subNet = localAddress[:strings.LastIndex(localAddress, ".")+1]

	log.Printf("Scan sub net %v\n", subNet)
	var results = make(chan WLEDNode)
	for i := 1; i < 255; i++ {
		go sendQuery(subNet, i, results)
	}
	time.Sleep(time.Second * 6)
	var nodes []WLEDNode
	for true {
		select {
		case x, ok := <-results:
			if ok {
				nodes = append(nodes, x)
			} else {
				//channel closed
				return nodes
			}
			break
		default:
			fmt.Println("No value ready, moving on.")
			return nodes
		}
	}
	return nodes
}

func sendQuery(subNet string, i int, results chan WLEDNode) {
	var ip = subNet + strconv.Itoa(i)
	var url = "http://" + ip + "/json"
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%e", err)
		return
	}
	var wled WLEDNode
	if err := json.Unmarshal(bytes, &wled); err != nil {
		log.Printf("%s, parsing error %s", ip, err.Error())
		return
	}
	wled.IP = ip
	log.Printf("Found %s\n", ip)
	results <- wled
}

func getLocalNetworkAddress() *net.IPNet {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	//TODO what to do if its not the second one?
	return addrs[1].(*net.IPNet)
}
