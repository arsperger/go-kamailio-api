package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/voip-services/go-kamailio-api/api/controllers"
	"gitlab.com/voip-services/go-kamailio-api/internal/jconf"
	"gitlab.com/voip-services/go-kamailio-api/internal/utils"

	log "github.com/romana/rlog"
)

func main() {

	if os.Getenv("RLOG_LOG_LEVEL") == "DEBUG" {
		log.Debug("Debug ON")
	}
	confPath := "config.json"

	configFilePath := flag.String("config", "", "path to config file")
	versionFlag := flag.Bool("v", false, "print the current version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Fprintf(os.Stdout, utils.Version())
		return
	}
	if *configFilePath != "" {
		confPath = *configFilePath
	}

	c := jconf.ServiceConfig{}

	err := c.LoadConfigFile(confPath)
	if err != nil {
		log.Errorf("could load config [%s]", err.Error())
		return
	}

	app := controllers.App{}

	if kamDbURL := os.Getenv("KAM_DB_URL"); kamDbURL == "" {
		app.Initialize(c.DB.DbURL)
	} else {
		app.Initialize(fmt.Sprintf("%s", kamDbURL))
	}
	app.NewClient(c.Kamailio.ServerAddr)
	app.RunServer(c.HTTP.ListenAddr)
}
