mackerel-plugin-photon-stats
=====================

Photon Status custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-photon-stats [-region=<region>] [-appid=<application id>] [-token=<auth token>] [-timeout=<seconds>] [-log] 
```

## Install

install mackerel-plugin-photon-stats by mkr

please check mackerel-plugin-photon-stats/releases 

```
mkr plugin install takihito/mackerel-plugin-photon-stats@v0.1.0
```

## Example of mackerel-agent.conf

```
[plugin.metrics.photon-stats]
command = "mackerel-plugin-photon-stats -region jp -appid xxxxxxxx -token **********
```


