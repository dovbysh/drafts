package sequence

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type Sequence struct {
	name      string
	filePath  string
	i         int64
	saveMtx   sync.Mutex
	lastSaved int64
	DiskSaved int64
}

func NewSequence(name string, dir string) *Sequence {
	s := &Sequence{name: name}
	s.filePath = fmt.Sprintf("%s/%x.txt", dir, sha256.Sum256([]byte(fmt.Sprintf("%s_%s", name, dir))))
	f, err := os.Open(s.filePath)
	if err == nil {
		buf := make([]byte, 65536)
		_, err = f.Read(buf)
		if err != nil {
			return s
		}
		if f.Close() != nil {
			return s
		}
		st := strings.Split(string(buf), "\n")
		var i int64
		i, err = strconv.ParseInt(st[0], 10, 64)
		if err != nil {
			return s
		}
		if name != st[1] {
			return s
		}
		if fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d%s", i, s.name)))) != st[2] {
			return s
		}
		s.i = i
	}

	return s
}

func (s Sequence) Name() string {
	return s.name
}

// Each goroutine get atomic value
func (s *Sequence) Next() (int64, error) {
	return atomic.AddInt64(&s.i, 1), s.safeSave()
}

func (s *Sequence) Current() int64 {
	return atomic.LoadInt64(&s.i)
}

func (s *Sequence) safeSave() error {
	s.saveMtx.Lock()
	defer s.saveMtx.Unlock()

	return s.save()
}

func (s *Sequence) save() error {
	i := atomic.LoadInt64(&s.i)
	if i == s.lastSaved {
		return nil
	}

	f, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%d\n%s\n%x\n", i, s.name, sha256.Sum256([]byte(fmt.Sprintf("%d%s", i, s.name)))))
	if err != nil {
		return err
	}
	s.lastSaved = i
	s.DiskSaved++

	return nil
}
