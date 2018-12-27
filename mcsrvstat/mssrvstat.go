package mcsrvstat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ServerStatus struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Debug struct {
		Ping          bool `json:"ping"`
		Query         bool `json:"query"`
		Srv           bool `json:"srv"`
		Querymismatch bool `json:"querymismatch"`
		Ipinsrv       bool `json:"ipinsrv"`
		Animatedmotd  bool `json:"animatedmotd"`
		Proxypipe     bool `json:"proxypipe"`
		Cachetime     int  `json:"cachetime"`
		DNS           struct {
			Srv []interface{} `json:"srv"`
			A   []struct {
				Host  string `json:"host"`
				Class string `json:"class"`
				TTL   int    `json:"ttl"`
				Type  string `json:"type"`
				IP    string `json:"ip"`
			} `json:"a"`
		} `json:"dns"`
	} `json:"debug"`
	Motd struct {
		Raw   []string `json:"raw"`
		Clean []string `json:"clean"`
		HTML  []string `json:"html"`
	} `json:"motd"`
	Players struct {
		Online int      `json:"online"`
		Max    int      `json:"max"`
		List   []string `json:"list"`
	} `json:"players"`
	Version  string `json:"version"`
	Protocol int    `json:"protocol"`
	Hostname string `json:"hostname"`
}

func Query(address string) (ServerStatus, error) {
	response, err := http.Get("https://api.mcsrvstat.us/1/" + address)
	if err != nil {
		return ServerStatus{}, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ServerStatus{}, err
	}

	status := ServerStatus{}
	err = json.Unmarshal(data, &status)
	if err != nil {
		return ServerStatus{}, err
	}

	return status, nil
}
