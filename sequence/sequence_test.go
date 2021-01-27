package sequence

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewSequence(t *testing.T) {
	const name = "load test"
	fp := fmt.Sprintf("%s/zzz%d", os.TempDir(), rand.Uint64())
	var err error
	err = os.MkdirAll(fp, 0700)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fp)

	{
		s := NewSequence(name, fp)
		i, err := s.Next()
		assert.Equal(t, int64(1), i)
		assert.NoError(t, err)
		assert.NoError(t, s.save())
	}
	{
		s := NewSequence(name, fp)
		assert.Equal(t, int64(1), s.Current())
	}
}

func TestNext(t *testing.T) {
	const name = "load test TestNext"
	fp := fmt.Sprintf("%s/zzz%d", os.TempDir(), rand.Uint64())
	var err error
	err = os.MkdirAll(fp, 0700)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fp)
	var sum, expectedSum, runs, waits int64
	const maxIterations = 100000
	{
		s := NewSequence(name, fp)
		ch := make(chan int64, 1000000)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range ch {
				sum += i
			}
		}()
		var wg2 sync.WaitGroup
		wg2.Add(maxIterations)
		for i := 1; i < maxIterations+1; i++ {
			expectedSum += int64(i)
			go func() {
				defer wg2.Done()
				atomic.AddInt64(&runs, 1)
				defer atomic.AddInt64(&runs, -1)
				z, err := s.Next()
				assert.NoError(t, err)
				ch <- z
			}()
			for atomic.LoadInt64(&runs) > 1000000 {
				runtime.Gosched()
				waits++
			}
		}
		wg2.Wait()
		close(ch)
		wg.Wait()
		t.Log("disk saved times: ", s.DiskSaved)
		assert.Equal(t, expectedSum, sum)
		t.Log("waits: ", waits)
	}
	{
		s := NewSequence(name, fp)
		assert.Equal(t, int64(maxIterations), s.Current())
	}
}
