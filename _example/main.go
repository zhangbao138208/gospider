package main

import (
	"gospider"
	log "github.com/sirupsen/logrus"
	_ "gospider/_example/rule/Leisusport"
	_ "gospider/_example/rule/baidunews"
	_ "gospider/_example/rule/dianping"
	_ "gospider/_example/rule/mojitianqi"
	_ "gospider/_example/rule/stackoverflow"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
	log.SetLevel(log.DebugLevel)
}

func main() {
	gs := gospider.New()
	log.Fatal(gs.Run())
}
