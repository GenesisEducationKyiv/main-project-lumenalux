package transport

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockController struct{}

func (m *mockController) GetRate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getRate"))
}

func (m *mockController) SubscribeEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("subscribeEmail"))
}

func (m *mockController) SendEmails(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sendEmails"))
}

func TestHttpRouter(t *testing.T) {
	mux := http.NewServeMux()
	controller := &mockController{}
	router := NewHTTPRouter(controller)
	router.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	tests := []struct {
		route string
		want  string
	}{
		{route: "/api/rate", want: "getRate"},
		{route: "/api/subscribe", want: "subscribeEmail"},
		{route: "/api/sendEmails", want: "sendEmails"},
	}

	for _, tt := range tests {
		res, err := http.Get(server.URL + tt.route)
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		got := string(body)
		if got != tt.want {
			t.Errorf("got %q, want %q", got, tt.want)
		}
	}
}
