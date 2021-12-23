package main

import (
	"strings"

	"github.com/owenthereal/jqplay/cli"
	"github.com/owenthereal/jqplay/config"
	"github.com/owenthereal/jqplay/jq"
	"github.com/owenthereal/jqplay/server"
	log "github.com/sirupsen/logrus"
)

var GinMode = "debug"

func main() {
	log.SetLevel(log.WarnLevel)

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = jq.Init()
	if err != nil {
		log.Fatal(err)
	}
	conf.JQVer = jq.Version

	log.WithFields(log.Fields{
		"version": jq.Version,
		"path":    jq.Path,
	}).Info("initialized jq")

	log.WithFields(log.Fields{
		"host": conf.Host,
		"port": conf.Port,
	}).Infof("Starting server at %s:%s", conf.Host, conf.Port)

	if conf.Cli {
		// run CLI interface
		cli.New(conf).Start()
		return
	}

	srv := server.New(conf)
	err = srv.Start(GinMode)
	if err != nil && !strings.Contains(err.Error(), "Server closed") {
		log.WithError(err).Fatal("error starting sever")
	}
}
