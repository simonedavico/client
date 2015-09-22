package deploy

import (
	benchFlowDocker "cloud/benchflow/docker"
	"cloud/benchflow/environment"
	"cloud/benchflow/structs"
	"cloud/benchflow/utils/configuration"
	benchFlowDockerUtils "cloud/benchflow/utils/docker"
	"cloud/benchflow/utils/logging"
	"cloud/benchflow/utils/servers"
	//	"encoding/json"
	//	"fmt"
	"github.com/fsouza/go-dockerclient"
	//	"github.com/gocql/gocql"
	"bytes"
	"strconv"
	"strings"
	"time"
)

//TODO: refactor
//TODO: add TLS for communication https://www.consul.io/docs/agent/encryption.html

func InstallConsul() {

	//BenchFlow configuration
	benchFlowConf := environment.Env.BenchFlow

	consulName, consulTag := benchFlowDockerUtils.GetImageAndTag(benchFlowConf.DockerImages.Consul.Image)

	installConsul(consulName, consulTag)

}

func installConsul(consulName string, consulTag string) {

	//Deploy Consul on Master and Analyser Servers
	var master = servers.GetMasterServer(environment.Env.Servers)

		//Pull image on master and analyser servers (TODO: When swarm support image pull with affinity, change it to Swarm)
		deployConsulOnServer(master, consulName, consulTag)
	
		//Iterate analyzer and deploy consul on them
	var serversList = servers.GetAnalyserNotMasterServers(environment.Env.Servers)
	//
		var server structs.Server
	
		//Install consul
		for _, server = range serversList {
	
			//Pull image on server (TODO: When swarm support image pull with affinity, change it to Swarm)
			deployConsulOnServer(server, consulName, consulTag)
	
		}

	//Start consul on master
	var numServers = len(serversList) + 1 //Analysers server + master server
	logging.Log.Debug("numServers: ")
	logging.Log.Debug(numServers)
	masterContainerID := runConsul(master, consulName, consulTag, numServers, "master")

	waitUntilTheConsulMasterIsStarted(master, masterContainerID)

	//Iterate analyzer and start consul on them
	for _, server = range serversList {

		runConsul(server, consulName, consulTag, numServers, "analyser")

	}
}

func deployConsulOnServer(server structs.Server, consulName string, consulTag string) {

	var client = benchFlowDocker.GetNewServerDockerClient(server)

	benchFlowDockerUtils.DeployImage(client, consulName, consulTag)

}

//Reference: https://registry.hub.docker.com/u/gliderlabs/consul/
func runConsul(server structs.Server, consulName string, consulTag string, totalNode int, role string) string {

	var client = benchFlowDocker.GetNewServerDockerClient(server)

	var serverIP = configuration.GetServerIP(server)
	var masterIP = configuration.GetServerIP(servers.GetMasterServer(environment.Env.Servers))

	consulMasterConfigToBuild := benchFlowDockerUtils.DockerConfigBuilder.
		Hostname("consul_"+role+"_"+strings.ToLower(server.Name)).
		Image(consulName, consulTag).
		ExposedPorts(map[docker.Port]struct{}{
		"8300/tcp": {},
		"8301/tcp": {},
		"8301/udp": {},
		"8302/tcp": {},
		"8302/udp": {},
		"8400/tcp": {},
		"8500/tcp": {},
		"8600/tcp": {},
		"8600/udp": {},
	}).
		Volumes(map[string]struct{}{
		environment.Env.BenchFlow.DockerImages.Consul.DataRootDirectory + "/consul": {},
	}).
		AddEnv("constraint:node==" + server.Hostname).
		AddCmd("-server").
		AddCmd("-advertise").
		AddCmd(serverIP)

	if totalNode == 1 {
		consulMasterConfigToBuild = consulMasterConfigToBuild.AddCmd("-bootstrap")
	} else if strings.EqualFold(role, "master") {
		consulMasterConfigToBuild = consulMasterConfigToBuild.AddCmd("-bootstrap-expect")
		consulMasterConfigToBuild = consulMasterConfigToBuild.AddCmd(strconv.Itoa(totalNode))
	} else {
		consulMasterConfigToBuild = consulMasterConfigToBuild.AddCmd("-join")
		consulMasterConfigToBuild = consulMasterConfigToBuild.AddCmd(masterIP)
	}

	consulMasterConfig := consulMasterConfigToBuild.Build()

	logging.Log.Debug(consulMasterConfig)

	consulMasterHostConfig := benchFlowDockerUtils.DockerHostConfigBuilder.
		AddBind(environment.Env.BenchFlow.DockerImages.Consul.DataRootDirectory + "/consul:/data:rw").
		PortBindings(getConsulPortBindings()).
		Build()

	logging.Log.Debug(consulMasterHostConfig)

	consul := benchFlowDockerUtils.CreateContainerOptionsBuilder.
		Name("consul_" + role).
		Config(&consulMasterConfig).
		HostConfig(&consulMasterHostConfig).
		Build()

	res, err := client.CreateContainer(consul)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(res)

	logging.Log.Debug("Err: ")
	logging.Log.Debug(err)

	// Start Container
	logging.Log.Debug("Test Run Container")

	startRes := client.StartContainer(res.ID, nil)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(startRes)

	return res.ID
}

//TODO: abstract to a builder allowing to define the port number, port type and binding
func getConsulPortBindings() map[docker.Port][]docker.PortBinding {

	//Ports
	ports := make(map[docker.Port][]docker.PortBinding, 8)

	numPorts := len(environment.Env.BenchFlow.DockerImages.Consul.RunPorts)

	//Bindings
	bindings := make([]docker.PortBinding, numPorts, numPorts)

	//Bindings
	for i, port := range environment.Env.BenchFlow.DockerImages.Consul.RunPorts {

		bindings[i] = docker.PortBinding{
			HostPort: port,
		}

	}

	ports["8300/tcp"] = []docker.PortBinding{bindings[0]}
	ports["8301/tcp"] = []docker.PortBinding{bindings[1]}
	ports["8301/udp"] = []docker.PortBinding{bindings[1]}
	ports["8302/tcp"] = []docker.PortBinding{bindings[2]}
	ports["8302/udp"] = []docker.PortBinding{bindings[2]}
	ports["8400/tcp"] = []docker.PortBinding{bindings[3]}
	ports["8500/tcp"] = []docker.PortBinding{bindings[4]}
	ports["8600/udp"] = []docker.PortBinding{bindings[5]}

	return ports
}

//TODO: refactor to a shared functionality
func waitUntilTheConsulMasterIsStarted(master structs.Server, masterContainerID string) {
	for {

		time.Sleep(5000 * time.Millisecond)

		var client = benchFlowDocker.GetNewServerDockerClient(master)

		var buf bytes.Buffer

		cassandraLogs := benchFlowDockerUtils.LogOptionsBuilder.
			Container(masterContainerID).
			OutputStream(&buf).
			Follow(false).
			Stdout(true).
			Stderr(true).
			Timestamps(false).
			Build()

		logErr := client.Logs(cassandraLogs)

		logging.Log.Debug(">>>>Log")
		logging.Log.Debug(buf.String())
		logging.Log.Debug(logErr)

		//Check for the last log message showed during the initialization
		if strings.Contains(buf.String(), "agent: failed to sync remote state: No cluster leader") {
			break
		}
	}
}
