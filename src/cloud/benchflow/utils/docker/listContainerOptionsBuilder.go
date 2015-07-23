package docker

import (
	"github.com/lann/builder"
	"github.com/fsouza/go-dockerclient"
)

type listContainerOptionsBuilder builder.Builder

func (b listContainerOptionsBuilder) All(all bool) listContainerOptionsBuilder {
    return builder.Set(b, "All", all).(listContainerOptionsBuilder)
}

func (b listContainerOptionsBuilder) Filters(filters map[string][]string) listContainerOptionsBuilder {
    return builder.Set(b, "Filters", filters).(listContainerOptionsBuilder)
}


func (b listContainerOptionsBuilder) Build() docker.ListContainersOptions {
    return builder.GetStruct(b).(docker.ListContainersOptions)
}

var ListContainerOptionsBuilder = builder.Register(listContainerOptionsBuilder{}, docker.ListContainersOptions{}).(listContainerOptionsBuilder)