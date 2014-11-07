// Copyright 2012 Andrew 'Diddymus' Rolfe. All rights reserved.
// 2014 Modified by Jason Del Ponte for 'Explorer' project
//
package uid

import (
	"runtime"
	"testing"
)

const (
	loops        = 100
	countPerLoop = 100
	maxLoops     = loops * countPerLoop
)

func TestSequence(t *testing.T) {
	for x := 0; x < loops; x++ {
		have := <-Next
		want := <-Next - 1
		if have != want {
			t.Errorf("Corrupt sequence: Case %d, have %d wanted %d", x, have, want)
		}
	}
}

func TestConcurrency(t *testing.T) {

	uids := make([]UID, 0)
	results := make(chan UID, maxLoops)

	// Fire off a number of Goroutines to grab a bunch of UIDs each
	for x := 0; x < loops; x++ {
		go func(results chan UID) {
			for x := 0; x < countPerLoop; x++ {
				results <- <-Next
				runtime.Gosched()
			}
		}(results)
		runtime.Gosched()
	}

	// Wait for results
	for x := 0; x < maxLoops; x++ {
		uids = append(uids, <-results)
	}

	// Make sure all results are unique
	for x, have := range uids {
		for y, found := range uids {
			if have == found && x != y {
				t.Errorf("Duplicate UID: Cases %d & %d, have %d found %d", x, y, have, found)
			}
		}
	}
}

func TestIsAlso(t *testing.T) {

	testSubjects := [countPerLoop]UID{}
	for i := range testSubjects {
		testSubjects[i] = <-Next
	}

	// Match each test subject with every other - should only match itself
	for i1, subject1 := range testSubjects {
		for i2, subject2 := range testSubjects {
			have := subject1.IsAlso(subject2)
			want := i1 == i2
			if have != want {
				t.Errorf("Corrupt IsAlso: Case %d, have %t wanted %t", i1, have, want)
			}
		}
	}
}

func TestUniqueId(t *testing.T) {
	testSubjects := [countPerLoop]UID{}
	for i := range testSubjects {
		testSubjects[i] = <-Next
	}

	for i, s := range testSubjects {
		have := s.UniqueId()
		want := s
		if have != want {
			t.Errorf("Corrupt UniqueId: Case %d, have %d wanted %d", i, have, want)
		}
	}
}
