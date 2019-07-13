package mcsrvstat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	backgroundGraphic image.Image
)

// ServerStatus is JSON Data returned by mcsrvstat api
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

	userIP string
}

// Load fonts and graphics
func init() {
	var err error
	backgroundGraphic, err = gg.LoadImage("./resources/background.jpg")
	if err != nil {
		panic("background.jpg not found")
	}
}

// Query return server status from it's address
func Query(address string) (ServerStatus, error) {
	response, err := http.Get("https://api.mcsrvstat.us/1/" + address)
	if err != nil {
		return ServerStatus{}, err
	}
	// Server not found
	if response.StatusCode != 200 {
		return ServerStatus{}, fmt.Errorf("server responded with code other then 200")
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

	if status.Debug.Ping == false {
		return ServerStatus{}, errors.New("server is offline")
	}
	status.userIP = address
	return status, nil
}

// GenerateStatusImage generates image from ServerStatus
func (s ServerStatus) GenerateStatusImage() (imageBuf bytes.Buffer, err error) {
	fontLocation := "./resources/mcfont.ttf"
	about := "Generated using Minecraft Server Status Bot [Discord]"

	// Prepare strings
	trimmedMotd := strings.TrimSpace(s.Motd.Clean[0])
	players := fmt.Sprintf("%d players online of %d max", s.Players.Online, s.Players.Max)
	version := "Version: " + strings.TrimSpace(s.Version)
	// Create new image
	img := gg.NewContextForImage(backgroundGraphic)
	// Draw MOTD
	img.LoadFontFace(fontLocation, 50)
	img.SetRGB(1, 1, 1)
	img.DrawString(trimmedMotd, 50, 75)
	// Draw IP below MOTD
	img.LoadFontFace(fontLocation, 20)
	img.SetRGB(1, 1, 1)
	img.DrawString(s.userIP, 50, 140)
	// Draw number of players online
	img.LoadFontFace(fontLocation, 35)
	img.SetRGB(0, 1, 0)
	img.DrawString(players, 50, 265)
	// Draw version
	img.SetRGB(1, 1, 0)
	img.DrawString(version, 50, 335)
	// Draw about
	img.LoadFontFace(fontLocation, 10)
	img.SetRGB(1, 1, 1)
	img.DrawString(about, 1120-350, 700-30)
	// Draw players name's
	if len(s.Players.List) >= 1 && len(s.Players.List) < 11 {
		// Load font
		img.LoadFontFace(fontLocation, 15)
		// Set color
		img.SetRGB(1, 1, 1)
		// Add annotation
		img.DrawString("Player list:", 50, 400)
		// And goooo!
		for i, player := range s.Players.List {
			img.DrawString(player, 50, 425.0+float64(i)*25.0)
		}
	}
	// Encode image into png
	buf := bytes.Buffer{}
	png.Encode(&buf, img.Image())
	// Draw players name's
	return buf, nil
}
