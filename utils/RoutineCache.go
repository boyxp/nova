package RoutineCache

import "sync"
import "strings"
import "runtime"

var cache sync.Map

func Set(key string, value interface{}) bool {
    if strings.Contains(key, ":") {
        return false
    }

    _rid := getRoutineId()
    _key := "r"+_rid+":"+key

    cache.Store(_key, value)

    return true
}

func Get(key string) interface{} {
    _rid := getRoutineId()
    _key := "r"+_rid+":"+key

    value, ok := cache.Load(_key)
    if !ok {
        return nil
    }

    return value
}

func Delete(key string) bool {
    _rid := getRoutineId()
    _key := "r"+_rid+":"+key

    cache.Delete(_key)

    return true
}

func getRoutineId() string {
    var (
        buf [64]byte
        n   = runtime.Stack(buf[:], false)
        stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
    )

    return strings.Fields(stk)[0]
}
