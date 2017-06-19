package photonstats

import (
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"io/ioutil"
	"net/http"
)

const (
	photonUrl    = "https://counter.photonengine.com/Counter/api/data/app/"
	photonRegion = "jp"
)

type PhotonStatsPlugin struct {
	Url    string
	AppId  string
	Region string
	Token  string
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

func getPhotonStats(url string, appId string, region string, token string, name string) (string, error) {
	apiUrl := url + appId + "/" + region + "/" + name
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", token)

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
	ccu, err := getPhotonStats(u.Url, u.AppId, u.Region, u.Token, "ccu")
	rooms, err := getPhotonStats(u.Url, u.AppId, u.Region, u.Token, "rooms")
	rejects, err := getPhotonStats(u.Url, u.AppId, u.Region, u.Token, "rejects")
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
	optAppid := flag.String("appid", "", "App Id")
	optUrl := flag.String("url", photonUrl, "Photon analytivs api url")
	optRegion := flag.String("region", photonRegion, "region")
	optToken := flag.String("token", "", "Authorization Token")
	flag.Parse()

	var photon PhotonStatsPlugin
	photon.AppId = *optAppid
	photon.Url = *optUrl
	photon.Region = *optRegion
	photon.Token = *optToken

	helper := mp.NewMackerelPlugin(photon)
	helper.Run()
}
