package deploy

import (
	benchFlowDocker "cloud/benchflow/docker"
	"cloud/benchflow/environment"
	"cloud/benchflow/structs"
	"cloud/benchflow/utils/logging"
	"cloud/benchflow/utils/servers"
	"cloud/benchflow/utils/configuration"
	benchFlowDockerUtils "cloud/benchflow/utils/docker"
	"github.com/fsouza/go-dockerclient"
	"github.com/gocql/gocql"
	"strings"
	"strconv"
	"time"
	"bytes"
//	"fmt"
)

//TODO: refactor

//Tested correct presence of the data on different nodes (but not optimized)

func InstallCassandra() {

	//TODO: refactor to a function splitting the image to Name and Tag and reuse in the project

	//BenchFlow configuration
//	benchFlowConf := environment.Env.BenchFlow

//	cassandraName, cassandraTag := benchFlowDockerUtils.GetImageAndTag(benchFlowConf.DockerImages.Cassandra.Image)

//	installCassandra(cassandraName, cassandraTag)
	
	configureCassandraSuperuser()
	
	//TODO: remove - Test Connection
//	testConnection()

}

func configureCassandraSuperuser() {
	
	//TODO: verify according to: http://docs.datastax.com/en/cassandra/2.0/cassandra/security/secure_config_native_authorize_t.html

	hosts := getCassandraHosts()

	session := connectCluster(hosts, "cassandra", "cassandra", "system_auth")

    session.Query(`ALTER KEYSPACE system_auth WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : ` + strconv.Itoa(len(hosts)) + `}`).Exec()
	
	//Change credentials
	
	logging.Log.Info(`CREATE USER IF NOT EXISTS ` + environment.Env.Credentials.Cassandra.Username + ` WITH PASSWORD '` + environment.Env.Credentials.Cassandra.Password + `' SUPERUSER`)
	
	session.Query(`CREATE USER IF NOT EXISTS ` + environment.Env.Credentials.Cassandra.Username + ` WITH PASSWORD '` + environment.Env.Credentials.Cassandra.Password + `' SUPERUSER`).Exec()
	
	disconnectCluster(session)
	
	sessionNewUser := connectCluster(hosts, environment.Env.Credentials.Cassandra.Username, environment.Env.Credentials.Cassandra.Password, "system_auth")
	
	sessionNewUser.Query(`ALTER USER cassandra WITH PASSWORD 'VriHThalDBrOwF2wzvaMcVW4H8P9XUMvuzf8k3OCoWtpAfMgByL6D1QzvJCKYrL8BkhIZuaPZKQs2RFfek0mmg8MKpe8N0ZLtqAK' NOSUPERUSER`).Exec()
	
	disconnectCluster(sessionNewUser)
	
}

func connectCluster(hosts []string, username string, password string, keyspace string) *gocql.Session {
	
	cluster := gocql.NewCluster() //TODO: change to dinamically determined
	cluster.Authenticator = gocql.PasswordAuthenticator {
		Username: username,
		Password: password,
	}
	cluster.Hosts = hosts
	cluster.Timeout = time.Duration(1 * time.Minute)
	cluster.Port, _ = strconv.Atoi(environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[3])
	cluster.Keyspace = keyspace
    cluster.Consistency = gocql.Quorum
    session, _ := cluster.CreateSession()
	
	return session
}

func disconnectCluster(session *gocql.Session) {
	
	session.Close()
	
}

func configureCassandraReplicationFactorForSystemAuth() {
	
	
	
}

func installCassandra(cassandraName string, cassandraTag string) {

//	//Deploy the Cassandra master on the master server
	var master = servers.GetMasterServer(environment.Env.Servers)
//	
//	//The master is the Cassandra Seed (No need to define seed)
//	//Pull image on master (TODO: When swarm support image pull with affinity, change it to Swarm)
	deployCassandraServer(master, cassandraName, cassandraTag)
//
//	//Start cassandra on master
	masterContainerID := runCassandra(master, cassandraName, cassandraTag, "master")
//	
//	//Wait until cassandra is correctly started on the master
	waitUntilTheMasterIsStarted(master, masterContainerID)
	
	//-------------------------//
	
	
	//Deploy a clustered Cassandra, with a node on every not SUT, DRIVERS, MASTER server
	//Servers
	var serversList = servers.GetNotSutNotDriversNotMasterServers(environment.Env.Servers)

	logging.Log.Debug("Servers on which to install Cassandra")

	logging.Log.Debug(serversList)
	
	var server structs.Server

	//Install cassandra
	for _, server = range serversList {

		//Pull image on server (TODO: When swarm support image pull with affinity, change it to Swarm)
		deployCassandraServer(server, cassandraName, cassandraTag)

		//Start cassandra on server
		runCassandra(server, cassandraName, cassandraTag, "node")

	}

}

func deployCassandraServer(server structs.Server, cassandraName string, cassandraTag string) {

	var client = benchFlowDocker.GetNewServerDockerClient(server)

	benchFlowDockerUtils.DeployImage(client, cassandraName, cassandraTag)

}

//TODO: refactor to a parametrized method, using the builder design pattern

func runCassandra(server structs.Server, cassandraName string, cassandraTag string, role string) string {

	var client = benchFlowDocker.GetNewSwarmClient()
	
	var serverIP = configuration.GetServerIP(server)
	var masterIP = configuration.GetServerIP(servers.GetMasterServer(environment.Env.Servers))
	
	cassandraMasterConfigToBuild := benchFlowDockerUtils.DockerConfigBuilder.
														Hostname("cassandra_" + role).
														Image(cassandraName,cassandraTag).
														ExposedPorts(map[docker.Port]struct{}{
																"7000/tcp": {},
																"7001/tcp": {},
																"7199/tcp": {},
																"9042/tcp": {},
																"9160/tcp": {},
														}).
														Volumes(map[string]struct{}{
																environment.Env.BenchFlow.DockerImages.Cassandra.DataRootDirectory + "/cassandra": {},
																environment.Env.BenchFlow.DockerImages.Cassandra.DataRootDirectory + "/cassandra/data": {},
														}).
														AddEnv("constraint:node==" + server.Hostname).
														AddEnv("CASSANDRA_CLUSTER_NAME=benchflow").
														AddEnv("CASSANDRA_BROADCAST_ADDRESS=" + serverIP)
														
														
	if !strings.EqualFold(role,"master") {
			cassandraMasterConfigToBuild = cassandraMasterConfigToBuild.AddEnv("CASSANDRA_SEEDS=" + masterIP)
	}
	
	cassandraMasterConfig := cassandraMasterConfigToBuild.Build()
														
	//		"AttachStdout": true,
//	var JSONcassandraMasterConfig string = `{
//		"Hostname": "cassandra_` + role + `",
//		"Image": "` + cassandraName + ":" + cassandraTag + `",
//		"ExposedPorts": {
//			       	 "7000/tcp": {},
//			       	 "7001/tcp": {},
//			       	 "7199/tcp": {},
//			       	 "9042/tcp": {},
//			       	 "9160/tcp": {}
//			      },
//		"Volumes": {
//			         "` + environment.Env.BenchFlow.DockerImages.Cassandra.DataRootDirectory + `/cassandra": {},
//			         "` + environment.Env.BenchFlow.DockerImages.Cassandra.DataRootDirectory + `/cassandra/data": {}
//			      },
//		"Env": ["constraint:node==` + server.Hostname + `","CASSANDRA_CLUSTER_NAME=benchflow","CASSANDRA_BROADCAST_ADDRESS=` + masterIP + `"`
//			
//			if !strings.EqualFold(role,"master") {
//				JSONcassandraMasterConfig = JSONcassandraMasterConfig + ",\"CASSANDRA_SEEDS=" + serverIP + "\""
//			}
//		
//		JSONcassandraMasterConfig = JSONcassandraMasterConfig + `]
//	}`

//	var cassandraMasterConfig docker.Config
//	errConfigMaster := json.Unmarshal([]byte(JSONcassandraMasterConfig), &cassandraMasterConfig)
	
	
//	logging.Log.Debug(JSONcassandraMasterConfig)	
	
	logging.Log.Debug(cassandraMasterConfig)
//	logging.Log.Debug("Errr: ")
//	logging.Log.Debug(errConfigMaster)
	
//	, 
//				"7001/tcp": [
//						{ "HostPort":"` + strconv.Itoa(environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[1]) + `" }
//					],
//				"7199/tcp": [
//						{ "HostPort":"` + strconv.Itoa(environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[2]) + `" }
//					],
//				"9042/tcp": [
//						{ "HostPort":"` + strconv.Itoa(environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[3]) + `" }
//					],
//				"9160/tcp": [
//						{ "HostPort":"` + strconv.Itoa(environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[4]) + `" }
//					]
//

//,
//			"ExtraHosts": [` + "\"" + server.Alias + ":" + serverIP + "\"" + `]

	//TODO: open 9042 only on master
	
	cassandraMasterHostConfig := benchFlowDockerUtils.DockerHostConfigBuilder.
																AddBind(environment.Env.BenchFlow.DockerImages.Cassandra.DataRootDirectory + "/cassandra:/var/lib/cassandra:rw").
																AddBind(environment.Env.BenchFlow.DockerImages.Cassandra.DataRootDirectory + "/cassandra/data:/var/lib/cassandra/data:rw").
																PortBindings(getCassandraPortBindings()).
																Build()
	
//	var JSONcassandraMasterHostConfig string = `{
//			"Binds": [
//				"` + environment.Env.BenchFlow.DockerImages.Cassandra.DataRootDirectory + `/cassandra:/var/lib/cassandra:rw",
//				"` + environment.Env.BenchFlow.DockerImages.Cassandra.DataRootDirectory + `/cassandra/data:/var/lib/cassandra/data:rw"
//			],
//			"PortBindings": { 
//				"7000/tcp": [
//						{ "HostPort":"` + environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[0] + `" }
//					],
//				"9042/tcp": [ 
//						{ "HostPort":"` + environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[3] + `" }
//					]
//			}
//	}`
//	
//	var cassandraMasterHostConfig docker.HostConfig
//	errHostConfig := json.Unmarshal([]byte(JSONcassandraMasterHostConfig), &cassandraMasterHostConfig)
	
//	logging.Log.Debug(JSONcassandraMasterHostConfig)	
	
	logging.Log.Debug(cassandraMasterHostConfig)
//	logging.Log.Debug("Errr: ")
//	logging.Log.Debug(errHostConfig)
	
	cassandra := benchFlowDockerUtils.CreateContainerOptionsBuilder.
											Name("cassandra_" + role).
											Config(&cassandraMasterConfig).
											HostConfig(&cassandraMasterHostConfig).
											Build()

//	cassandra := docker.CreateContainerOptions{
//		Name:   "cassandra_" + role,
//		Config: &cassandraMasterConfig,
//		HostConfig: &cassandraMasterHostConfig,
//	}

	res, err := client.CreateContainer(cassandra)

	logging.Log.Debug("Res: ")
	logging.Log.Debug(res)

	logging.Log.Debug("Err: ")
	logging.Log.Debug(err)
	
	// Start Container
	logging.Log.Debug("Test Run Container")
    
    startRes := client.StartContainer(res.ID,nil)
    
    logging.Log.Debug("Res: ")
    logging.Log.Debug(startRes)
    
    return res.ID

}

func getCassandraPortBindings() map[docker.Port][]docker.PortBinding {
	
	//Ports
	ports := make(map[docker.Port][]docker.PortBinding, 2)
	
	//Bindings
	var bindings7000 = docker.PortBinding{
		HostPort: environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[0],
	}
	
	var bindings9042 = docker.PortBinding{
		HostPort: environment.Env.BenchFlow.DockerImages.Cassandra.RunPorts[3],
	}
	
	ports["7000/tcp"] = []docker.PortBinding{bindings7000}
	ports["9042/tcp"] = []docker.PortBinding{bindings9042}

	
	return ports
}

func waitUntilTheMasterIsStarted(master structs.Server, masterContainerID string) {
	for {
		
		time.Sleep(15000 * time.Millisecond)
		
		var client = benchFlowDocker.GetNewServerDockerClient(master)
		
		var buf bytes.Buffer
	
		cassandraLogs := benchFlowDockerUtils.LogOptionsBuilder.
									Container(masterContainerID).
									OutputStream(&buf).
									Follow(false).
									Stdout(true).
									Stderr(true).
									Timestamps(false).
									Build()
	
		logErr := client.Logs(cassandraLogs)
	
		logging.Log.Debug(">>>>Log")
		logging.Log.Debug(buf.String())
		logging.Log.Debug(logErr)
		
		//Check for the last log message showed during the initialization
		if strings.Contains(buf.String(),"Created default superuser") {
			break
		}
	}
}

func getCassandraHosts() []string {
	var nodes = servers.GetNotSutNotDriversNotMasterServers(environment.Env.Servers)
	var hosts []string
	
	var server structs.Server
	
	for _, server = range nodes {

		hosts = append(hosts,configuration.GetServerAddress(server))

	}
	
	return hosts
}