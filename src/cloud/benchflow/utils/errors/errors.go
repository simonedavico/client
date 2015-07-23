package errors

import (
	"cloud/benchflow/utils/logging"
	log "github.com/Sirupsen/logrus"
)

func CheckFatal(e error) {
	
    if e != nil {
       	logging.Log.WithFields(log.Fields{
			  "error": e,
			  }).Fatal()
    }
}