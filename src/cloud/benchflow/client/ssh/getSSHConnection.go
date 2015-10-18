package ssh

import (
	"strconv"
	"golang.org/x/crypto/ssh"
	"cloud/benchflow/client/structs"
	"cloud/benchflow/client/utils/logging"
)

func OpenSshConnection(server structs.Server) *ssh.Session {
	
//	fmt.Println("Server: ")
//	fmt.Println(server)
	
	sshConfig := &ssh.ClientConfig{
		User: server.SudoUserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.SudoUserPassword),
		},
	}
	
//	fmt.Println("User: " + server.SudoUserName)

	client, err := ssh.Dial("tcp", server.ExternalIP + ":" + strconv.Itoa(server.SshPort), sshConfig)

	if err != nil {
		logging.Log.Fatal("Failed to dial: " + err.Error())
		//TODO: show the following only when the error state that the server is not ssh enabled
		logging.Log.Fatal("The server must be ssh enabled")
	}

	session, err := client.NewSession()
	if err != nil {
		logging.Log.Fatal("Failed to create session: " + err.Error())
	}

	return session
}

func CloseSshConnection(session *ssh.Session)  {
	session.Close()
}
