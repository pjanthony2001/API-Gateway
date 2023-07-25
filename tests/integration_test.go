package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestIntegration(t *testing.T) {
	var tests = []struct {
		name             string
		url              string
		wantMessage      string
		wantFlag         int32
		additionalInputs string
	}{
		{"empty json should be empty", "/echo/query", "", 0, ""},
		{"message in json should be message", "/echo/query", "Message", 0, ""},
		{"flag in json should be flag", "/echo/query", "", 300, ""},
		{"message and flag in json should be message and flag json", "/echo/query", "Hello", 211, ""},
		{"Additional inputs with empty request should have no effect", "/echo/query", "", 0, "Additional"},
		{"Additional inputs with full request should have no effect", "/echo/query", "Hallo", 212, "Additional"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := json.Marshal(map[string]interface{}{"Message": tt.wantMessage, "Flag": tt.wantFlag, "Additional": tt.additionalInputs})
			if err != nil {
				t.Fatalf("Error from creating json data %s", err.Error())
			}
			reqBody := bytes.NewReader(data)
			req, err := http.NewRequest(http.MethodGet, hertzURL+tt.url, reqBody)
			if err != nil {
				t.Fatalf("Error from creating request %s", err.Error())
			}

			req.Header.Set("Content-Type", "application/json")

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

			if sb["Message"] != tt.wantMessage {
				t.Errorf("got %s, want %s", sb["Message"], tt.wantMessage)
			}
			if int32(sb["Flag"].(float64)) != tt.wantFlag {
				t.Errorf("got %f, want %d", sb["Flag"].(float64), tt.wantFlag)
			}
			if len(sb) != 2 {
				t.Errorf("got %d, want %d", len(sb), 2)
			}
		})
	}
}
