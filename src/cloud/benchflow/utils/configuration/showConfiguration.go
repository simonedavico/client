package configuration

import (
	"strconv"
	"fmt"
	"cloud/benchflow/structs"
)

func ShowServerConfiguration(serversMap map[string]structs.Server) {
	var key string
	var valServer structs.Server

	fmt.Println("The configuration file defines", len(serversMap), "servers.")
	fmt.Println("")

	for key, valServer = range serversMap {
		fmt.Println(key + ": ")
		fmt.Println("\tName: " + valServer.Name)
		fmt.Println("\tAlias: " + valServer.Alias)
		fmt.Println("\tExternal IP: " + valServer.ExternalIP)
		fmt.Println("\tLocal IP: " + valServer.LocalIP)
		fmt.Println("\tHostname: " + valServer.Hostname)
		fmt.Println("\tDocker Port: " + strconv.Itoa(valServer.DockerPort))
		fmt.Println("\tSSH Port: " + strconv.Itoa(valServer.SshPort))
		fmt.Println("\tSudo User Username: " + valServer.SudoUserName)
		fmt.Println("\tSudo User Password: " + valServer.SudoUserPassword)
		fmt.Println("\tPurpose: " + valServer.Purpose)
		fmt.Println("")
	}
}

func ShowBenchFlowConfiguration(benchflow structs.BenchFlow) {
	
	fmt.Println("The configuration file defines the following configuration for BenchFlow: ")
	fmt.Println("")

	fmt.Println("\tDocker Swarm Image: " + benchflow.DockerImages.DockerSwarm.Image)
	fmt.Print("\tDocker Swarm Ports: ")
	fmt.Println(benchflow.DockerImages.DockerSwarm.RunPorts)
	if benchflow.DockerImages.DockerSwarm.DataRootDirectory != "" {
		fmt.Println("\tDocker Swarm Data Root Directory: " + benchflow.DockerImages.DockerSwarm.DataRootDirectory)
	}
	fmt.Println("")
	fmt.Println("\tDocker Compose Image: " + benchflow.DockerImages.DockerCompose.Image)
	fmt.Print("\tDocker Compose Ports: ")
	fmt.Println(benchflow.DockerImages.DockerCompose.RunPorts)
	if benchflow.DockerImages.DockerCompose.DataRootDirectory != "" {
		fmt.Println("\tDocker Compose Data Root Directory: " + benchflow.DockerImages.DockerCompose.DataRootDirectory)
	}
	fmt.Println("")
	fmt.Println("\tCassandra Image: " + benchflow.DockerImages.Cassandra.Image)
	fmt.Print("\tCassandra Ports: ")
	fmt.Println(benchflow.DockerImages.Cassandra.RunPorts)
	if benchflow.DockerImages.Cassandra.DataRootDirectory != "" {
		fmt.Println("\tCassandra Data Root Directory: " + benchflow.DockerImages.Cassandra.DataRootDirectory)
	}
	fmt.Println("")
	fmt.Println("\tRancher Image: " + benchflow.DockerImages.Rancher.Image)
	fmt.Print("\tRancher Ports: ")
	fmt.Println(benchflow.DockerImages.Rancher.RunPorts)
	if benchflow.DockerImages.Rancher.DataRootDirectory != "" {
		fmt.Println("\tRancher Data Root Directory: " + benchflow.DockerImages.Rancher.DataRootDirectory)
	}
	fmt.Println("")
	fmt.Println("\tRancher Compose Image: " + benchflow.DockerImages.RancherCompose.Image)
	fmt.Print("\tRancher compose Ports: ")
	fmt.Println(benchflow.DockerImages.RancherCompose.RunPorts)
	if benchflow.DockerImages.RancherCompose.DataRootDirectory != "" {
		fmt.Println("\tRancher Compose Data Root Directory: " + benchflow.DockerImages.RancherCompose.DataRootDirectory)
	}
	fmt.Println("")
	fmt.Println("\tConsul Image: " + benchflow.DockerImages.Consul.Image)
	fmt.Print("\tConsul Ports: ")
	fmt.Println(benchflow.DockerImages.Consul.RunPorts)
	if benchflow.DockerImages.Consul.DataRootDirectory != "" {
		fmt.Println("\tConsul Data Root Directory: " + benchflow.DockerImages.Consul.DataRootDirectory)
	}
	fmt.Println("")
	
}

func ShowCredentialsConfiguration(credentials structs.Credentials) {
	
	fmt.Println("The configuration file defines the following credentials: ")
	fmt.Println("")

	fmt.Println("\tCassandra Username: " + credentials.Cassandra.Username)
	fmt.Println("\tCassandra Password: " + credentials.Cassandra.Password)
	
	fmt.Println("")
	
	fmt.Println("\tRancher Username: " + credentials.Rancher.Username)
	fmt.Println("\tRancher Password: " + credentials.Rancher.Password)
	
	fmt.Println("")
	
}