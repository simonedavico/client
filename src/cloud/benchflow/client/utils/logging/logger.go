package logging

import (
	"os"
	log "github.com/Sirupsen/logrus"
)

var Log *log.Logger

//TODO: Configure by using the environment that reads the configuration
func InitialiseLogger() {
	
    if Log == nil {
       	Log = log.New()
       	Log.Out = os.Stdout
	    Log.Level = log.DebugLevel
	//    Log.Formatter = new(log.JSONFormatter)
	//    Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
	//        log.InfoLevel : "/var/log/info.log",
	//        log.ErrorLevel : "/var/log/error.log",
	//    }))
    }
    
}