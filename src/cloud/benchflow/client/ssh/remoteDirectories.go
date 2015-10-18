package ssh

import (
	"bytes"
	"strings"
	"cloud/benchflow/client/structs"
	"cloud/benchflow/client/utils/logging"
	crypto "golang.org/x/crypto/ssh"
)

//TODO: move in conf or in tmp file the first time we retrieve it
func GetHomeDirectory(server structs.Server) string {
	var session *crypto.Session = OpenSshConnection(server)
	var home string
	
	var b bytes.Buffer
	session.Stdout = &b
	
	if err := session.Run("echo $HOME"); err != nil {
		logging.Log.Fatal("Failed to run \"echo $HOME\": " + err.Error())
	}
	
	home = b.String()
	home = strings.TrimRight(home, "\n")
	
	logging.Log.Debug("Home: " + home)
	
	CloseSshConnection(session)
	
	return home
}

func DeleteRemoteDirectoryInsideHome(server structs.Server, directory string) {
	
	var directoryToRemove string

	directoryToRemove = GetHomeDirectory(server) + "/" + directory + "/"
	
	logging.Log.Debug("directoryToRemove: " + directoryToRemove)
	
	DeleteRemoteDirectory(server, directoryToRemove)
	
}

func CreateRemoteDirectory(server structs.Server, insideFolder string, folderPathTo string) {
	
	//	fmt.Println("folderPathTo: " + folderPathTo)

	var session *crypto.Session = OpenSshConnection(server)

	//	fmt.Println("cd " + home + " && mkdir -p " + folderNameTo)

	if err := session.Run("cd " + insideFolder + " && mkdir -p " + folderPathTo); err != nil {
		logging.Log.Fatal("Failed to run \"mkdir\": " + err.Error())
	}

	CloseSshConnection(session)

	
}


func DeleteRemoteDirectory(server structs.Server, directory string) {

	var session *crypto.Session = OpenSshConnection(server)
	
	session = OpenSshConnection(server)
	
	if err := session.Run("rm -fR " + directory); err != nil {
		logging.Log.Fatal("Failed to run: " + err.Error())
	}
	
	CloseSshConnection(session)
	
}