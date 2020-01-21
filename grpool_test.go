package grpool

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

func TestJobs(t *testing.T) {
	var wg sync.WaitGroup
	var m sync.Map
	wg.Add(150)
	pool := New(10, func() Job {
		return &mockJob{wg: &wg, m: &m}
	})
	for i := 0; i < 150; i++ {
		pool.Process(i)
	}
	wg.Wait()
	length := 0
	m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	exp := 10
	act := length
	if exp != act {
		t.FailNow()
	}
}

type mockJob struct {
	wg *sync.WaitGroup
	m  *sync.Map
}

func (j *mockJob) Process(value interface{}) {
	j.m.Store(getGID(), "")
	j.wg.Done()
}

func (j *mockJob) Exit() {
	// nothing
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
