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
}

var graphdef = map[string]mp.Graphs{
	"photon.rooms": {
		Label: "Photon Rooms",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "rooms", Label: "room", Diff: false},
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
}

func getPhotonStats(p PhotonStatsPlugin, name string) (string, error) {
	endPointUrl := p.Url + p.AppId + "/" + p.Region + "/" + name
	u, err := url.Parse(endPointUrl)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	now := time.Now()
	end := now.UTC()
	start := end.Add(-time.Duration(p.SecondsAgo) * time.Second)
	q.Set("start", start.Format("2006-01-02T03:04:05"))
	q.Set("end", end.Format("2006-01-02T03:04:05"))
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", p.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP status error: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(body[:]), nil
}

// FetchMetrics interface for mackerelplugin
func (u PhotonStatsPlugin) FetchMetrics() (stats map[string]interface{}, err error) {
	ccu, err := getPhotonStats(u, "ccu")
	rooms, err := getPhotonStats(u, "rooms")
	rejects, err := getPhotonStats(u, "rejects")
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"rooms": rooms, "ccu": ccu, "rejects": rejects}, nil
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
	flag.Parse()

	var photon PhotonStatsPlugin
	photon.SecondsAgo = *optSecondsAgo
	photon.AppId = *optAppid
	photon.Url = *optUrl
	photon.Region = *optRegion
	photon.Token = *optToken

	helper := mp.NewMackerelPlugin(photon)
	helper.Run()
}
