package docker

import (
	"os"
	"fmt"
	"cloud/benchflow/environment"
	"cloud/benchflow/structs"
	"cloud/benchflow/utils/logging"
    "cloud/benchflow/utils/configuration"
    "cloud/benchflow/utils/servers"
    "github.com/fsouza/go-dockerclient"
)

func GetNewSwarmClient() *docker.Client {
	
	//Master server: where we deploy the swarm master
    master := servers.GetMasterServer(environment.Env.Servers)
    
    var serverEndPoint = "tcp://" + configuration.GetServerAddress(master) + ":" + environment.Env.BenchFlow.DockerImages.DockerSwarm.RunPorts[0]
    
    return buildNewClient(serverEndPoint)
	
}

//TODO: think about keep only the GetNewServerDockerClient one
func GetNewMasterDockerClient() *docker.Client {
	
	//Master server
    master := servers.GetMasterServer(environment.Env.Servers)
    
    var serverEndPoint = "tcp://" + configuration.GetServerEndPoint(master)
    
    return buildNewClient(serverEndPoint)
	
}

func GetNewServerDockerClient(server structs.Server) *docker.Client {
	
    var serverEndPoint = "tcp://" + configuration.GetServerEndPoint(server)
    
    return buildNewClient(serverEndPoint)
	
}

func buildNewClient(serverEndPoint string) *docker.Client {
    
//    var serverEndPoint = "tcp://" + configuration.GetServerEndPoint(server)
//    var serverName = strings.ToLower(server.Name)
    
    logging.Log.Debug("Server Endpoint: " + serverEndPoint)
    
    //The certificate must be valid for the selected NAME
    path := environment.Env.Settings.LocalCertificateFolder
//    fmt.Println(path)
//    fmt.Println(fmt.Sprintf("%s" + string(os.PathSeparator) + "local" + string(os.PathSeparator) + "cert.pem", path))
    ca := fmt.Sprintf("%s" + "local" + string(os.PathSeparator) + "ca.pem", path)
    cert := fmt.Sprintf("%s" + "local" + string(os.PathSeparator) + "cert.pem", path)
    key := fmt.Sprintf("%s" + "local" + string(os.PathSeparator) + "key.pem", path)
    serverClient, err := docker.NewTLSClient(serverEndPoint, cert, key, ca)
    
//    res, err := serverClient.Info()
//    
//    fmt.Println("Res: ")
//    fmt.Println(res)
    logging.Log.Debug("Err: ")
    logging.Log.Debug(err)
    
    return serverClient

}