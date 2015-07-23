package structs

type Environment struct {
     Servers map[string]Server
//     Docker structs.Docker
     BenchFlow BenchFlow
     Settings Settings
     Credentials Credentials
}

type Server struct {
	Name             string
	Alias			 string
	ExternalIP       string
	LocalIP          string
	Hostname		 string
	DockerPort		 int
	SshPort			 int
	SudoUserName     string
	SudoUserPassword string
	Purpose		 	 string
}

type Settings struct {
	//TODO: change all structs and configurations for using the "yaml" mapping explicit declaration
	LocalCertificateFolder string `yaml:"localCertificateFolder"`
    RemoteCertificatesFolder string `yaml:"remoteCertificatesFolder"`
    RemoteDockerComposeProjects string `yaml:"remoteDockerComposeProjects"`
}

type BenchFlow struct {
	DockerImages	DockerImages
}

type DockerImages struct {
	DockerSwarm		ImageDetails
	DockerCompose	ImageDetails
	Cassandra		ImageDetails
	Rancher			ImageDetails
	RancherCompose	ImageDetails
	Consul			ImageDetails
}

type ImageDetails struct {
	Image				string
	RunPorts			[]string
	DataRootDirectory	string
}

type Credentials struct {
	Rancher		CredentialsDetails
	Cassandra	CredentialsDetails
}

type CredentialsDetails struct {
	Username	string
	Password	string
}