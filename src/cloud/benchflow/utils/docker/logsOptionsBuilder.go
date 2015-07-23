package docker

import (
	"bytes"
	"github.com/lann/builder"
	"github.com/fsouza/go-dockerclient"
)

type logOptionsBuilder builder.Builder

func (b logOptionsBuilder) Container(id string) logOptionsBuilder {
    return builder.Set(b, "Container", id).(logOptionsBuilder)
}

func (b logOptionsBuilder) ErrorStream(buffer *bytes.Buffer) logOptionsBuilder {
    return builder.Set(b, "ErrorStream", buffer).(logOptionsBuilder)
}

func (b logOptionsBuilder) OutputStream(buffer *bytes.Buffer) logOptionsBuilder {
    return builder.Set(b, "OutputStream", buffer).(logOptionsBuilder)
}

func (b logOptionsBuilder) Follow(follow bool) logOptionsBuilder {
    return builder.Set(b, "Follow", follow).(logOptionsBuilder)
}

func (b logOptionsBuilder) Stdout(attachStdout bool) logOptionsBuilder {
    return builder.Set(b, "Stdout", attachStdout).(logOptionsBuilder)
}

func (b logOptionsBuilder) Stderr(attachStderr bool) logOptionsBuilder {
    return builder.Set(b, "Stderr", attachStderr).(logOptionsBuilder)
}

func (b logOptionsBuilder) Timestamps(timestamps bool) logOptionsBuilder {
    return builder.Set(b, "Timestamps", timestamps).(logOptionsBuilder)
}


func (b logOptionsBuilder) Build() docker.LogsOptions {
    return builder.GetStruct(b).(docker.LogsOptions)
}

var LogOptionsBuilder = builder.Register(logOptionsBuilder{}, docker.LogsOptions{}).(logOptionsBuilder)