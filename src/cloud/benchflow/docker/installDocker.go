package docker
//
//import (
//	"bytes"
//	"os"
//	"golang.org/x/crypto/ssh"
////	"cloud/benchflow/environment"
//)

//TODO: refactor

////TODO: Update to the new procedure that uses docker-machine
////TODO: extend also to other machine types, automatically or by defining another property in the configuration file and using ubuntu as default
//func InstallDocker() {
//	
//	if(len(os.Args) <= 2){
//		fmt.Println("Expected server Name as parameter")
//		return
//	}
//	
//	var serverName = os.Args[2]
//	
//	//TODO: validate server name
//	fmt.Println("Installing and Configuring Docker on: " + serverName)
//	fmt.Println()
//	
////	var serversMap = Env.Servers
//	var serversMap = ""
//
//	sshConfig := &ssh.ClientConfig{
//		User: serversMap[serverName].SudoUserName,
//		Auth: []ssh.AuthMethod{
//			ssh.Password(serversMap[serverName].SudoUserPassword),
//		},
//	}
//	
//	fmt.Println("User: " + serversMap[serverName].SudoUserName)
//
//	client, err := ssh.Dial("tcp", serversMap[serverName].ExternalIP + ":22", sshConfig)
//
//	if err != nil {
//		fmt.Println("Failed to dial: " + err.Error())
//		//TODO: show the following only when the error state that the server is not ssh enabled
//		fmt.Println("The server must be ssh enabled")
//		return
//	}
//
//	session, err := client.NewSession()
//	if err != nil {
//		fmt.Println("Failed to create session: " + err.Error())
//		return
//	}
//
//	defer session.Close()
//
//	var b bytes.Buffer
//	session.Stdout = &b
//	if err := session.Run("ls /"); err != nil {
//		fmt.Println("Failed to run: " + err.Error())
//		return
//	}
//
//	fmt.Println(b.String())
//}
//
//func openSshConnection() {
//	
//}
//
//func installDocker() {
//	
//}
//
//func closeSshConnection() {
//	
//}
