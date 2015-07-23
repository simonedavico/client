package main

import (
	"github.com/codegangsta/cli"
	"cloud/benchflow/utils/logging"
	"os"
	"cloud/benchflow/commands"
)

var Commands = []cli.Command{
	commandConf,
//	commandDocker,
	commandDeploy,
	commandUnDeploy,
}

var commandConf = cli.Command{
	Name:        "conf",
	Usage:       "Show the configuration",
	Description: ``,
	Action:      commands.DoConf,
	Subcommands: []cli.Command{
		{
			Name:        "servers",
			Usage:       "Shows the servers configuration",
			Description: ``,
			Action:      commands.DoConfServer,
		},
//		{
//			Name:        "docker",
//			Usage:       "Shows the docker configuration",
//			Description: ``,
//			Action:      commands.DoConfDocker,
//		},
		{
			Name:        "benchflow",
			Usage:       "Shows the benchflow configuration",
			Description: ``,
			Action:      commands.DoConfBenchFlow,
		},
		{
			Name:        "credentials",
			Usage:       "Shows the credentials configuration",
			Description: ``,
			Action:      commands.DoConfCredentials,
		},
	},
}

//var commandDocker = cli.Command{
//	Name:        "docker",
//	Usage:       "Handles docker setup and configuration on local or remote machines",
//	Description: `The docker command install docker on the defined servers (benchflow conf servers) and enable it for remote connections by means of Docker Remote API.
//   The server must use the Ubuntu operating system supported by docker: https://docs.docker.com/installation/ubuntulinux/.
//   By default the command install docker on the declared serverName passed by argument.`,
//	Action:      commands.DoInstallDocker,
//}

//TODO: Add description to expose all the installation details
var commandDeploy = cli.Command{
	Name:        "install",
	Usage:       "Handles the deployment of BenchFlow and its configuration",
	Description: ``,
	Action:      commands.DoInstallBenchFlow,
}

var commandUnDeploy = cli.Command{
	Name:        "uninstall",
	Usage:       "Handles the undeployment of BenchFlow",
	Description: ``,
	Action:      commands.DoUninstallBenchFlow,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		logging.Log.Debug(v...)
	}
}

func assert(err error) {
	if err != nil {
		logging.Log.Fatal(err)
	}
}
