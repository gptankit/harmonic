package harmonic

import "testing"

func TestServiceIndexWithinBoundR1(t *testing.T) {

	var params = []struct {
		cs *ClusterState
		rt int
		ic string
	}{
		{&ClusterState{
			errormap: map[string]uint64{"s0": 1, "s1": 2},
			servicelist: []string{"s0","s1"},
			numservices: 2,
					},
					0,
					"s5"},
		{&ClusterState{
			errormap: map[string]uint64{"s0": 1, "s1": 2},
			servicelist: []string{"s0","s1"},
			numservices: 2,
					},
					0,
					"s8"},
		{&ClusterState{
			errormap: map[string]uint64{"s0": 1, "s1": 2},
			servicelist: []string{"s0","s1"},
			numservices: 2,
					},
					0,
					"s1"},
		{&ClusterState{
			errormap: map[string]uint64{"s0": 1, "s1": 2},
			servicelist: []string{"s0","s1"},
			numservices: 2,
					},
					0,
					"s2"},
	}

	for _, prm := range params {
		ce, err := SelectService(prm.cs, prm.rt, prm.ic)
		ns := len((*prm.cs).servicelist)
		if (ns > 0 && err != nil) || ns <= 0{
			t.Errorf("service index out of bound, ns=%d, rt=%d, ic=%s --> ce=%s\n", ns, prm.rt, prm.ic, ce)
		}
	}
}

func TestServiceIndexWithinBoundRn(t *testing.T) {

	var params = []struct {
		cs *ClusterState
		rt int
		ic string
	}{
		{&ClusterState{
			errormap: map[string]uint64{"s0": 1, "s1": 2},
			servicelist: []string{"s0","s1"},
			numservices: 2,
					},
					3,
					"s59"},
		{&ClusterState{
			errormap: map[string]uint64{"s0": 1, "s1": 2},
			servicelist: []string{"s0","s1"},
			numservices: 2,
					},
					2,
					"s8"},
		{&ClusterState{
			errormap: map[string]uint64{"s0": 1, "s1": 2},
			servicelist: []string{"s0","s1"},
			numservices: 2,
					},
					-1,
					"s14"},
		{&ClusterState{
			errormap: map[string]uint64{"s0": 1, "s1": 2},
			servicelist: []string{"s0","s1"},
			numservices: 2,
					},
					-1,
					"s20"},
	}

	for _, prm := range params {
		ce, err := SelectService(prm.cs, prm.rt, prm.ic)

		ns := len((*prm.cs).servicelist)
		if (ns > 0 && err != nil) || ns <= 0 {
			t.Errorf("service index out of bound, ns=%d, rt=%d, ic=%s --> ce=%s\n", ns, prm.rt, prm.ic, ce)
		}
	}
}
