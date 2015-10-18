package deploy

import (
	benchFlowDocker "cloud/benchflow/client/docker"
	benchFlowDockerUtils "cloud/benchflow/client/utils/docker"
	"cloud/benchflow/client/environment"
	"cloud/benchflow/client/utils/servers"
	"cloud/benchflow/client/structs"
	"github.com/fsouza/go-dockerclient"
	"cloud/benchflow/client/utils/logging"
	"cloud/benchflow/client/utils/configuration"
//	rancherClient "github.com/rancherio/go-rancher/client"
	"strings"
	"github.com/spf13/viper"
	"bytes"
	"io/ioutil"
	"time"
	"fmt"
	"bufio"
    "os"
)

func InstallRancher() {
	
	//BenchFlow configuration
	benchFlowConf := environment.Env.BenchFlow
	
	rancherName, rancherTag := benchFlowDockerUtils.GetImageAndTag(benchFlowConf.DockerImages.Rancher.Image)
	
	installRancher(rancherName, rancherTag)
	
}

func installRancher(rancherName string, rancherTag string) {
	
	//Deploy the Rancher server on the master server
	var master = servers.GetMasterServer(environment.Env.Servers)
	
	//Pull image on master (TODO: When swarm support image pull with affinity, change it to Swarm)
	deployRancherServer(master, rancherName, rancherTag)

	//Start rancher server on master
	rancherID := runRancher(master, rancherName, rancherTag)
	
	waitUntilRancherIsStarted(master, rancherID)
	
	//Prompt user to configure Access Control
	fmt.Println()
	fmt.Println("CONFIGURE THE ACCESS CONTROL")
	fmt.Print("-> Connect to the following url and follow the provided instructions: ")
	fmt.Println(configuration.GetServerExternalIP(master) + ":" + environment.Env.BenchFlow.DockerImages.Rancher.RunPorts[0] + "/static/admin/access/github")
	
	waitForUserEnter()
	
	//Prompt user to create BenchFlow environment
	fmt.Println("CREATE THE BenchFlow ENVIRONMENT")
	fmt.Print("-> Connect to the following url and follow the provided instructions: ")
	fmt.Println(configuration.GetServerExternalIP(master) + ":" + environment.Env.BenchFlow.DockerImages.Rancher.RunPorts[0] + "/static/settings/environments")
	fmt.Println("AFTER CREATING THE ENVIRONMENT: Switch to the newly created BenchFlow environment (Click on the arrow in the upper right corner (next to the user account image))")
	fmt.Println("NOTE: Set BenchFlow as the environment name")
	
	waitForUserEnter()
	
	//Prompt user to create the API key and secret and write here to save
	fmt.Println("CREATE AN API KEY FOR THE BenchFlow ENVIRONMENT")
	fmt.Print("-> Connect to the following url and follow the provided instructions: ")
	fmt.Println(configuration.GetServerExternalIP(master) + ":" + environment.Env.BenchFlow.DockerImages.Rancher.RunPorts[0] + "/static/settings/api")
	fmt.Println()
	fmt.Print("Paste here the 'USERNAME (ACCESS KEY)' and press Enter: ")
	username := readUserInput()
	fmt.Print("Paste here the 'PASSWORD (SECRET KEY)' and press Enter: ")
	password := readUserInput()
	
	updateConfigurationWithRancherCredentials(username,password)
	
	fmt.Println()
	fmt.Println()
	
	//Prompt user to add hosts, for each of the hosts. It also suggests the data to be used
	fmt.Println("ADD THE SERVERS TO THE BenchFlow ENVIRONMENT")
	fmt.Print("-> Connect to the following url and follow the provided instructions: ")
	fmt.Println(configuration.GetServerExternalIP(master) + ":" + environment.Env.BenchFlow.DockerImages.Rancher.RunPorts[0] + "/static/infra/hosts")
	fmt.Println()
	fmt.Println("ADD THE FOLLOWING HOSTS, AS \"Custom\" HOSTS:")
	
	var serversList = environment.Env.Servers

//	logging.Log.Debug("Servers to add to Rancher")

//	logging.Log.Debug(serversList)
	
	var server structs.Server

	//Install cassandra
	for _, server = range serversList {

		fmt.Println()
		fmt.Println("\tName: " + server.Name)
		fmt.Println("\tLabel: Key=name,Value=" + strings.ToLower(server.Name))
		
	}
	
	fmt.Println()
	fmt.Println("\tNOTE: Step (4) - remove sudo from the command")
	
	waitForUserEnter()
	
//	logging.Log.Debug("Username: " + username)
//	logging.Log.Debug("Password: " + password)

	//Deploy Rancher-Compose on Master
	rancherComposeName, rancherComposeTag := benchFlowDockerUtils.GetImageAndTag(environment.Env.BenchFlow.DockerImages.RancherCompose.Image)
	deployRancherServer(master, rancherComposeName, rancherComposeTag)

}

//TODO: change to a better way to handle yaml files updates
func updateConfigurationWithRancherCredentials(username string, password string) {
	
	//Update Configuration File
	input, err := ioutil.ReadFile(viper.ConfigFileUsed())
        if err != nil {
                logging.Log.Fatalln(err)
        }
 
        lines := strings.Split(string(input), "\n")
        
        credentialFound := false
        rancherFound := false
        usernameUpdated := false
        passwordUpdated := false
 
        for i, line := range lines {
                if !credentialFound && strings.Contains(line, "credentials:") {
//                        lines[i] = "LOL"
					fmt.Println(line)
					credentialFound = true
                } else if credentialFound && !rancherFound && strings.Contains(line, "rancher:") {
//                        lines[i] = "LOL"
					fmt.Println(line)
					rancherFound = true
                }  else if credentialFound && rancherFound && !usernameUpdated && strings.Contains(line, "username:") {
                  	lines[i] = strings.Replace(line, "\"\"", "\"" + username + "\"",1)
					fmt.Println(line)
					usernameUpdated = true
                } else if credentialFound && rancherFound && !passwordUpdated && strings.Contains(line, "password:") {
                        lines[i] = strings.Replace(line, "\"\"", "\"" + password + "\"",1)
						fmt.Println(line)
						passwordUpdated = true
                }
        }
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(viper.ConfigFileUsed(), []byte(output), 0644)
        if err != nil {
                logging.Log.Fatalln(err)
        }
	
	//Update Environment Configuration
	viper.ReadInConfig()
	

	//Validate the configuration file
	//TODO: validate config file structure
	logging.Log.Debug(viper.Get("credentials.rancher.username"))
	environment.UpdateCredentials()
	logging.Log.Debug(environment.Env.Credentials.Rancher)
	
}

func waitForUserEnter() {
	fmt.Println()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press Enter to proceed to the next step...")
	reader.ReadString('\n')
	fmt.Println()
}

func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	
	text = strings.TrimRight(text, "\n")
	
	return text
}

func deployRancherServer(server structs.Server, rancherName string, rancherTag string) {

	var client = benchFlowDocker.GetNewServerDockerClient(server)

	benchFlowDockerUtils.DeployImage(client, rancherName, rancherTag)

}

func runRancher(server structs.Server, rancherName string, rancherTag string) string {

	var client = benchFlowDocker.GetNewSwarmClient()
	
	rancherServerConfig := benchFlowDockerUtils.DockerConfigBuilder.
														Hostname("rancher_server").
														Image(rancherName,rancherTag).
														ExposedPorts(map[docker.Port]struct{}{
																"8080/tcp": {},
																"3306/tcp": {},
														}).
														//It is not yet supported and cause problem at startup
//														Volumes(map[string]struct{}{
//																environment.Env.BenchFlow.DockerImages.Rancher.DataRootDirectory + "/rancher/cattle": {},
//																environment.Env.BenchFlow.DockerImages.Rancher.DataRootDirectory + "/rancher/mysql": {},
//														}).
														AddEnv("constraint:node==" + server.Hostname).
														Build()

	logging.Log.Debug(rancherServerConfig)

	
	rancherServerHostConfig := benchFlowDockerUtils.DockerHostConfigBuilder.
																//It is not yet supported and cause problem at startup
//																AddBind(environment.Env.BenchFlow.DockerImages.Rancher.DataRootDirectory + "/rancher/cattle:/var/lib/cattle:rw").
//																AddBind(environment.Env.BenchFlow.DockerImages.Rancher.DataRootDirectory + "/rancher/mysql:/var/lib/mysql:rw").
																PortBindings(getRancherPortBindings()).
																RestartPolicy(docker.RestartPolicy{
																	Name: "always",
																}).
																Build()
	
	logging.Log.Debug(rancherServerHostConfig)
	
	rancher := benchFlowDockerUtils.CreateContainerOptionsBuilder.
											Name("rancher_server").
											Config(&rancherServerConfig).
											HostConfig(&rancherServerHostConfig).
											Build()


	res, err := client.CreateContainer(rancher)

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

//TODO: refactor to utils and make it parametrizable
func waitUntilRancherIsStarted(server structs.Server, masterContainerID string) {
	for {
		
		time.Sleep(15000 * time.Millisecond)
		
		var client = benchFlowDocker.GetNewServerDockerClient(server)
		
		var buf bytes.Buffer
	
		rancherLogs := benchFlowDockerUtils.LogOptionsBuilder.
									Container(masterContainerID).
									OutputStream(&buf).
									Follow(false).
									Stdout(true).
									Stderr(true).
									Timestamps(false).
									Build()
	
		logErr := client.Logs(rancherLogs)
	
		logging.Log.Debug(">>>>Log")
		logging.Log.Debug(buf.String())
		logging.Log.Debug(logErr)
		
		//Check for the last log message showed during the initialization
		if strings.Contains(buf.String(),"Startup Succeeded") {
			break
		}
	}
}

func getRancherPortBindings() map[docker.Port][]docker.PortBinding {
	
	//Ports
	ports := make(map[docker.Port][]docker.PortBinding, 1)
	
	//Bindings
	var bindings8080 = docker.PortBinding{
		HostPort: environment.Env.BenchFlow.DockerImages.Rancher.RunPorts[0],
	}
	
//	var bindings3306 = docker.PortBinding{
//		HostPort: environment.Env.BenchFlow.DockerImages.Rancher.RunPorts[1],
//	}
	
	ports["8080/tcp"] = []docker.PortBinding{bindings8080}
//	ports["3306/tcp"] = []docker.PortBinding{bindings3306}

	
	return ports
}