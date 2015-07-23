package docker

import (
	"github.com/lann/builder"
	"github.com/fsouza/go-dockerclient"
)

type dockerConfigBuilder builder.Builder

func (b dockerConfigBuilder) Hostname(hostname string) dockerConfigBuilder {
    return builder.Set(b, "Hostname", hostname).(dockerConfigBuilder)
}

func (b dockerConfigBuilder) Image(name string, tag string) dockerConfigBuilder {
    return builder.Set(b, "Image", name + ":" + tag).(dockerConfigBuilder)
}

func (b dockerConfigBuilder) AddCmd(cmd string) dockerConfigBuilder {
    return builder.Append(b, "Cmd", cmd).(dockerConfigBuilder)
}

func (b dockerConfigBuilder) AddEnv(env string) dockerConfigBuilder {
    return builder.Append(b, "Env", env).(dockerConfigBuilder)
}

func (b dockerConfigBuilder) Volumes(volumes map[string]struct{}) dockerConfigBuilder {
    return builder.Set(b, "Volumes", volumes).(dockerConfigBuilder)
}

func (b dockerConfigBuilder) ExposedPorts(exposedPorts map[docker.Port]struct{}) dockerConfigBuilder {
    return builder.Set(b, "ExposedPorts", exposedPorts).(dockerConfigBuilder)
}


func (b dockerConfigBuilder) Build() docker.Config {
    return builder.GetStruct(b).(docker.Config)
}

var DockerConfigBuilder = builder.Register(dockerConfigBuilder{}, docker.Config{}).(dockerConfigBuilder)