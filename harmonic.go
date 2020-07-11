package harmonic

import (
	"errors"
	"math"
)

// SelectService implements the routing logic to the cluster of downstream services. On
// first try, an error log lookup is done to determine the service-wise error count and effective
// error is calculated. If no error found for any service, random service selection (equal probability)
// is done, else weighted random service selection is done, where weights are inversely proportional
// to error count on the particular service. If the request to the selected service fails, round robin
// selection is done to deterministically select the next service.
func SelectService(cs *ClusterState, retryIndex int, prevService string) (string, error) {

	//invalid num of endpoints
	if cs.numServices == 0 {
		return "", errors.New("harmonic: service list is empty")
	}

	// single endpoint
	if cs.numServices == 1 {
		return getIndexedService(cs, 0)
	}

	if retryIndex == 0 { // first try
		cs.reMutex.Lock()
		defer cs.reMutex.Unlock()

		maxErr := uint64(0)

		for _, svc := range cs.serviceList {
			errCnt := cs.errorMap[svc]
			effectiveErr := uint64(math.Floor(math.Pow(float64(1+errCnt), 1.5)))
			if effectiveErr >= maxErr {
				maxErr = effectiveErr
			}
		}

		if maxErr == 1 {
			return getIndexedService(cs, randomize(0, cs.numServices))
		} else {
			weights := make([]float64, cs.numServices)
			prefixes := make([]float64, cs.numServices)

			for i, svc := range cs.serviceList {
				errCnt := cs.errorMap[svc]
				weights[i] = math.Ceil(float64(maxErr) / float64(errCnt+1))
			}

			for i, _ := range weights {
				if i == 0 {
					prefixes[i] = weights[i]
				} else {
					prefixes[i] = weights[i] + prefixes[i-1]
				}
			}

			prLen := cs.numServices - 1
			randx := randomize64(1, int64(prefixes[prLen])+1)
			ceil := findCeilIn(randx, prefixes, 0, prLen)

			if ceil >= 0 {
				return getIndexedService(cs, ceil)
			}
		}

		return getIndexedService(cs, randomize(0, cs.numServices))
	} else { // retries
		prevServiceIndex := -1
		for psi, svc := range cs.serviceList {
			if svc == prevService {
				prevServiceIndex = psi
			}
		}

		return getIndexedService(cs, roundrobin(cs.numServices, prevServiceIndex))
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
func getIndexedService(cs *ClusterState, serviceIndex int) (string, error) {

	if serviceIndex < 0 || serviceIndex >= cs.numServices {
		return "", errors.New("harmonic: service index out of bounds")
	}

	return cs.serviceList[serviceIndex], nil
}
