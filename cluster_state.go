package harmonic

import (
	"sync"
	"errors"
)

type ClusterState struct {
	servicelist             []string
	numservices		int
	errormap                map[string]uint64
	remutex                 sync.Mutex
}

// InitClusterState initializes cluster state with user provided 
// servicelist, and zeroing error for each service.
func InitClusterState(servicelist []string) (*ClusterState, error) {

	if servicelist == nil || len(servicelist) == 0{
		return nil, errors.New("harmonic: invalid service list")
	}

	errormap := make(map[string]uint64)
	for _, svc := range servicelist{
		errormap[svc] = 0
	}

	return &ClusterState{
		servicelist: servicelist,
		numservices: len(servicelist),
		errormap: errormap,
	}, nil
}

// IncrementError increments errorcount by 1 for a service.
func (cp *ClusterState) IncrementError(service string) error{

	cp.remutex.Lock()
	defer cp.remutex.Unlock()

	if _, ok := cp.errormap[service];ok{
		cp.errormap[service]++
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// UpdateError updates user provided errorcount for a service.
func (cp *ClusterState) UpdateError(service string, errorcount uint64) error{

	cp.remutex.Lock()
	defer cp.remutex.Unlock()

	if _, ok := cp.errormap[service];ok{
		cp.errormap[service] = errorcount
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// ResetError resets errorcount for a service.
func (cp *ClusterState) ResetError(service string) error{

	cp.remutex.Lock()
	defer cp.remutex.Unlock()

	if _, ok := cp.errormap[service];ok{
		cp.errormap[service] = 0
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// ResetAllErrors resets errorcount for all services.
func (cp *ClusterState) ResetAllErrors(){

	cp.remutex.Lock()
	defer cp.remutex.Unlock()

	for k, _ := range cp.errormap{
		cp.errormap[k] = 0
	}
}
