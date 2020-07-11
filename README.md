# Harmonic [![Build Status](https://travis-ci.com/gptankit/harmonic.svg?branch=master)](https://travis-ci.com/gptankit/harmonic) [![GoDoc](https://godoc.org/github.com/gptankit/harmonic?status.svg)](https://pkg.go.dev/github.com/gptankit/harmonic?tab=doc)

Harmonic is the request dispatch algorithm powering ServiceQ (https://github.com/gptankit/serviceq). It is exposed in this repository as a package with enhanced support for initializing cluster state and error management.

## Introduction

Assume a cluster of homogenous services that can serve user requests for a resource. One or more services might exhibit errors during their run ranging from connectivity loss, latency, misconfigurations and so on. Harmonic's role, in this scenario, is to select a service that has the best chance of processing an incoming request, thus improving the total probability of success for a group of requests.

<b>How To Get</b>

<pre>
go get github.com/gptankit/harmonic
</pre>

Then, import the module - 

<pre>
import "github.com/gptankit/harmonic"
</pre>

<b>How To Use</b>

Harmonic works on user-supplied cluster state which includes a list of services (can be IP, fully-qualified URL or any unique identifier). By default, harmonic will zero the error count on each service and return a reference to cluster state, that can further be used to manage error counts.

<pre>
// Initialize cluster state
serviceList := []string{"s0", "s1", "s2"}
cs, _ := harmonic.InitClusterState(serviceList)

// Call SelectService with <i>retryIndex</i>=0 and <i>prevService</i>=""
svc, _ := harmonic.SelectService(cs, 0, "")
</pre>

The request can now be forwarded to the selected <i>svc</i>. Usually, the above call to <i>SelectService</i> should be enough to select a healthy service. If the request to <i>svc</i> succeedes, reset the error count -

<pre>
// Reset error count for a service
cs.ResetError(svc)
</pre>

As a good practice, it is advisable to do retries if the request to <i>svc</i> fails. The retry call to <i>SelectService</i> can be made after incrementing error count, incrementing <i>retryIndex</i> and passing previously selected <i>svc</i> - 

<pre>
// Increment error count for a service
cs.IncrementError(svc)
svc, _ := harmonic.SelectService(cs, 1, svc)
</pre>

<b>Note</b>: What constitutes a <i>success</i> or <i>failure</i> response will depend on the type of application and use case, and needs to be defined accordingly by the application owner.

A complete example is added in <b>sample</b> folder.

<b>Error Management</b>

Listed are functions to manage error counts, and can be used depending on application needs - 

<pre>
func (cs *ClusterState) GetError(service string) (uint64, error)
func (cs *ClusterState) IncrementError(service string) error
func (cs *ClusterState) UpdateError(service string, errorcount uint64) error
func (cs *ClusterState) ResetError(service string) error
func (cs *ClusterState) ResetAllErrors() error
</pre>

## Benchmarks

Go bench results (env: linux/amd64)

Call to <i>SelectService</i>, first try - 

|  Ops  |  Time  |
| ----- |  ----- |
| 2000  | 0.089s |
| 20000 | 0.575s |
| 200000| 2.027s |

Avg execution time: <b>8889 ns/op</b>

Call to <i>SelectService</i>, retries - 

|  Ops  |  Time  |
| ----- |  ----- |
| 5000  | 0.063s |
| 50000 | 0.082s |
| 500000| 0.089s |

Avg execution time: <b>32.4 ns/op</b>

Feel free to play around and post feedbacks
