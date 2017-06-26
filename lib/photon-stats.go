package photonstats

import (
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	photonUrl    = "https://counter.photonengine.com/Counter/api/data/app/"
	photonRegion = "jp"
	secondsAgo   = 90
)

type PhotonStatsPlugin struct {
	Url        string
	AppId      string
	Region     string
	Token      string
	SecondsAgo int
	Timeout    int
	Log        bool
}

var graphdef = map[string]mp.Graphs{
	"photon.rooms": {
		Label: "Photon Rooms",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "rooms", Label: "room", Diff: false},
		},
	},
	"photon.channel": {
		Label: "Photon Channel",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "channels", Label: "channel", Diff: false},
		},
	},
	"photon.stats": {
		Label: "Photon Stats",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ccu", Label: "ccu", Diff: false},
			{Name: "rejects", Label: "reject", Diff: false},
		},
	},
	"photon.message": {
		Label: "Photon Message",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "messages", Label: "message", Diff: false},
		},
	},
	"photon.bandwidth": {
		Label: "Photon Bandwidth",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "bandwidth", Label: "bandwidth", Diff: false},
			{Name: "bandwidthchat", Label: "bandwidth chat", Diff: false},
		},
	},
}

func getPhotonStats(p PhotonStatsPlugin, name string) (string, error) {
	endPointUrl := p.Url + p.AppId + "/" + p.Region + "/" + name
	u, err := url.Parse(endPointUrl)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	end := time.Now()
	start := end.Add(-time.Duration(p.SecondsAgo) * time.Second)
	q.Set("start", start.UTC().Format("2006-01-02T15:04:05"))
	q.Set("end", end.UTC().Format("2006-01-02T15:04:05"))
	u.RawQuery = q.Encode()
	if p.Log {
		log.Printf("request_url:%s", u.String())
		log.Printf("appid:%s", p.AppId)
		log.Printf("token:%s", p.Token)
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", p.Token)

	client := &http.Client{Timeout: time.Duration(p.Timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("URL:%s, Range:%s(%s) - %s(%s), HTTP status error: %d",
			u.String(), start.Format("2006-01-02T15:04:05"), start.UTC().Format("2006-01-02T15:04:05"),
			end.Format("2006-01-02T15:04:05"), end.UTC().Format("2006-01-02T15:04:05"), resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if p.Log {
		log.Printf("URL:%s, Range:%s(%s) - %s(%s), HTTP status error: %d",
			u.String(), start.Format("2006-01-02T15:04:05"), start.UTC().Format("2006-01-02T15:04:05"),
			end.Format("2006-01-02T15:04:05"), end.UTC().Format("2006-01-02T15:04:05"), resp.StatusCode)
		log.Printf("status:%d", resp.StatusCode)
		log.Printf("body:%s", string(body[:]))
	}
	return string(body[:]), nil
}

// FetchMetrics interface for mackerelplugin
func (u PhotonStatsPlugin) FetchMetrics() (stats map[string]interface{}, err error) {
	ccu, err := getPhotonStats(u, "ccu")
	rooms, err := getPhotonStats(u, "rooms")
	channels, err := getPhotonStats(u, "channels")
	rejects, err := getPhotonStats(u, "rejects")
	messages, err := getPhotonStats(u, "messages")
	bandwidth, err := getPhotonStats(u, "bandwidth")
	bandwidthchat, err := getPhotonStats(u, "bandwidthchat")
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ccu":           ccu,
		"rooms":         rooms,
		"channels":      channels,
		"rejects":       rejects,
		"messages":      messages,
		"bandwidth":     bandwidth,
		"bandwidthchat": bandwidthchat,
	}, nil
}

// GraphDefinition interface for mackerelplugin
func (u PhotonStatsPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func Do() {
	optSecondsAgo := flag.Int("seconds", secondsAgo, "seconds")
	optAppid := flag.String("appid", "", "App Id")
	optUrl := flag.String("url", photonUrl, "Photon analytivs api url")
	optRegion := flag.String("region", photonRegion, "region")
	optToken := flag.String("token", "", "Authorization Token")
	optLog := flag.Bool("log", false, "Use logging")
	optTimeout := flag.Int("timeout", 10, "timeout")
	flag.Parse()

	var photon PhotonStatsPlugin
	photon.SecondsAgo = *optSecondsAgo
	photon.AppId = *optAppid
	photon.Url = *optUrl
	photon.Region = *optRegion
	photon.Token = *optToken
	photon.Log = *optLog
	photon.Timeout = *optTimeout

	helper := mp.NewMackerelPlugin(photon)
	helper.Run()
}
