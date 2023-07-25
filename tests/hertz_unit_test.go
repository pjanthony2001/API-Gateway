package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

const hertzURL = "http://127.0.0.1:8080"

func TestHertz(t *testing.T) {
	var tests = []struct {
		name  string
		query string
		data  string
		want  string
	}{
		{"/ping should be pong", "/ping", "", "pong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reqBody := strings.NewReader(tt.data)
			req, err := http.NewRequest(http.MethodGet, hertzURL+tt.query, reqBody)
			if err != nil {
				t.Fatalf("Error from creating request %s", err.Error())
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Error from making request %s", err.Error())
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error from reading response %s", err.Error())
			}

			t.Log(string(body))
			var sb map[string]interface{}
			json.Unmarshal(body, &sb)

			if sb["message"] != tt.want {
				t.Errorf("got %s, want %s", sb, tt.want)
			}
		})
	}
}
