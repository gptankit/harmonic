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

	retryIndex, retryLimit, svc := 0, len(servicelist)-1, ""

	for retryIndex <= retryLimit {

		// call SelectService
		svc, _ = harmonic.SelectService(cs, retryIndex, svc)

		// send request to resource located at svc
		response := makeRequestToSvc()

		// if success, then reset error for service and break
		if response == "SUCCESS" {
			cs.ResetError(svc)
			break
		} else { // if failed, then increment error for service and retryIndex
			cs.IncrementError(svc)
			retryIndex++
		}

		// optional (test for current errorcount for svc)
		// ec, _ := cs.GetError(svc)
	}
}

// replace below function with your service call
func makeRequestToSvc() string {

	// return "SUCCESS"
	return "FAIlED"
}
