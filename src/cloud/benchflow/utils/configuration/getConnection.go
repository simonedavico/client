package configuration

import (
	"strings"
	"strconv"
	"cloud/benchflow/structs"
)

func GetServerEndPoint(server structs.Server) string {
	
	var serverEndPoint string
	
    serverEndPoint = GetServerAddress(server) + ":" + strconv.Itoa(server.DockerPort)
    
    return serverEndPoint
    
}

func GetServerAddress(server structs.Server) string {
	
	var serverAddress string
	
	//NAME Preference order: alias, LocalIP, ExternalIP
    if len(strings.TrimSpace(server.Alias)) > 0 {
    	serverAddress = server.Alias
    } else {
    	serverAddress = GetServerIP(server)
    }
    
    return serverAddress
    
}

func GetServerExternalIP(server structs.Server) string {
	
	var serverEndPoint string
	
	serverEndPoint = server.ExternalIP

    return serverEndPoint
    
}

func GetServerIP(server structs.Server) string {
	
	var serverEndPoint string
	
	//IP Preference order: LocalIP, ExternalIP
    if len(strings.TrimSpace(server.LocalIP)) > 0 {
    	serverEndPoint = server.LocalIP
    } else {
    	serverEndPoint = GetServerExternalIP(server)
    }

    return serverEndPoint
    
}