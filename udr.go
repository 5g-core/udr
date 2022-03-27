package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/5g-core/udr/logger"
	udr_service "github.com/5g-core/udr/service"
	"github.com/5g-core/version"
)

var UDR = &udr_service.UDR{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "udr"
	appLog.Infoln(app.Name)
	appLog.Infoln("UDR version: ", version.GetVersion())
	app.Usage = "-n5gccfg common configuration file -udrcfg udr configuration file"
	app.Action = action
	app.Flags = UDR.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("UDR Run error: %v", err)
	}
}

func action(c *cli.Context) error {
	if err := UDR.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	UDR.Start()

	return nil
}
