package filesystem

import (
//	"fmt"
	"io/ioutil"
	"cloud/benchflow/client/utils/errors"
)

func GetFileContent(filePath string) string {
	
	dat, err := ioutil.ReadFile(filePath)
	//Reading files requires checking most calls for errors. This helper will streamline our error checks below.
    errors.CheckFatal(err)
//    fmt.Print(string(dat))
    return string(dat)
    
}