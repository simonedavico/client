package commands

import (
	"github.com/codegangsta/cli"
//	"cloud/benchflow/docker"
	"cloud/benchflow/deploy"
	"cloud/benchflow/environment"
	"cloud/benchflow/utils/configuration"
)

func DoInstallDocker(c *cli.Context) {
//	docker.InstallDocker()
}

func DoInstallBenchFlow(c *cli.Context) {
	deploy.InstallBenchFlow()
}

func DoUninstallBenchFlow(c *cli.Context) {
	deploy.UninstallBenchFlow()
}

func DoConf(c *cli.Context) {
	//TODO: print everything
}

func DoConfServer(c *cli.Context) {
	
	configuration.ShowServerConfiguration(environment.Env.Servers)
}

//TODO: move in utils/configuration
//func DoConfDocker(c *cli.Context) {
//
//	var docker structs.Docker = environment.Env.Docker
//
//	fmt.Println("The configuration file defines the following configuration for Docker: ")
//	fmt.Println("")
//
//	fmt.Println("\tLocal Certificates Path: " + docker.LocalCertificatesPath)
//	fmt.Println("\tContainer Hostname: " + docker.ContainerHostname)
//	fmt.Println("")
//
//}

func DoConfBenchFlow(c *cli.Context) {

	configuration.ShowBenchFlowConfiguration(environment.Env.BenchFlow)

}

func DoConfCredentials(c *cli.Context) {

	configuration.ShowCredentialsConfiguration(environment.Env.Credentials)

}
