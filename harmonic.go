package harmonic

import (
	"math"
	"errors"
)

// SelectService implements the routing logic to the cluster of downstream services. On 
// first try, an error log lookup is done to determine the service-wise error count and effective
// error is calculated. If no error found for any service, random service selection (equal probability)
// is done, else weighted random service selection is done, where weights are inversely proportional 
// to error count on the particular service. If the request to the selected service fails, round robin
// selection is done to deterministically select the next service.
func SelectService(cs *ClusterState, retryindex int, prevservice string) (string, error) {

	//invalid num of endpoints
	if cs.numservices == 0{
		return "", errors.New("harmonic: servicelist is empty")
	}

	// single endpoint 
	if cs.numservices == 1 {
		return getIndexedService(0, cs)
	}

	if retryindex == 0 { // first try
		cs.remutex.Lock()
		defer cs.remutex.Unlock()
		maxErr := uint64(0)
		for _, svc := range cs.servicelist {
			errCnt := cs.errormap[svc]
			effectiveErr := uint64(math.Floor(math.Pow(float64(1 + errCnt), 1.5)))
			if effectiveErr >= maxErr {
				maxErr = effectiveErr
			}
		}
		if maxErr == 1 {
			return getIndexedService(randomize(0, cs.numservices), cs)
		} else {
			weights := make([]float64, cs.numservices)
			prefixes := make([]float64, cs.numservices)
			for i, svc := range cs.servicelist {
				errCnt := cs.errormap[svc]
				weights[i] = math.Ceil(float64(maxErr) / float64(errCnt + 1))
			}
			for i, _ := range weights {
				if i == 0 {
					prefixes[i] = weights[i]
				} else {
					prefixes[i] = weights[i] + prefixes[i-1]
				}
			}
			prLen := len(prefixes) - 1
			randx := randomize64(1, int64(prefixes[prLen]) + 1)
			ceil := findCeilIn(randx, prefixes, 0, prLen)
			if ceil >= 0 {
				return getIndexedService(ceil, cs)
			}
		}
		return getIndexedService(randomize(0, cs.numservices), cs)
	} else {
		prevserviceindex := -1
		for psi, svc := range cs.servicelist{
			if svc == prevservice{
				prevserviceindex = psi
			}
		}
		return getIndexedService(roundrobin(cs.numservices, prevserviceindex), cs)
	}
}

// findCeilIn does a binary search to find position of selected random
// number and returns corresponding ceil index in prefixes array.
func findCeilIn(randx int64, prefixes []float64, start int, end int) int {

	var mid int
	for {
		if start >= end {
			break
		}
		mid = start + ((end - start) >> 1)
		if randx > int64(prefixes[mid]) {
			start = mid + 1
		} else {
			end = mid
		}
	}

	if randx <= int64(prefixes[start]) {
		return start
	}
	return -1
}

// getIndexedService returns service name at an index. Error is returned
// if index is found to be invalid.
func getIndexedService(index int, cs *ClusterState) (string, error){

	if index < 0 || index >= cs.numservices{
		return "", errors.New("harmonic: service index out of bounds")
	}

	return cs.servicelist[index], nil
}
