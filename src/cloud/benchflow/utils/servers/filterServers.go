package servers

import (
	"cloud/benchflow/structs"
	"strings"
)

func GetMasterServer(servers map[string]structs.Server) structs.Server { //structs.Server

	master := getSingleServerByPurposeMap("master", servers)

	//	fmt.Println(master)

	return master
}

func GetAnalyserNotMasterServers(servers map[string]structs.Server) []structs.Server { //structs.Server

	list := getListOfServersByPurposeMap("master", false, servers)

	updatedList := getListOfServersByPurposeList("analyser", true, list)

	//	fmt.Println(master)

	return updatedList
}

func GetNotMasterServers(servers map[string]structs.Server) []structs.Server { //structs.Server

	list := getListOfServersByPurposeMap("master", false, servers)

	return list
}

func GetNotSutNotDriversNotMasterServers(servers map[string]structs.Server) []structs.Server { //structs.Server

	list := getListOfServersByPurposeMap("sut,drivers,master", false, servers)

	return list
}

func getSingleServerByPurposeMap(purpose string, servers map[string]structs.Server) structs.Server {

	var matched structs.Server
	var value structs.Server

	for _, value = range servers {

		if strings.Contains(value.Purpose, purpose) {
			matched = value
			break
		}
	}

	return matched

}

//Get a list of server with the specifies purposes (keep = true) or without the specified purpose (keep = false)
//It is sufficient that the purpose is present in the list of purposes assigned to the server
func getListOfServersByPurposeMap(purposes string, keep bool, servers map[string]structs.Server) []structs.Server {

	var matched []structs.Server
	var value structs.Server

	var pusposesList = strings.Split(purposes, ",")

	for _, value = range servers {

		found := filterByPurpose(pusposesList, keep, value)

		if keep && found {
			matched = append(matched, value)
		} else if !keep && !found {
			matched = append(matched, value)
		}
	}

	return matched

}

//The same as getListOfServersByPurpose, but for arrays
func getListOfServersByPurposeList(purposes string, keep bool, servers []structs.Server) []structs.Server {
	
	var matched []structs.Server
	var value structs.Server

	var pusposesList = strings.Split(purposes, ",")

	for	_, value = range servers {

		found := filterByPurpose(pusposesList, keep, value)

		if keep && found {
			matched = append(matched, value)
		} else if !keep && !found {
			matched = append(matched, value)
		}
	}

	return matched
	
}

func filterByPurpose(pusposesList []string, keep bool, servers structs.Server) bool {

	var found = false

	for _, purpose := range pusposesList {

		if keep && strings.Contains(servers.Purpose, purpose) {
			found = true
			break
		} else if !keep && strings.Contains(servers.Purpose, purpose) {
			found = true
			break
		}

	}

	return found

}
