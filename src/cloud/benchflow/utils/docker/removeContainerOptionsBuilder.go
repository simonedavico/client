package docker

import (
	"github.com/lann/builder"
	"github.com/fsouza/go-dockerclient"
)

type removeContainerOptionsBuilder builder.Builder

func (b removeContainerOptionsBuilder) ID(id string) removeContainerOptionsBuilder {
    return builder.Set(b, "ID", id).(removeContainerOptionsBuilder)
}

func (b removeContainerOptionsBuilder) RemoveVolumes(remove bool) removeContainerOptionsBuilder {
    return builder.Set(b, "RemoveVolumes", remove).(removeContainerOptionsBuilder)
}

func (b removeContainerOptionsBuilder) Force(force bool) removeContainerOptionsBuilder {
    return builder.Set(b, "Force", force).(removeContainerOptionsBuilder)
}


func (b removeContainerOptionsBuilder) Build() docker.RemoveContainerOptions {
    return builder.GetStruct(b).(docker.RemoveContainerOptions)
}

var RemoveContainerOptionsBuilder = builder.Register(removeContainerOptionsBuilder{}, docker.RemoveContainerOptions{}).(removeContainerOptionsBuilder)