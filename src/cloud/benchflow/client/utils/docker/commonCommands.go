package docker

import (
	"bytes"
	"cloud/benchflow/client/utils/logging"
	"github.com/fsouza/go-dockerclient"
	"strings"
)

func DeployImage(client *docker.Client, name string, tag string) {

	logging.Log.Debug("Test Import Image")
	
    var buf bytes.Buffer
    
    image := docker.PullImageOptions{
    			Repository: name,
    			Tag: tag,
    			OutputStream: &buf,
    		 }
    
    res := client.PullImage(image, docker.AuthConfiguration{})
    
    logging.Log.Debug("Output: " + buf.String()) //Output of the command
    
    logging.Log.Debug("Res: ")
    logging.Log.Debug(res)
}


func GetImageAndTag(image string) (string, string) {
	
	Image := strings.Split(image, ":")
    Name, Tag := Image[0], Image[1]
    
    return Name,Tag

}