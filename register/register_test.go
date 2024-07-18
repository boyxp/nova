package register

import "testing"
import "net/http/httptest"

func TestRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.Header.Set("Test", "OK")

	SetRequest(req)

	tmp := GetRequest()

	if tmp.Header.Get("Test")=="OK" {
		t.Log(tmp)
	} else {
		t.Fail()
	}
}

func TestResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	w.Header().Set("Test", "OK")

	SetResponseWriter(w)

	tmp := GetResponseWriter()

	if tmp.Header().Get("Test")=="OK" {
		t.Log(tmp)
	} else {
		t.Fail()
	}
}
