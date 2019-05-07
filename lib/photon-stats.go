package photonstats

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	photonUrl     = "https://counter.photonengine.com/Counter/api/data/app/"
	photonRegion  = "jp"
	endSecondsAgo = 180
	secondsAgo    = 90
)

type PhotonStatsPlugin struct {
	Url           string
	AppId         string
	Region        string
	Token         string
	EndSecondsAgo int
	SecondsAgo    int
	Timeout       int
	Log           bool
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

func (u PhotonStatsPlugin) getPhotonStats(name string) (string, error) {
	endPointUrl := u.Url + u.AppId + "/" + u.Region + "/" + name
	photonUrl, err := url.Parse(endPointUrl)
	if err != nil {
		return "", err
	}
	q := photonUrl.Query()
	now := time.Now()
	end := now.Add(-time.Duration(u.EndSecondsAgo) * time.Second)
	start := end.Add(-time.Duration(u.SecondsAgo) * time.Second)
	q.Set("start", start.UTC().Format("2006-01-02T15:04:05"))
	q.Set("end", end.UTC().Format("2006-01-02T15:04:05"))
	photonUrl.RawQuery = q.Encode()
	if u.Log {
		log.Printf("request_url:%s", photonUrl.String())
		log.Printf("appid:%s", u.AppId)
		log.Printf("token:%s", u.Token)
	}
	req, err := http.NewRequest("GET", photonUrl.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", u.Token)

	client := &http.Client{Timeout: time.Duration(u.Timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("URL:%s, Range:%s(%s) - %s(%s), HTTP status error: %d",
			photonUrl.String(), start.Format("2006-01-02T15:04:05"), start.UTC().Format("2006-01-02T15:04:05"),
			end.Format("2006-01-02T15:04:05"), end.UTC().Format("2006-01-02T15:04:05"), resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if u.Log {
		log.Printf("URL:%s, Range:%s(%s) - %s(%s), HTTP status error: %d",
			photonUrl.String(), start.Format("2006-01-02T15:04:05"), start.UTC().Format("2006-01-02T15:04:05"),
			end.Format("2006-01-02T15:04:05"), end.UTC().Format("2006-01-02T15:04:05"), resp.StatusCode)
		log.Printf("status:%d", resp.StatusCode)
		log.Printf("body:%s", string(body[:]))
	}
	return string(body[:]), nil
}

// FetchMetrics interface for mackerelplugin
func (u PhotonStatsPlugin) FetchMetrics() (stats map[string]interface{}, err error) {
	ccu, err := u.getPhotonStats("ccu")
	if err != nil {
		log.Printf("ccu_error:%s", err)
	}
	rooms, err := u.getPhotonStats("rooms")
	if err != nil {
		log.Printf("rooms_error:%s", err)
	}
	// NOTE Chatアプリケーションでのみ表示されます
	// 非Chatアプリはエラーステータスが返ります
	channels, err := u.getPhotonStats("channels")
	if err != nil {
		log.Printf("channels_error:%s", err)
	}
	rejects, err := u.getPhotonStats("rejects")
	if err != nil {
		log.Printf("reject_error:%s", err)
	}
	messages, err := u.getPhotonStats("messages")
	if err != nil {
		log.Printf("messages_error:%s", err)
	}
	bandwidth, err := u.getPhotonStats("bandwidth")
	if err != nil {
		log.Printf("bandwidth_error:%s", err)
	}
	bandwidthchat, err := u.getPhotonStats("bandwidthchat")
	if err != nil {
		log.Printf("bandwidthchat_error:%s", err)
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
	optEndSecondsAgo := flag.Int("end_seconds", endSecondsAgo, "seconds")
	optSecondsAgo := flag.Int("seconds", secondsAgo, "seconds")
	optAppid := flag.String("appid", "", "App Id")
	optUrl := flag.String("url", photonUrl, "Photon analytivs api url")
	optRegion := flag.String("region", photonRegion, "region")
	optToken := flag.String("token", "", "Authorization Token")
	optLog := flag.Bool("log", false, "Use logging")
	optTimeout := flag.Int("timeout", 10, "timeout")
	flag.Parse()

	var photon PhotonStatsPlugin
	photon.EndSecondsAgo = *optEndSecondsAgo
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
