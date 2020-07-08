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

func TestInitClusterState(t *testing.T) {

	_, err := InitClusterState([]string{"s0","s1","s2"})
	if err != nil{
		t.Errorf("init cluster failed")
	}
}

func TestIncrementError(t *testing.T){

	cs, err := InitClusterState([]string{"s0","s1","s2"})
	if err != nil{
		return
	}

	err1 := cs.IncrementError("s1")
	if err1 != nil{
		t.Errorf("increment error failed")
	}

	err2 := cs.IncrementError("s5")
	if err2 == nil{
		t.Errorf("non-existent service error incremented")
	}
}

func TestUpdateError(t *testing.T){

	cs, err := InitClusterState([]string{"s0","s1","s2"})
	if err != nil{
		return
	}

	err1 := cs.UpdateError("s1", 10)
	if err1 != nil{
		t.Errorf("update error failed")
	}

	err2 := cs.UpdateError("s5", 10)
	if err2 == nil{
		t.Errorf("non-existent service error updated")
	}
}

func TestResetError(t *testing.T){

	cs, err := InitClusterState([]string{"s0","s1","s2"})
	if err != nil{
		return
	}

	err1 := cs.ResetError("s1")
	if err1 != nil{
		t.Errorf("reset error failed")
	}

	err2 := cs.ResetError("s5")
	if err2 == nil{
		t.Errorf("non-existent service error reset")
	}
}

func BenchmarkSelectServiceR1(b *testing.B){

	cs, err := InitClusterState([]string{"s0","s1","s2"})
	if err != nil{
		return
	}

	for i:=0;i<b.N;i++{
		SelectService(cs, 0, "s1") // first try
	}
}

func BenchmarkSelectServiceRn(b *testing.B){

	cs, err := InitClusterState([]string{"s0","s1","s2"})
	if err != nil{
		return
	}

	for i:=0;i<b.N;i++{
		SelectService(cs, 2, "s1") // subsequent try
	}
}
