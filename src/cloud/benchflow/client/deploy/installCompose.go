package deploy

import (
	benchFlowDocker "cloud/benchflow/client/docker"
	"cloud/benchflow/client/environment"
	benchFlowDockerUtils "cloud/benchflow/client/utils/docker"
)

//TODO: refactor
//TODO: for now I don't use it --> remove/comment out (maybe false, I use it to deploy some compoments, also benchflow components)

func InstallCompose() {

	//BenchFlow configuration
	benchFlowConf := environment.Env.BenchFlow

	composeName, composeTag := benchFlowDockerUtils.GetImageAndTag(benchFlowConf.DockerImages.DockerCompose.Image)

	deployCompose(composeName, composeTag)

}

func deployCompose(composeName string, composeTag string) {
	
	var client = benchFlowDocker.GetNewMasterDockerClient()
	
	benchFlowDockerUtils.DeployImage(client, composeName, composeTag)

}
