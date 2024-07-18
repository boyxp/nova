package cache

import "time"
import "testing"

func TestSetAndGet(t *testing.T) {
	Memory{}.Set("test", 1024)
	v := Memory{}.Get("test")
	if v!=1024 {
		t.FailNow()
	}

	t.Log(v)
}

func TestSetAndGetTtl(t *testing.T) {
	Memory{}.Set("test", 1024, 2)

	time.Sleep(3 * time.Second)

	v := Memory{}.Get("test")
	if v!=nil {
		t.FailNow()
	}

	t.Log(v)
}

func TestExist(t *testing.T) {
	e := Memory{}.Exist("test")
	if e!=true {
		t.FailNow()
	}

	t.Log(e)
}

func TestTtl(t *testing.T) {
	l := Memory{}.Ttl("test")

	if l<=1 {
		t.FailNow()
	}

	t.Log(l)
}

func TestDelete(t *testing.T) {
	Memory{}.Delete("test")

	e := Memory{}.Exist("test")

	if e!=false {
		t.FailNow()
	}

	t.Log(e)
}

func TestFlush(t *testing.T) {
	Memory{}.Set("test2", 2048)

	v := Memory{}.Get("test2")

	t.Log(v)

	Memory{}.Flush()

	e := Memory{}.Exist("test")

	if e != false {
		t.FailNow()
	}

	t.Log(e)
}
