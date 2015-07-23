package deploy

import (
//    "strings"
//    "bytes"
    "cloud/benchflow/environment"
    "cloud/benchflow/structs"
    "cloud/benchflow/utils/servers"
    "cloud/benchflow/utils/logging"
    benchFlowSSH "cloud/benchflow/ssh"
	benchFlowDocker "cloud/benchflow/docker"
	benchFlowDockerUtils "cloud/benchflow/utils/docker"
	"github.com/fsouza/go-dockerclient"
	//Test Rancher
	rancherClient "github.com/rancherio/go-rancher/client"
)

//TODO: refactor

//TODO, remove: testing Rancher
const (
	PROJECT_URL = "http://195.176.181.55:8081/v1/projects/1a13"
	URL         = "http://195.176.181.55:8081/v1/"
//	ACCESS_KEY  = "admin"
//	SECRET_KEY  = "adminpass"
//	MAX_WAIT    = time.Duration(time.Second * 10)
)

func newClient(url string) *rancherClient.RancherClient {
	client, err := rancherClient.NewRancherClient(&rancherClient.ClientOpts{
		Url:       url,
//		AccessKey: ACCESS_KEY,
//		SecretKey: SECRET_KEY,
	})

	if err != nil {
		logging.Log.Fatal("Failed to create client", err)
	}

	return client
}

func TestClientLoad() {
	client := newClient(URL)
	if client.Schemas == nil {
		logging.Log.Debug("Failed to load schema")
	}

	if len(client.Schemas.Data) == 0 {
		logging.Log.Debug("Schemas is empty")
	}

	if _, ok := client.Types["container"]; !ok {
		logging.Log.Debug("Failed to find container type")
	}
	
	logging.Log.Debug()
	logging.Log.Debug()
	logging.Log.Debug(client.ApiKey.Create(&rancherClient.ApiKey{AccountId: "1a5",Name: "Remote",}))
	logging.Log.Debug()
	logging.Log.Debug()
	logging.Log.Debug(client.Environment.List(&rancherClient.ListOpts{}))
}

func testRancherAPI() {
	
	TestClientLoad()
	
}

func UninstallSwarm() {
	
//	testRancherAPI()
	
	 undeploySwarm()

	//TODO: move, test contact swarm
//	var masterClient = benchFlowDocker.GetNewMasterDockerClient()
//	var masterClient = benchFlowDocker.GetNewSwarmClient()
	
//	res, err := masterClient.Info()
//    
//    logging.Log.Debug("Res: ")
//    logging.Log.Debug(res)
//    logging.Log.Debug("Err: ")
//    logging.Log.Debug(err)
//
//	logging.Log.Debug("Images")
//    imgs, _ := masterClient.ListImages(docker.ListImagesOptions{All: false})
//    for _, img := range imgs {
//        logging.Log.Debug("ID: ", img.ID)
//        logging.Log.Debug("RepoTags: ", img.RepoTags)
//        logging.Log.Debug("Created: ", img.Created)
//        logging.Log.Debug("Size: ", img.Size)
//        logging.Log.Debug("VirtualSize: ", img.VirtualSize)
//        logging.Log.Debug("ParentId: ", img.ParentID)
//    }

//	logging.Log.Debug("Test Import Image")
//	
//    var buf bytes.Buffer
//    
//    swarm := docker.PullImageOptions{
//    			Repository: "swarm",
//    			Tag: "0.3.0",
//    			OutputStream: &buf,
//    		 }
//    
//    res := masterClient.PullImage(swarm, docker.AuthConfiguration{})
//    
//    logging.Log.Debug("Output: " + buf.String()) //Output of the command
//    
//    logging.Log.Debug("Res: ")
//    logging.Log.Debug(res)
		
//		//Test create container with Affinity
//
//		swarmJoinConfig := docker.Config{
//				    			Hostname: "swarm_test",
//				    			AttachStdout: true,
//				    			Image: "swarm:0.3.0",
//				    			Env: []string{"constraint:node==neha"},
//				    			Cmd: []string{"list"},
//				    		 }
//		
//	    swarmJoin := docker.CreateContainerOptions{
//		    			Name: "swarm_test",
//		    			Config: &swarmJoinConfig,
//		    		 }
//	    
//	    res,err := masterClient.CreateContainer(swarmJoin)
//    
//	    logging.Log.Debug("Res Join Create: ")
//	    logging.Log.Debug(res)
//	    
//	    logging.Log.Debug("Err Join Create: ")
//	    logging.Log.Debug(err)
//		
//		resC,errC := masterClient.ListContainers(docker.ListContainersOptions{All: true})
//		
//		logging.Log.Debug("Res List Containers: ")
//	    logging.Log.Debug(resC)
//	    
//	    logging.Log.Debug("Err List Containers: ")
//	    logging.Log.Debug(errC)

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