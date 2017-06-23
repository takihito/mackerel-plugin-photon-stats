mackerel-plugin-photon-stats
=====================

Photon Status custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-photon-stats [-region=<region>] [-appid=<application id>] [-token=<auth token>] [-timeout=<seconds>] [-log] 
```

## Reference documentRequirements

- [Photon Analytics API : Introduction](http://doc-api.photonengine.com/ja-jp/analytics-api/current/intro.html)
- [Photon Analytics API : Detail](http://doc-api.photonengine.com/ja-jp/analytics-api/current/)

## Example of mackerel-agent.conf

```
[plugin.metrics.photon-stats]
command = "mackerel-plugin-photon-stats -region jp -appid xxxxxxxx -token **********
```


