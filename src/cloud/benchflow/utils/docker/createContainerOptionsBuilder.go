package docker

import (
	"github.com/lann/builder"
	"github.com/fsouza/go-dockerclient"
)

type createContainerOptionsBuilder builder.Builder

func (b createContainerOptionsBuilder) Name(name string) createContainerOptionsBuilder {
    return builder.Set(b, "Name", name).(createContainerOptionsBuilder)
}

func (b createContainerOptionsBuilder) Config(config *docker.Config) createContainerOptionsBuilder {
    return builder.Set(b, "Config", config).(createContainerOptionsBuilder)
}

func (b createContainerOptionsBuilder) HostConfig(hostConfig *docker.HostConfig) createContainerOptionsBuilder {
    return builder.Set(b, "HostConfig", hostConfig).(createContainerOptionsBuilder)
}

func (b createContainerOptionsBuilder) Build() docker.CreateContainerOptions {
    return builder.GetStruct(b).(docker.CreateContainerOptions)
}

var CreateContainerOptionsBuilder = builder.Register(createContainerOptionsBuilder{}, docker.CreateContainerOptions{}).(createContainerOptionsBuilder)