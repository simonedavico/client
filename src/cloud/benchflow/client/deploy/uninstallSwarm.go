package deploy

import (
//    "strings"
//    "bytes"
    "cloud/benchflow/client/environment"
    "cloud/benchflow/client/structs"
    "cloud/benchflow/client/utils/servers"
    "cloud/benchflow/client/utils/logging"
    benchFlowSSH "cloud/benchflow/client/ssh"
	benchFlowDocker "cloud/benchflow/client/docker"
	benchFlowDockerUtils "cloud/benchflow/client/utils/docker"
	"github.com/fsouza/go-dockerclient"
)

func UninstallSwarm() {
	
	 undeploySwarm()
}

//TODO: update to undeploy to all the nodes, not only on master
func undeploySwarm() {
	
	var client = benchFlowDocker.GetNewMasterDockerClient()
	
	//Docker configuration
    benchFlowConf := environment.Env.BenchFlow
	
	logging.Log.Debug("Test Remove Containers")
	
	swarmContainersConf := benchFlowDockerUtils.ListContainerOptionsBuilder.
														All(true).
														Filters(makeSwarmFilter()).
														Build()
	
//	swarmContainersConf := docker.ListContainersOptions{
//							All: true,
//							Filters: makeSwarmFilter(),
//						}
	
	//Remove all the started containers using the swarm image
	res, err := client.ListContainers(swarmContainersConf)
	
	var cont docker.APIContainers
	
	for _, cont = range res {
		
		swarmRemoveConf := benchFlowDockerUtils.RemoveContainerOptionsBuilder.
														ID(cont.ID).
														RemoveVolumes(true).
														Force(true).
														Build()
		
//		swarmRemoveConf := docker.RemoveContainerOptions{
//							ID: cont.ID,
//							RemoveVolumes: true,
//							Force: true,
//						}

		client.RemoveContainer(swarmRemoveConf)
	}
	
	logging.Log.Debug(res)
	logging.Log.Debug(err)
	
	logging.Log.Debug("Test Remove Image")
	
	client = benchFlowDocker.GetNewMasterDockerClient()

    //Remove the image
    errRemove := client.RemoveImage(benchFlowConf.DockerImages.DockerSwarm.Image)
    
    if errRemove != nil {
		logging.Log.Debug(errRemove)
	}
    
    //Remove certificates
    var serversMap map[string]structs.Server = environment.Env.Servers
	
	var master structs.Server = servers.GetMasterServer(serversMap)
	
	benchFlowSSH.DeleteRemoteDirectoryInsideHome(master,environment.Env.Settings.RemoteCertificatesFolder)
    
}

func makeSwarmFilter() map[string][]string {
	
	filter := make(map[string][]string,1)
	
	filter["name"] = []string{"swarm"}
	
	return filter
}