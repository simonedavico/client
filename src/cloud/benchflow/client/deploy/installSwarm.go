package deploy

import (
	"bytes"
	benchFlowDocker "cloud/benchflow/client/docker"
	"cloud/benchflow/client/environment"
	benchFlowSSH "cloud/benchflow/client/ssh"
	"cloud/benchflow/client/structs"
	"cloud/benchflow/client/utils/configuration"
	benchFlowDockerUtils "cloud/benchflow/client/utils/docker"
	"cloud/benchflow/client/utils/logging"
	"cloud/benchflow/client/utils/servers"
	"github.com/fsouza/go-dockerclient"
	"os"
	"strings"
)

//TODO: refactor

//TODO: reduce heartbeat time by setting the corresponding flag: https://github.com/docker/swarm/blob/28cd51b65f1e721d073ee69befd72dd09467d54b/cli/flags.go

func InstallSwarm() {

	//BenchFlow configuration
	benchFlowConf := environment.Env.BenchFlow

	swarmName, swarmTag := benchFlowDockerUtils.GetImageAndTag(benchFlowConf.DockerImages.DockerSwarm.Image)

	deploySwarm(swarmName, swarmTag)
	runSwarm(swarmName, swarmTag)

}

func deploySwarm(swarmName string, swarmTag string) {
	
	//Deploy swarm on all servers
	var server structs.Server

	//Install cassandra
	for _, server = range environment.Env.Servers {
		
		var client = benchFlowDocker.GetNewServerDockerClient(server)
		
		benchFlowDockerUtils.DeployImage(client, swarmName, swarmTag)

	}
	

}

func runSwarm(swarmName string, swarmTag string) {

	logging.Log.Debug("Test Start Swarm")

	clusterID := createCluster(swarmName, swarmTag)

	joinNodes(swarmName, swarmTag, clusterID)

	copyCertsOnMasterServer()

	startSwarmMaster(swarmName, swarmTag, clusterID)

}

func createCluster(swarmName string, swarmTag string) string {

	var client = benchFlowDocker.GetNewMasterDockerClient()

	swarmConfig := benchFlowDockerUtils.DockerConfigBuilder.
		Hostname("smarm").
		Image(swarmName, swarmTag).
		AddCmd("create").
		Build()

	swarm := benchFlowDockerUtils.CreateContainerOptionsBuilder.
		Name("swarm").
		Config(&swarmConfig).
		Build()

	res, err := client.CreateContainer(swarm)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(res.ID)

	logging.Log.Debug("Err: ")
	logging.Log.Debug(err)

	logging.Log.Debug("Test Run Container")

	startRes := client.StartContainer(res.ID, nil)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(startRes)

	client = benchFlowDocker.GetNewMasterDockerClient()

	//Get the Cluster ID from the container log
	var buf bytes.Buffer
	//TODO, check also for the ErrorStream in case of errors
	swarmLogs := benchFlowDockerUtils.LogOptionsBuilder.
		Container(res.ID).
		OutputStream(&buf).
		Follow(true).
		Stdout(true).
		Stderr(true).
		Timestamps(false).
		Build()

	logErr := client.Logs(swarmLogs)

	logging.Log.Debug("Cluster ID")
	logging.Log.Debug(buf.String())
	logging.Log.Debug(logErr)

	clusterID := buf.String()

	clusterID = strings.TrimRight(clusterID, "\n")

	client = benchFlowDocker.GetNewMasterDockerClient()

	//    logging.Log.Debug("Rest Remove Container")

	swarmRemove := benchFlowDockerUtils.RemoveContainerOptionsBuilder.
		ID(res.ID).
		RemoveVolumes(true).
		Force(true).
		Build()

	removeRes := client.RemoveContainer(swarmRemove)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(removeRes)

	return clusterID

}

func joinNodes(swarmName string, swarmTag string, clusterID string) {

	//Servers
	servers := environment.Env.Servers

	logging.Log.Debug("Servers to Join")

	logging.Log.Debug(servers)

	var server structs.Server

	//Join all the servers to the cluster
	for _, server = range servers {

		var client = benchFlowDocker.GetNewServerDockerClient(server)

		var serverEndPoint = configuration.GetServerEndPoint(server)

		swarmJoinConfig := benchFlowDockerUtils.DockerConfigBuilder.
			Hostname("swarm_"+strings.ToLower(server.Name)).
			Image(swarmName, swarmTag).
			AddCmd("join").
			AddCmd("--addr=" + serverEndPoint).
			AddCmd("token://" + clusterID).
			Build()

		//		swarmJoinConfig := docker.Config{
		//				    			Hostname: "swarm_" + strings.ToLower(server.Name),
		////				    			AttachStdout: true,
		//				    			Image: swarmName + ":" + swarmTag,
		//				    			Cmd: []string{"join","--addr=" + serverEndPoint,"token://" + clusterID},
		//				    		 }

		swarmJoin := benchFlowDockerUtils.CreateContainerOptionsBuilder.
			Name("swarm_" + strings.ToLower(server.Name)).
			Config(&swarmJoinConfig).
			Build()

		//	    swarmJoin := docker.CreateContainerOptions{
		//		    			Name: "swarm_" + strings.ToLower(server.Name),
		//		    			Config: &swarmJoinConfig,
		//		    		 }

		res, err := client.CreateContainer(swarmJoin)

		logging.Log.Debug("Err Join Create: ")
		logging.Log.Debug(err)

		logging.Log.Debug("Res Join Create: ")
		logging.Log.Debug(res.ID)

		logging.Log.Debug("Test Run Join Container")

		startRes := client.StartContainer(res.ID, nil)

		logging.Log.Debug("Res Join Start: ")
		logging.Log.Debug(startRes)

	}

	logging.Log.Debug(servers)

}

func copyCertsOnMasterServer() {

	var serversMap map[string]structs.Server = environment.Env.Servers

	var master structs.Server = servers.GetMasterServer(serversMap)

	var folderPathFrom string = environment.Env.Settings.LocalCertificateFolder + strings.ToLower(master.Name) + string(os.PathSeparator)

	benchFlowSSH.ScpFolderInUserHome(master, folderPathFrom, environment.Env.Settings.RemoteCertificatesFolder)

}

func startSwarmMaster(swarmName string, swarmTag string, clusterID string) {

	var client = benchFlowDocker.GetNewMasterDockerClient()

	master := servers.GetMasterServer(environment.Env.Servers)

	var home = benchFlowSSH.GetHomeDirectory(master)

	//		"AttachStdout": true,
	//	var JSONswarmMasterConfig string = `{
	//		"Hostname": "swarm_master",
	//		"Image": "` + swarmName + ":" + swarmTag + `",
	//		"ExposedPorts": {
	//			       	 "2375/tcp": {}
	//			      },
	//		"Volumes": {
	//			         "` + home + "/" + environment.Env.Settings.RemoteCertificatesFolder + `/": {}
	//			      },
	//		"Cmd": [
	//			"manage",
	//			"--tlsverify",
	//			"--tlscacert=/certs/ca.pem",
	//			"--tlscert=/certs/cert.pem",
	//			"--tlskey=/certs/key.pem",
	//			"token://` + clusterID + `"
	//		]
	//	}`
	//	JSONswarmMasterConfig = JSONswarmMasterConfig + `
	//	}`

	//	logging.Log.Debug(JSONswarmMasterConfig)

	swarmMasterConfig := benchFlowDockerUtils.DockerConfigBuilder.
		Hostname("swarm_master").
		Image(swarmName, swarmTag).
		ExposedPorts(map[docker.Port]struct{}{"2375/tcp": {}}).
		Volumes(map[string]struct{}{home + environment.Env.Settings.RemoteCertificatesFolder: {}}).
		AddCmd("manage").
		AddCmd("--tlsverify").
		AddCmd("--tlscacert=/certs/ca.pem").
		AddCmd("--tlscert=/certs/cert.pem").
		AddCmd("--tlskey=/certs/key.pem").
		AddCmd("token://" + clusterID).
		Build()

	logging.Log.Debug(swarmMasterConfig)

	//	var swarmMasterConfig docker.Config
	//	errConfig := json.Unmarshal([]byte(JSONswarmMasterConfig), &swarmMasterConfig)

	//	logging.Log.Debug(swarmMasterConfig)
	//	logging.Log.Debug("Errr: ")
	//	logging.Log.Debug(errConfig)
	
	joinedServers := servers.GetNotMasterServers(environment.Env.Servers)

	swarmMasterHostConfig := benchFlowDockerUtils.DockerHostConfigBuilder.
		AddBind(home + "/" + environment.Env.Settings.RemoteCertificatesFolder + "/:/certs/:ro").
		PortBindings(getSwarmMasterPortBindings()).
		ExtraHosts(joinedServers).
		Build()

//	joinedServers := servers.GetNotMasterServers(environment.Env.Servers)
//	var listJoinedServers []string
//
//	var server structs.Server
//
//	for _, server = range joinedServers {
//
//		//We need to define additional host only if we use aliases
//		if server.Alias != "" {
//			//Get the utilized IP
//			var serverIP = configuration.GetServerIP(server)
//			listJoinedServers = append(listJoinedServers, server.Alias+":"+serverIP)
//		}
//
//	}
//
//	var JSONswarmMasterHostConfig string = `{
//			"Binds": ["` + home + "/" + environment.Env.Settings.RemoteCertificatesFolder + `/:/certs/:ro"],
//			"PortBindings": { 
//				"2375/tcp": [
//						{ "HostPort":"` + environment.Env.BenchFlow.DockerImages.DockerSwarm.RunPorts[0] + `" }
//					] 
//			},
//			"ExtraHosts": [`
//
//	for _, addr := range listJoinedServers {
//		JSONswarmMasterHostConfig = JSONswarmMasterHostConfig + "\"" + addr + "\","
//	}
//
//	JSONswarmMasterHostConfig = strings.TrimRight(JSONswarmMasterHostConfig, ",")
//
//	JSONswarmMasterHostConfig = JSONswarmMasterHostConfig + `]
//	}`
//
//	logging.Log.Debug(JSONswarmMasterHostConfig)
//
//	var swarmMasterHostConfig docker.HostConfig
//	errHostConfig := json.Unmarshal([]byte(JSONswarmMasterHostConfig), &swarmMasterHostConfig)

	logging.Log.Debug(swarmMasterHostConfig)
//	logging.Log.Debug("Errr: ")
//	logging.Log.Debug(errHostConfig)

	swarmMaster := benchFlowDockerUtils.CreateContainerOptionsBuilder.
											Name("swarm_master").
											Config(&swarmMasterConfig).
											HostConfig(&swarmMasterHostConfig).
											Build()


//	swarmMaster := docker.CreateContainerOptions{
//		Name:       "swarm_master",
//		Config:     &swarmMasterConfig,
//		HostConfig: &swarmMasterHostConfig,
//	}

	logging.Log.Debug("Final Config: ")
	logging.Log.Debug(swarmMaster)

	res, err := client.CreateContainer(swarmMaster)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(res)

	logging.Log.Debug("Err: ")
	logging.Log.Debug(err)

	logging.Log.Debug("Test Run Container")

	startRes := client.StartContainer(res.ID, nil)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(startRes)

}

func getSwarmMasterPortBindings() map[docker.Port][]docker.PortBinding {
	
	//Ports
	ports := make(map[docker.Port][]docker.PortBinding, 2)
	
	//Bindings
	var bindings2375 = docker.PortBinding{
		HostPort: environment.Env.BenchFlow.DockerImages.DockerSwarm.RunPorts[0],
	}
	
	ports["2375/tcp"] = []docker.PortBinding{bindings2375}
	
	return ports
}
