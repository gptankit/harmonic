package main

import (
	"github.com/gptankit/harmonic"
)

func main() {

	servicelist := []string{"s0", "s1", "s2"}
	cs, err := harmonic.InitClusterState(servicelist)
	if err != nil {
		return
	}

	// get a request to process (e.g. from a queue, or datastore, or client request if exposed as an api)
	// req := getRequest()

	// initialize parameters
	retryindex, svc := 0, ""
	// retrylimit is recommended to be equal to size of service list
	retrylimit := len(cs.GetServices())

	for retryindex < retrylimit {

		// call SelectService
		svc, _ = harmonic.SelectService(cs, retryindex, svc)

		// send request to resource located at svc (e.g. execute query, or call external api)
		response := makeRequestToSvc()

		// if success, then reset error for service and break
		if response == "SUCCESS" {
			cs.ResetError(svc)
			break
		} else { // if failed, then increment error for service and retryIndex
			cs.IncrementError(svc)
			retryindex++
		}
	}
}

// replace below function with your service call
func makeRequestToSvc() string {

	// return "SUCCESS"
	return "FAIlED"
}
