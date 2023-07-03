package transport

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubController struct{}

func (m *stubController) GetRate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getRate"))
}

func (m *stubController) SubscribeEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("subscribeEmail"))
}

func (m *stubController) SendEmails(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sendEmails"))
}

func TestHttpRouter(t *testing.T) {
	mux := http.NewServeMux()
	controller := &stubController{}
	router := NewHTTPRouter(controller)
	router.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	tests := []struct {
		name  string
		route string
		want  string
	}{
		{name: "Test rate", route: "/api/rate", want: "getRate"},
		{name: "Test subscribe", route: "/api/subscribe", want: "subscribeEmail"},
		{name: "Test sendEmails", route: "/api/sendEmails", want: "sendEmails"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Get(server.URL + tt.route)
			require.NoError(t, err)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			got := string(body)
			require.Equal(t, tt.want, got)
		})
	}
}
