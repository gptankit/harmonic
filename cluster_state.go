package harmonic

import (
	"errors"
	"sync"
)

type ClusterState struct {
	serviceList []string
	numServices int
	errorMap    map[string]uint64
	reMutex     sync.Mutex
}

// InitClusterState initializes cluster state with user provided
// serviceList, and zeroing error for each service.
func InitClusterState(serviceList []string) (*ClusterState, error) {

	if serviceList == nil || len(serviceList) == 0 {
		return nil, errors.New("harmonic: invalid service list")
	}

	errorMap := make(map[string]uint64)
	for _, svc := range serviceList {
		errorMap[svc] = 0
	}

	return &ClusterState{
		serviceList: serviceList,
		numServices: len(serviceList),
		errorMap:    errorMap,
	}, nil
}

// GetServices returns current list of services.
func (cs *ClusterState) GetServices() []string {

	return cs.serviceList
}

// GetError returns current errorcount for a service.
func (cs *ClusterState) GetError(service string) (uint64, error) {

	cs.reMutex.Lock()
	defer cs.reMutex.Unlock()

	if _, ok := cs.errorMap[service]; ok {
		return cs.errorMap[service], nil
	}
	return 0, errors.New("harmonic: service " + service + " not found")
}

// IncrementError increments errorcount by 1 for a service.
func (cs *ClusterState) IncrementError(service string) error {

	cs.reMutex.Lock()
	defer cs.reMutex.Unlock()

	if _, ok := cs.errorMap[service]; ok {
		cs.errorMap[service]++
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// UpdateError updates user provided errorcount for a service.
func (cs *ClusterState) UpdateError(service string, errorcount uint64) error {

	cs.reMutex.Lock()
	defer cs.reMutex.Unlock()

	if _, ok := cs.errorMap[service]; ok {
		cs.errorMap[service] = errorcount
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// ResetError resets errorcount for a service.
func (cs *ClusterState) ResetError(service string) error {

	cs.reMutex.Lock()
	defer cs.reMutex.Unlock()

	if _, ok := cs.errorMap[service]; ok {
		cs.errorMap[service] = 0
		return nil
	}
	return errors.New("harmonic: service " + service + " not found")
}

// ResetAllErrors resets errorcount for all services.
func (cs *ClusterState) ResetAllErrors() error {

	cs.reMutex.Lock()
	defer cs.reMutex.Unlock()

	for k, _ := range cs.errorMap {
		cs.errorMap[k] = 0
	}
	return nil
}
