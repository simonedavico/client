package docker

import (
	"cloud/benchflow/client/structs"
	"github.com/lann/builder"
	"github.com/fsouza/go-dockerclient"
	"cloud/benchflow/client/utils/configuration"
)

type dockerHostConfigBuilder builder.Builder

func (b dockerHostConfigBuilder) AddBind(bind string) dockerHostConfigBuilder {
    return builder.Append(b, "Binds", bind).(dockerHostConfigBuilder)
}

func (b dockerHostConfigBuilder) ExtraHosts(servers []structs.Server) dockerHostConfigBuilder {
    return builder.Set(b, "ExtraHosts", buildExtraHosts(servers)).(dockerHostConfigBuilder)
}

func (b dockerHostConfigBuilder) PortBindings(portbindings map[docker.Port][]docker.PortBinding) dockerHostConfigBuilder {
    return builder.Set(b, "PortBindings", portbindings).(dockerHostConfigBuilder)
}

func (b dockerHostConfigBuilder) RestartPolicy(restartPolicy docker.RestartPolicy) dockerHostConfigBuilder {
    return builder.Set(b, "RestartPolicy", restartPolicy).(dockerHostConfigBuilder)
}

func (b dockerHostConfigBuilder) NetworkMode(networkMode string) dockerHostConfigBuilder {
    return builder.Set(b, "NetworkMode", networkMode).(dockerHostConfigBuilder)
}

func (b dockerHostConfigBuilder) Build() docker.HostConfig {
    return builder.GetStruct(b).(docker.HostConfig)
}


func buildExtraHosts(servers []structs.Server) []string {
	
	var extraHosts []string
	
	var server structs.Server
	
	for _, server = range servers {
		
		//We need to define additional host only if we use aliases
		if server.Alias != "" {
			//Get the utilized IP
			var serverIP = configuration.GetServerIP(server)
			extraHosts = append(extraHosts,server.Alias + ":" + serverIP)
		}
		
	}
	
	return extraHosts		
}

var DockerHostConfigBuilder = builder.Register(dockerHostConfigBuilder{}, docker.HostConfig{}).(dockerHostConfigBuilder)