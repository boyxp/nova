package cache

import "sync"
import "time"

func init() {
	go Memory{}.GC()
}

var cube sync.Map

type cell struct {
	Data any
	Ttl time.Time
}

type Memory struct {
}

func (M Memory) Set(key string, value any, ttl ...uint) bool {
	var sec uint = 600
	if len(ttl)>0 {
		sec = ttl[0]
	}

	cube.Store(key, cell{
		Data : value,
		Ttl  : time.Now().Add(time.Second * time.Duration(sec)),
	})

	return true
}

func (M Memory) Get(key string) any {
	v, ok := cube.Load(key)

	if !ok {
		return nil
	}

	if time.Now().After(v.(cell).Ttl) {
		return nil
	}

	return v.(cell).Data
}

func (M Memory) Delete(key string) bool {
	cube.Delete(key)

	return true
}

func (M Memory) Exist(key string) bool {
	_, ok := cube.Load(key)

	return ok
}

func (M Memory) Ttl(key string) uint {
	v, ok := cube.Load(key)

	if !ok {
		return 0
	}

	return uint(v.(cell).Ttl.Sub(time.Now()).Seconds())
}

func (M Memory) Flush() bool {
	cube.Range(func(key, value any) bool {
		cube.Delete(key)
		return true
	})

	return true
}

func (M Memory) GC() {
	now := time.Now()

	cube.Range(func(key, value any) bool {
		if now.After(value.(cell).Ttl) {
			cube.Delete(key)
		}
		return true
	})

	time.Sleep(600 * time.Second)

	go M.GC()
}
