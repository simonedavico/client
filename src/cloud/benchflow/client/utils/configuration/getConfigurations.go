package configuration

import (
	"cloud/benchflow/client/structs"
	"cloud/benchflow/client/utils/logging"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

//TODO: find a better way to handle it

//Update when settings.yml changes by running go generate ./src/cloud/benchflow/utils/configuration/ in the project's root

//go:generate ../../../../../_tools/embed file -var settings --source ../../../../../configuration/settings.yml
var settings = "localCertificateFolder: \"certs\"\nremoteCertificatesFolder: \"benchflow/certs\"\nremoteDockerComposeProjects: \"benchflow/compose\""

//TODO: refactor to avoid code cloning

func GetConfServer() map[string]structs.Server {

	//TODO: refactor to an utility function and avoid code cloning and struct when not needed
	var servers = viper.GetStringMap("servers")
	var serversMap map[string]structs.Server = make(map[string]structs.Server)

	var key string
	var val interface{}

	//TODO: read in the order used in the config file
	for key, val = range servers {
		var server structs.Server
		err := mapstructure.Decode(val, &server)

		if err != nil {
			logging.Log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Unable to decode into struct")
		} else {
			serversMap[key] = server
		}
	}

	return serversMap
}

func GetSettings() structs.Settings {

	var settMap structs.Settings

	err := yaml.Unmarshal([]byte(settings), &settMap)

	if err != nil {
		logging.Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to decode into struct")
	}

	return settMap
}

func GetConfBenchFlow() structs.BenchFlow {

	//TODO: refactor to an utility function and avoid code cloning and struct when not needed
	var dockerImagesConf = viper.GetStringMap("benchflow.dockerImages")
	var benchflow structs.BenchFlow
	var dockerImages structs.DockerImages

	var key string
	var val interface{}

	//TODO: read in the order used in the config file
	for key, val = range dockerImagesConf {
		var imageDetails structs.ImageDetails
		err := mapstructure.Decode(val, &imageDetails)

		if err != nil {
			logging.Log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Unable to decode into struct")
		} else {
			if key == "dockerSwarm" {
				dockerImages.DockerSwarm = imageDetails
			} else if key == "dockerCompose" {
				dockerImages.DockerCompose = imageDetails
			} else if key == "cassandra" {
				dockerImages.Cassandra = imageDetails
			} else if key == "rancher" {
				dockerImages.Rancher = imageDetails
			} else if key == "rancherCompose" {
				dockerImages.RancherCompose = imageDetails
			} else if key == "consul" {
				dockerImages.Consul = imageDetails
			} else if key == "registrator" {
				dockerImages.Registrator = imageDetails
			}
		}
	}

	benchflow.DockerImages = dockerImages

	return benchflow
}

func GetConfCredentials() structs.Credentials {

	//TODO: refactor to an utility function and avoid code cloning and struct when not needed
	var credentialsConf = viper.GetStringMap("credentials")
	var credentials structs.Credentials

	var key string
	var val interface{}

	//TODO: read in the order used in the config file
	for key, val = range credentialsConf {
		var credentialsDetails structs.CredentialsDetails
		err := mapstructure.Decode(val, &credentialsDetails)

		if err != nil {
			logging.Log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Unable to decode into struct")
		} else {
			if key == "rancher" {
				credentials.Rancher = credentialsDetails
			} else if key == "cassandra" {
				credentials.Cassandra = credentialsDetails
			}
		}
	}

	return credentials
}

