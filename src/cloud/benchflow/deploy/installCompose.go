package deploy

import (
	benchFlowDocker "cloud/benchflow/docker"
	"cloud/benchflow/environment"
	benchFlowDockerUtils "cloud/benchflow/utils/docker"
)

//TODO: refactor

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
