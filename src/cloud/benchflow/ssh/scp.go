package ssh

import (
	"bytes"
	"cloud/benchflow/structs"
	"cloud/benchflow/utils/filesystem"
	"cloud/benchflow/utils/logging"
	"fmt"
//	"os"
//	"io"
//	"io/ioutil"
	crypto "golang.org/x/crypto/ssh"
)

func ScpFolderInUserHome(server structs.Server, folderPathFrom string, folderNameTo string) {

	var home = GetHomeDirectory(server)
	var folderPathTo = home + "/" + folderNameTo + "/"
	
	CreateRemoteDirectory(server,home,folderPathTo)
	
	ScpFilesInFolder(server, folderPathFrom, folderPathTo)

}

func ScpFileInUserHome(server structs.Server, content string, serverDirectoryName string, fileNameTo string) {

	var home = GetHomeDirectory(server)
	var folderPathTo = home + "/" + serverDirectoryName + "/"

	CreateRemoteDirectory(server,home,folderPathTo)

	ScpFileContent(server, content, folderPathTo, fileNameTo)

}

func ScpFilesInFolder(server structs.Server, folderPathFrom string, folderPathTo string) {

	var session *crypto.Session = OpenSshConnection(server)

	var b bytes.Buffer
	session.Stdout = &b

	//Create the folderPathTo directory if missing
	//Copy files to the folderPathTo directory
	go func() {

		var filesList = filesystem.GetFilesInDirectory(folderPathFrom)

		w, _ := session.StdinPipe()
		defer w.Close()

		//		fmt.Fprintln(w, "D0755", 0, "certs") // mkdir

		//Copy files
		for _, file := range filesList {
			//TODO: refactor in a function to avoid duplicated code with the next function ScpFile (if possible)
			if file.Mode().IsRegular() {

				var content = filesystem.GetFileContent(folderPathFrom + file.Name())
		
				fmt.Fprintln(w, "C0644", len(content), file.Name())
				fmt.Fprint(w, content)
				fmt.Fprint(w, "\x00") // transfer end with \x00
		
			}
		}

	}()

	if err := session.Run("/usr/bin/scp -tr " + folderPathTo); err != nil {
		logging.Log.Fatal("Failed to run: " + err.Error())
	}

	//	fmt.Println(b.String())

	CloseSshConnection(session)
}

func ScpFileContent(server structs.Server, content string, serverDirectory string, fileNameTo string) {

	var session *crypto.Session = OpenSshConnection(server)

	var b bytes.Buffer
	session.Stdout = &b

	//Create the filePathTo directory if missing
	//Copy file to the filePathTo file
	go func() {

				w, _ := session.StdinPipe()
				defer w.Close()
			
		//	if file.Mode().IsRegular() { #TODO: find out how to solve the error go shows on this line
		
//				logging.Log.Debug(content)
		
				fmt.Fprintln(w, "C0644", len(content), fileNameTo)
//				logging.Log.Debug(content)
				fmt.Fprint(w, content)
//				logging.Log.Debug(content)
				fmt.Fprint(w, "\x00") // transfer end with \x00
//				logging.Log.Debug(content)
		
		//	}
		
			logging.Log.Debug("Copy in: " + "/usr/bin/scp -tr " + serverDirectory)
	
	}()


	if err := session.Run("/usr/bin/scp -tr " + serverDirectory + "/"); err != nil {
		logging.Log.Fatal("Failed to run: " + err.Error())
	}

	//	fmt.Println(b.String())

	CloseSshConnection(session)
}

