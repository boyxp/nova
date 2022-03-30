package register

import "sync"
import "strings"
import "runtime"
import "net/http"

var cache sync.Map

func SetResponseWriter(w http.ResponseWriter) bool {
    _rid := getRoutineId()
    _key := "w"+_rid

    cache.Store(_key, w)

    return true
}

func SetRequest(r *http.Request) bool {
    _rid := getRoutineId()
    _key := "r"+_rid

    cache.Store(_key, r)

    return true
}

func GetResponseWriter() http.ResponseWriter {
    _rid := getRoutineId()
    _key := "w"+_rid

    value, ok := cache.Load(_key)
    if !ok {
        return nil
    }

    writer, ok := value.(http.ResponseWriter)
    if !ok {
        return nil
    }

    return writer
}

func GetRequest() *http.Request {
    _rid := getRoutineId()
    _key := "r"+_rid

    value, ok := cache.Load(_key)
    if !ok {
        return nil
    }

    request, ok := value.(*http.Request)
    if !ok {
        return nil
    }

    return request
}

func Clean() bool {
    _rid := getRoutineId()

    cache.Delete("w"+_rid)
    cache.Delete("r"+_rid)

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
