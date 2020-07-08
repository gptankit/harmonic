package harmonic

import (
	"errors"
	"sync"
)

type ClusterState struct {
	servicelist []string
	numservices int
	errormap    map[string]uint64
	remutex     sync.Mutex
}

// InitClusterState initializes cluster state with user provided
// servicelist, and zeroing error for each service.
func InitClusterState(servicelist []string) (*ClusterState, error) {

	if servicelist == nil || len(servicelist) == 0 {
		return nil, errors.New("harmonic: invalid service list")
	}

	errormap := make(map[string]uint64)
	for _, svc := range servicelist {
		errormap[svc] = 0
	}

	return &ClusterState{
		servicelist: servicelist,
		numservices: len(servicelist),
		errormap:    errormap,
	}, nil
}

// GetError returns current errorcount for a service.
func (cs *ClusterState) GetError(service string) (uint64, error) {

	cs.remutex.Lock()
	defer cs.remutex.Unlock()

	if _, ok := cs.errormap[service]; ok {
		return cs.errormap[service], nil
	}
	return 0, errors.New("harmonic: service " + service + " not found")
}

// IncrementError increments errorcount by 1 for a service.
func (cs *ClusterState) IncrementError(service string) error {

	cs.remutex.Lock()
	defer cs.remutex.Unlock()

	if _, ok := cs.errormap[service]; ok {
		cs.errormap[service]++
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// UpdateError updates user provided errorcount for a service.
func (cs *ClusterState) UpdateError(service string, errorcount uint64) error {

	cs.remutex.Lock()
	defer cs.remutex.Unlock()

	if _, ok := cs.errormap[service]; ok {
		cs.errormap[service] = errorcount
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// ResetError resets errorcount for a service.
func (cs *ClusterState) ResetError(service string) error {

	cs.remutex.Lock()
	defer cs.remutex.Unlock()

	if _, ok := cs.errormap[service]; ok {
		cs.errormap[service] = 0
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// ResetAllErrors resets errorcount for all services.
func (cs *ClusterState) ResetAllErrors() error {

	cs.remutex.Lock()
	defer cs.remutex.Unlock()

	for k, _ := range cs.errormap {
		cs.errormap[k] = 0
	}
	return nil
}
