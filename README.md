# Harmonic

Harmonic is the request dispatch algo powering ServiceQ (https://github.com/gptankit/serviceq). It is exposed in this repo as a package with enhanced support for initializing cluster state and error management.

Assume a cluster of homogenous services that can serve user requests for a resource. One or more services might exhibit errors during their run ranging from connectivity loss, latency, misconfigurations and so on. Harmonic's role in this scenario is to select a service that has the best chance of processing an incoming request thus improving the total probability of success for a group of requests.

<b>How To Import</b>

<pre>
import "github/gptankit/harmonic"
</pre>

<b>How To Use</b>

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

Feel free to play around and post feedbacks
