package session

import "io"
import "testing"
import "net/http"
import "net/http/httptest"
import "github.com/boyxp/nova/register"

func TestSessionIdHeader(t *testing.T) {
	var ssid string
	handler := func(w http.ResponseWriter, r *http.Request) {
		register.SetRequest(r)

		ssid = Id()
		io.WriteString(w, ssid)
	}

	w   := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.Header.Set("X-SESSID", "TEST-SESSION-HEADER-ID")

	handler(w, req)

	resp    := w.Result()
	body, _ := io.ReadAll(resp.Body)
	res     := string(body)

	if res=="TEST-SESSION-HEADER-ID" {
		t.Log(res)
	} else {
		t.Fail()
	}
}

func TestSessionIdCookie(t *testing.T) {
	var ssid string
	handler := func(w http.ResponseWriter, r *http.Request) {
		register.SetRequest(r)

		ssid = Id()
		io.WriteString(w, ssid)
	}

	cookie := &http.Cookie{
		Name    : "GOSESSID",
		Value   : "TEST-SESSION-COOOKIE-ID",
	}

	w   := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.AddCookie(cookie)

	handler(w, req)

	resp    := w.Result()
	body, _ := io.ReadAll(resp.Body)
	res     := string(body)

	if res=="TEST-SESSION-COOOKIE-ID" {
		t.Log(res)
	} else {
		t.Fail()
	}
}

func TestSetGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		register.SetRequest(r)

		Set("user", "lee")
		user := Get("user")

		if user=="lee" {
			t.Log(user)
		} else {
			t.Fail()
		}
	}

	w   := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost", nil)

	handler(w, req)
}
