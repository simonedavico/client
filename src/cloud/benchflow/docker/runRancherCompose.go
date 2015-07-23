package docker

import (

)

//TODO: think more why I also need docker-compose pointing at swarm and rancher-compose it is not enough
//Useful to user rancher services, but not for deploying SUT because add network that I do not need and I want a more fine grain control and define extra properties out of the docker-compose file so that it can be reused (that one compose pointing at swarm)

//docker run -v /Users/vincenzofe
//rme/Dropbox/Backup/Ph.D./Git/BenchFlow-Project/benchflow-docker-images/rancher-compose/examples/wordpress:/app:ro -e "RANCHER_URL=http://195.176.181.55:8085/v1" -e "RANCHER_A
//CCESS_KEY=5C23235C01F9C633AD2C" -e "RANCHER_SECRET_KEY=QgawkG8i87NY9j4KZPGAKUXJjvvrkjen9V2zVWh7" --rm 0eabb1c62abe -f "docker-compose.yml" -r "rancher-compose.yml" -p "wordpr
//ess" logs

//Move file, start and check started, remove file and folder