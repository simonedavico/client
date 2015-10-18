package environment

import (	
	"os"
	"cloud/benchflow/client/structs"
	"cloud/benchflow/client/utils/filesystem"
	"cloud/benchflow/client/utils/configuration"
)

//Share common functionalities by using the Environment Object Pattern (http://www.jerf.org/iri/post/2929) and a single instance
//It is initialized once in the main, then accessed where needed

//TODO, maybe remove and change to environment variables that can be passed to env and flags

//TODO: publish on consul only the needed ones (e.g., not server access and sensible data)

//TODO: refactor to avoid globals, and maybe direclty use configuration manager
var Env *structs.Environment = nil

func InitialiseEnvironment() {

	//This is useful when used as running service, not from command line
	if Env == nil {
	 	
	 	Env = &structs.Environment{}
	 	
		Env.Servers = configuration.GetConfServer()
		Env.Settings = configuration.GetSettings()
		Env.BenchFlow = configuration.GetConfBenchFlow()
		Env.Credentials = configuration.GetConfCredentials()
		
		//Set the local certificate folder:
		//certs
		//	general
		//		ca.pem
		//	serverName1
		//		cert.pem, key.pem
		Env.Settings.LocalCertificateFolder = filesystem.GetExecutableFolder() + string(os.PathSeparator) + Env.Settings.LocalCertificateFolder + string(os.PathSeparator) //Local to the directory where we run the benchflow-client
		
    }

}

func UpdateCredentials() {
	Env.Credentials = configuration.GetConfCredentials()
}