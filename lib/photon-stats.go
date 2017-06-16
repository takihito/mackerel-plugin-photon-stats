package photonstats

import (
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	region = "ja-jp"
)

type PhotonStatsPlugin struct {
	Url    string
	AppId  string
	Region string
}

var graphdef = map[string]mp.Graphs{
	"photon.rooms": {
		Label: "Photon Rooms",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "rooms", Label: "count", Diff: false},
		},
	},
	"photon.stats": {
		Label: "Photon Stats",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ccu", Label: "count", Diff: false},
			{Name: "rejects", Label: "count", Diff: false},
		},
	},
}

// FetchMetrics interface for mackerelplugin
func (u PhotonStatsPlugin) FetchMetrics() (stats map[string]interface{}, err error) {
	return map[string]interface{}{"rooms": 100, "ccu": 20, "reject": 5}, nil
}

// GraphDefinition interface for mackerelplugin
func (u PhotonStatsPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func Do() {
	var photon PhotonStatsPlugin
	helper := mp.NewMackerelPlugin(photon)
	helper.Run()
}
