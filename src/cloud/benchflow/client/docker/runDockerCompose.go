package docker

import (
	"cloud/benchflow/client/structs"
	"cloud/benchflow/client/environment"
//	"io/ioutil"
//	"syscall"
//	"os"
	"bytes"
	benchFlowDockerUtils "cloud/benchflow/client/utils/docker"
	benchFlowSSH "cloud/benchflow/client/ssh"
	"cloud/benchflow/client/utils/servers"
	"cloud/benchflow/client/utils/configuration"
	"cloud/benchflow/client/utils/logging"
	"github.com/satori/go.uuid"
)

func RunDockerCompose(composeFile string, projectName string) {
	
	//Upload the file in the "projectName" folder
	master := servers.GetMasterServer(environment.Env.Servers)
	copyProjectFile(composeFile, master, projectName)
	
	projectFolder := environment.Env.Settings.RemoteDockerComposeProjects + "/" + projectName
	
	//BenchFlow configuration
	benchFlowConf := environment.Env.BenchFlow

	dockerComposeName, dockerComposeTag := benchFlowDockerUtils.GetImageAndTag(benchFlowConf.DockerImages.DockerCompose.Image)
	
	//Run Compose and Get the Logs. After the execution it removes the container
	runCompose(master, dockerComposeName, dockerComposeTag, projectFolder, projectName)

	//Delete the "projectName" folder
	benchFlowSSH.DeleteRemoteDirectoryInsideHome(master, projectFolder)
	
}

func copyProjectFile(composeFile string, server structs.Server, projectName string) {
	
	logging.Log.Info(environment.Env.Settings.RemoteDockerComposeProjects)
	
	benchFlowSSH.ScpFileInUserHome(server, composeFile, environment.Env.Settings.RemoteDockerComposeProjects + "/" + projectName, "docker-compose.yml")

}


//TODO, refactor
func runCompose(server structs.Server, consulName string, consulTag string, projectFolder string, projectName string) {
	
	var client = GetNewSwarmClient()
	
	var bufOut bytes.Buffer
	var bufErr bytes.Buffer

	var swarmEndPoint = "tcp://" + configuration.GetServerAddress(server) + ":" + environment.Env.BenchFlow.DockerImages.DockerSwarm.RunPorts[0]
	var home = benchFlowSSH.GetHomeDirectory(server)
	
	dockerComposeMasterConfig := benchFlowDockerUtils.DockerConfigBuilder.
																Hostname("docker_compose").
																Image(consulName,consulTag).
																Volumes(map[string]struct{}{
																		home + "/" + environment.Env.Settings.RemoteCertificatesFolder: {},
																		home + "/" + projectFolder: {},
																}).
																AddEnv("constraint:node==" + server.Hostname).
																AddEnv("DOCKER_HOST=" + swarmEndPoint).
																AddEnv("DOCKER_TLS_VERIFY=1").
																AddEnv("DOCKER_CERT_PATH=/certs").
																AddEnv("COMPOSE_PROJECT_NAME=" + projectName).
																AddCmd("up").
																AddCmd("-d").
																Build()
																
	logging.Log.Debug(dockerComposeMasterConfig)
	
	
	dockerComposeMasterHostConfig := benchFlowDockerUtils.DockerHostConfigBuilder.
															AddBind(home + "/" + environment.Env.Settings.RemoteCertificatesFolder + ":/certs:ro").
															AddBind(home + "/" + projectFolder + ":/app:ro").
															Build()
	
	logging.Log.Debug(dockerComposeMasterHostConfig)
	
	//We use a random name in order to avoid possibe conflicts of parallel run of compose
	dockerCompose := benchFlowDockerUtils.CreateContainerOptionsBuilder.
													Name("docker_compose_" + uuid.NewV4().String()).
													Config(&dockerComposeMasterConfig).
													HostConfig(&dockerComposeMasterHostConfig).
													Build()
	
	res, err := client.CreateContainer(dockerCompose)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(res.ID)

	logging.Log.Debug("Err: ")
	logging.Log.Debug(err)
	
	// Start Container
	logging.Log.Debug("Test Run Container")

	startRes := client.StartContainer(res.ID, nil)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(startRes)
	
	client = GetNewSwarmClient()
	
	//TODO, understand how to get the log of compose because for same reasons it can be obtained in the same way used for swarm
	//We need the log in order to check for the correct setup of the system
	dockerComposeLogs := benchFlowDockerUtils.LogOptionsBuilder.
													Container(res.ID).
													ErrorStream(&bufOut).
													OutputStream(&bufErr).
													Follow(true). //-d option in starting compose quits the container after completing the compose operations
													Stdout(true).
													Stderr(true).
													Timestamps(false).
													Build()
	
	logErr := client.Logs(dockerComposeLogs)

	logging.Log.Debug("Compose Out Logs")
	logging.Log.Debug(bufOut.String())
	logging.Log.Debug("Compose Err Logs")
	logging.Log.Debug(bufErr.String())
	logging.Log.Debug(logErr)
	
	//TODO: decide if we need to get the actual logs of the containers by querying docker-compose logs
	
	dockerComposeRemove := benchFlowDockerUtils.RemoveContainerOptionsBuilder.
														ID(res.ID).
														RemoveVolumes(true).
														Force(true).
														Build()
	
	client = GetNewMasterDockerClient()
	
	removeRes := client.RemoveContainer(dockerComposeRemove)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(removeRes)
	
}