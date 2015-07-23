package deploy

import (
//	benchFlowDocker "cloud/benchflow/docker"
//	"github.com/fsouza/go-dockerclient"
)

//TODO: refactor

func UninstallBenchFlow() {
	
//	UninstallKafka()
	
//	UninstallRancher()
	
//	UninstallCassandra()
	
	UninstallSwarm()
	
	//TODO: Remove, testing code
	
//	var masterClient = benchFlowDocker.GetNewMasterDockerClient()
//	
//	res, err := masterClient.Info()
//    
//    fmt.Println("Res: ")
//    fmt.Println(res)
//    fmt.Println("Err: ")
//    fmt.Println(err)
//
//	fmt.Println("Images")
//    imgs, _ := masterClient.ListImages(docker.ListImagesOptions{All: false})
//    for _, img := range imgs {
//        fmt.Println("ID: ", img.ID)
//        fmt.Println("RepoTags: ", img.RepoTags)
//        fmt.Println("Created: ", img.Created)
//        fmt.Println("Size: ", img.Size)
//        fmt.Println("VirtualSize: ", img.VirtualSize)
//        fmt.Println("ParentId: ", img.ParentID)
//    }

//	UninstallConsul()

//	RestoreIptables()
    
}