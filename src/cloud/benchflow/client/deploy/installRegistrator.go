package deploy

import (
	benchFlowDocker "cloud/benchflow/client/docker"
	benchFlowDockerUtils "cloud/benchflow/client/utils/docker"
	"cloud/benchflow/client/environment"
	"cloud/benchflow/client/utils/servers"
	"cloud/benchflow/client/structs"
	"github.com/fsouza/go-dockerclient"
	"cloud/benchflow/client/utils/logging"
)

func InstallRegistrator() {
	
	//BenchFlow configuration
	benchFlowConf := environment.Env.BenchFlow
	
	registratorName, registratorTag := benchFlowDockerUtils.GetImageAndTag(benchFlowConf.DockerImages.Registrator.Image)
	
	deployRegistrator(registratorName, registratorTag)
	
	installRegistrator(registratorName, registratorTag)
	
}

func deployRegistrator(registratorName string, registratorTag string) {
	
	//Iterate non sut and deploy resitrator on them
	var serversList = servers.GetNotSutServers(environment.Env.Servers)
	
	var server structs.Server

	//Install registrator
	for _, server = range serversList {
		
		var client = benchFlowDocker.GetNewServerDockerClient(server)
		
		benchFlowDockerUtils.DeployImage(client, registratorName, registratorTag)

	}

}

func installRegistrator(registratorName string, registratorTag string) {
	
	//Iterate non sut and deploy resitrator on thme
	var serversList = servers.GetNotSutServers(environment.Env.Servers)
	
	var server structs.Server

	//Install registrator
	for _, server = range serversList {
		
		runRegistrator(server,registratorName,registratorTag)
		
	}

}

func runRegistrator(server structs.Server, registratorName string, registratorTag string) {
	
	var client = benchFlowDocker.GetNewServerDockerClient(server)
	
	//The server with the consul master
	var master = servers.GetMasterServer(environment.Env.Servers)
	
	var masterStruct []structs.Server
	masterStruct = append(masterStruct, master)
	
	registratorServerConfig := benchFlowDockerUtils.DockerConfigBuilder.
														Hostname("registrator").
														Image(registratorName,registratorTag).
														Volumes(map[string]struct{}{
																"/var/run/docker.sock": {},
														}).
														//TODO: convert to a dynamic one that use configurations to identify the hostname/ip and the port
														//even better by using the hostname mapped to the IP
														AddCmd("consulkv://neha:8500/services").
														Build()

	logging.Log.Debug(registratorServerConfig)

	
	registratorServerHostConfig := benchFlowDockerUtils.DockerHostConfigBuilder.
																AddBind("/var/run/docker.sock:/tmp/docker.sock:rw").
																NetworkMode("host").
																RestartPolicy(docker.RestartPolicy{
																	Name: "always",
																}).
																//TODO: It seems to have problem and does not add extra hosts anymore :(
																ExtraHosts(masterStruct).
																Build()
	
	logging.Log.Debug(registratorServerHostConfig)
	
	registrator := benchFlowDockerUtils.CreateContainerOptionsBuilder.
											Name("registrator").
											Config(&registratorServerConfig).
											HostConfig(&registratorServerHostConfig).
											Build()


	res, err := client.CreateContainer(registrator)

	logging.Log.Debug("Err: ")
	logging.Log.Debug(err)

	if res != nil {

		logging.Log.Debug("Res: ")
		logging.Log.Debug(res)
		
		// Start Container
		logging.Log.Debug("Test Run Container")
	    
	    startRes := client.StartContainer(res.ID,nil)
	    
	    logging.Log.Debug("Res: ")
	    logging.Log.Debug(startRes)
    
    }

}
