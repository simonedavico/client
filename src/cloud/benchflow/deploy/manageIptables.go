package deploy

import (
	crypto "golang.org/x/crypto/ssh"
	"cloud/benchflow/environment"
	"cloud/benchflow/utils/logging"
	"cloud/benchflow/structs"
	"cloud/benchflow/ssh"
)

//TODO: save in benchflow folder

func SaveIptables() {
	
	//Servers
    servers := environment.Env.Servers

	var server structs.Server
	
	//Save iptables for each of the servers
	for _, server = range servers {
		
		var home = ssh.GetHomeDirectory(server)
	
		var session *crypto.Session = ssh.OpenSshConnection(server)
		
		if err := session.Run("echo " + server.SudoUserPassword + " | sudo -S iptables-save > " + home + "/iptable.rules"); err != nil {
			logging.Log.Fatal("Failed to run: " + err.Error())
		}
		
		ssh.CloseSshConnection(session)
		    
	}
	
	
}


func RestoreIptables() {
	
	//Servers
    servers := environment.Env.Servers

	var server structs.Server
	
	//Save iptables for each of the servers
	for _, server = range servers {
		
		var home = ssh.GetHomeDirectory(server)
	
		var session *crypto.Session = ssh.OpenSshConnection(server)
		
		//The command run in an elevated privilege shell because it has < that conflicts with -S
		if err := session.Run("echo " + server.SudoUserPassword + " | sudo -S bash -c \"iptables-restore < " + home + "/iptable.rules && rm -f " + home + "/iptable.rules\""); err != nil {
			logging.Log.Fatal("Failed to run: " + err.Error())
		}
		
		ssh.CloseSshConnection(session)
		    
	}
}