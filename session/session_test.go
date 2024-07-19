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

func TestGoSessionIdCookie(t *testing.T) {
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

func TestPHPSessionIdCookie(t *testing.T) {
	var ssid string
	handler := func(w http.ResponseWriter, r *http.Request) {
		register.SetRequest(r)

		ssid = Id()
		io.WriteString(w, ssid)
	}

	cookie := &http.Cookie{
		Name    : "PHPSESSID",
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

func TestNoSessionId(t *testing.T) {
	var ssid string
	handler := func(w http.ResponseWriter, r *http.Request) {
		register.SetRequest(r)

		ssid = Id()
		io.WriteString(w, ssid)
	}

	w   := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost", nil)

	handler(w, req)

	resp    := w.Result()
	body, _ := io.ReadAll(resp.Body)
	res     := string(body)

	if len(res)>1 {
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

func TestGetAll(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		register.SetRequest(r)

		Set("user", "lee")
		Set("admin_id", "123")
		Set("entity_id", "4")
		Set("token", "abc")

		all  := All()

		if all["user"]=="lee" && all["admin_id"]=="123" && all["entity_id"]=="4" && all["token"]=="abc" {
			t.Log(all)
		} else {
			t.Fail()
		}
	}

	w   := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost", nil)

	handler(w, req)
}

func BenchmarkSetGet(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		register.SetRequest(r)

		Set("user", "lee")
		Set("admin_id", "123")
		Set("entity_id", "4")
		Set("token", "abc")

	    for i := 0; i < b.N; i++ {
			Get("user")
			Get("admin_id")
			Get("entity_id")
			Get("token")
	    }
	}

	w   := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost", nil)

	handler(w, req)
}

