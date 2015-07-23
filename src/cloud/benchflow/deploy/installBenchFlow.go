package deploy

import (
//	benchFlowDocker "cloud/benchflow/docker"
)


//TODO: make it a restartable procedure in case of error (skip the already done steps)

//TODO: test and configure the use of "always" as restart policy to enable autorestart of containers after reboot (if they correclty manage the dependencies)

func InstallBenchFlow() {
	
//	SaveIptables()

//	InstallConsul()
	
//	InstallSwarm()
	
//	InstallCassandra()

//	InstallRancher()

//	InstallCompose()

//	defaultConfig := `mywordpress:
// tty: true
// image: wordpress
// stdin_open: true
// environment:
//  - "constraint:node==grid"`

	//TODO, remove: test docker compose
//	benchFlowDocker.RunDockerCompose(defaultConfig,"mywordpress")


    
}