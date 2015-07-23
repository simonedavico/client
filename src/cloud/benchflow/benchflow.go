package main

import (
	"github.com/codegangsta/cli"
	"os"
	"github.com/spf13/viper"
	"cloud/benchflow/environment"
	"cloud/benchflow/utils/logging"
	log "github.com/Sirupsen/logrus"
	
//	"fmt"
//	"encoding/json"
//	benchFlowDockerUtils "cloud/benchflow/utils/docker"
//	"github.com/fsouza/go-dockerclient"
)

//TODO, idea for refactoring to env, flag etc... https://github.com/docker/swarm/tree/master/cli

func main() {
	
	//TODO: improve to support reading from ENV and flags (for portability)
	
	//Read configuration files
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")  //production
	err := viper.ReadInConfig()
	
	//Initialise the logger
	logging.InitialiseLogger()
	
//	logging.Log.Debug("Test")
	
	if err != nil { // Handle errors reading the config file
	    logging.Log.WithFields(log.Fields{
			  "error": err,
			  }).Fatal("Fatal error config file")
	}

	//Validate the configuration file
	//TODO: validate config file structure
	if !viper.IsSet("servers") {
	    logging.Log.Fatal("A config file is not present or is not valid")
	}
	
	//Initialise environment data (and application configurations)
	environment.InitialiseEnvironment()
	

	//Handle the benchflow start
	app := cli.NewApp()
	app.Name = "BenchFlow"
	app.Version = Version
	app.Usage = "Performance Benchmarking Made Easy"
	app.Author = "Vincenzo Ferme"
	app.Email = "info@vincenzoferme.it"
	app.Commands = Commands

	app.Run(os.Args)
}