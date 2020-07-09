# Harmonic

Harmonic is the request dispatch algo powering ServiceQ (https://github.com/gptankit/serviceq). It is exposed in this repo as a package with enhanced support for initializing cluster state and error management.

## Introduction

Assume a cluster of homogenous services that can serve user requests for a resource. One or more services might exhibit errors during their run ranging from connectivity loss, latency, misconfigurations and so on. Harmonic's role in this scenario is to select a service that has the best chance of processing an incoming request thus improving the total probability of success for a group of requests.

<b>How To Get</b>

<pre>
go get github/gptankit/harmonic
</pre>

Then, import the module - 

<pre>
import "github/gptankit/harmonic"
</pre>

<b>How To Use</b>

Harmonic works on user supplied cluster state which includes a list of services (can be IP, qualified url or any unique identifier). By default, harmonic will zero the error count on each service and return a reference to cluster state, that can further be used to manage error counts.

<pre>
// Initialize cluster state
servicelist := []string{"s0", "s1", "s2"}
cs, _ := harmonic.InitClusterState(servicelist)

// Call SelectService with <i>retryindex</i>=0 and <i>prevservice</i>=""
svc, _ := harmonic.SelectService(cs, 0, "")
</pre>

The request can now be forwarded to the selected <i>svc</i>. Usually, the above call to <i>SelectService</i> should be enough to select a healthy service. If the request to <i>svc</i> succeedes, reset the error count -

<pre>
// Reset error count for a service
cs.ResetError(svc)
</pre>

As a good practice, it is advisable to do retries if the request to <i>svc</i> fails. The retry call to <i>SelectService</i> can be made after incrementing error count, incrementing <i>retryindex</i> and passing previously selected <i>svc</i> - 

<pre>
// Increment error count for a service
cs.IncrementError(svc)
svc, _ := harmonic.SelectService(cs, 1, svc)
</pre>

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

Go bench results (OS: linux)

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
