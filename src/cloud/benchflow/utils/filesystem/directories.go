package filesystem

import (
	"os"
	"github.com/kardianos/osext"
	"cloud/benchflow/utils/errors"
)


func GetFilesInDirectory(directory string) []os.FileInfo {
	 
	 //Open the folderPathFrom directory
	 d, err := os.Open(directory)
     errors.CheckFatal(err)
     defer d.Close()

     files, err := d.Readdir(-1)
     errors.CheckFatal(err)
	
	 return files
}

func GetExecutableFolder() string {
	
	folderPath, err := osext.ExecutableFolder()
    errors.CheckFatal(err)
    
    return folderPath
	
}