package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestIntegration(t *testing.T) {
	var tests = []struct {
		name             string
		url              string
		wantFlag         int32
		additionalInputs string
		sendMessage      string
		service          int32
		method           int32
	}{
		{"empty json should be empty [default]", "/echo/query", 0, "", "", 1, 1},
		{"message in json should be message [default]", "/echo/query", 0, "", "Message", 1, 1},
		{"flag in json should be flag [default]", "/echo/query", 300, "", "", 1, 1},
		{"message and flag in json should be message and flag json [default]", "/echo/query", 211, "", "Hello", 1, 1},
		{"Additional inputs with empty request should have no effect [default]", "/echo/query", 0, "Additional", "", 1, 1},
		{"Additional inputs with full request should have no effect [default]", "/echo/query", 212, "Additional", "Hallo", 1, 1},

		{"empty json should be empty [service=1]", "/echo/query?service=1", 0, "", "", 1, 1},
		{"message in json should be message [service=1]", "/echo/query?service=1", 0, "", "Message", 1, 1},
		{"flag in json should be flag [service=1]", "/echo/query?service=1", 300, "", "", 1, 1},
		{"message and flag in json should be message and flag json [service=1]", "/echo/query?service=1", 211, "", "Hello", 1, 1},
		{"Additional inputs with empty request should have no effect [service=1]", "/echo/query?service=1", 0, "Additional", "", 1, 1},
		{"Additional inputs with full request should have no effect [service=1]", "/echo/query?service=1", 212, "Additional", "Hallo", 1, 1},

		{"empty json should be empty [service=2]", "/echo/query?service=2", 0, "", "", 2, 1},
		{"message in json should be message [service=2]", "/echo/query?service=2", 0, "", "Message", 2, 1},
		{"flag in json should be flag [service=2]", "/echo/query?service=2", 300, "", "", 2, 1},
		{"message and flag in json should be message and flag json [service=2]", "/echo/query?service=2", 211, "", "Hello", 2, 1},
		{"Additional inputs with empty request should have no effect [service=2]", "/echo/query?service=2", 0, "Additional", "", 2, 1},
		{"Additional inputs with full request should have no effect [service=2]", "/echo/query?service=2", 212, "Additional", "Hallo", 2, 1},

		{"empty json should be empty [method=1]", "/echo/query?method=1", 0, "", "", 1, 1},
		{"message in json should be message [method=1]", "/echo/query?method=1", 0, "", "Message", 1, 1},
		{"flag in json should be flag [method=1]", "/echo/query?method=1", 300, "", "", 1, 1},
		{"message and flag in json should be message and flag json [method=1]", "/echo/query?method=1", 211, "", "Hello", 1, 1},
		{"Additional inputs with empty request should have no effect [method=1]", "/echo/query?method=1", 0, "Additional", "", 1, 1},
		{"Additional inputs with full request should have no effect [method=1]", "/echo/query?method=1", 212, "Additional", "Hallo", 1, 1},

		{"empty json should be empty [method=2]", "/echo/query?method=2", 0, "", "", 1, 2},
		{"message in json should be message [method=2]", "/echo/query?method=2", 0, "", "Message", 1, 2},
		{"flag in json should be flag [method=2]", "/echo/query?method=2", 300, "", "", 1, 2},
		{"message and flag in json should be message and flag json [method=2]", "/echo/query?method=2", 211, "", "Hello", 1, 2},
		{"Additional inputs with empty request should have no effect [method=2]", "/echo/query?method=2", 0, "Additional", "", 1, 2},
		{"Additional inputs with full request should have no effect [method=2]", "/echo/query?method=2", 212, "Additional", "Hallo", 1, 2},

		{"empty json should be empty [service=1&method=1]", "/echo/query?service=1&method=1", 0, "", "", 1, 1},
		{"message in json should be message [service=1&method=1]", "/echo/query?service=1&method=1", 0, "", "Message", 1, 1},
		{"flag in json should be flag [service=1&method=1]", "/echo/query?service=1&method=1", 300, "", "", 1, 1},
		{"message and flag in json should be message and flag json [service=1&method=1]", "/echo/query?service=1&method=1", 211, "", "Hello", 1, 1},
		{"Additional inputs with empty request should have no effect [service=1&method=1]", "/echo/query?service=1&method=1", 0, "Additional", "", 1, 1},
		{"Additional inputs with full request should have no effect [service=1&method=1]", "/echo/query?service=1&method=1", 212, "Additional", "Hallo", 1, 1},

		{"empty json should be empty [service=1&method=2]", "/echo/query?service=1&method=2", 0, "", "", 1, 2},
		{"message in json should be message [service=1&method=2]", "/echo/query?service=1&method=2", 0, "", "Message", 1, 2},
		{"flag in json should be flag [service=1&method=2]", "/echo/query?service=1&method=2", 300, "", "", 1, 2},
		{"message and flag in json should be message and flag json [service=1&method=2]", "/echo/query?service=1&method=2", 212, "", "Hello", 1, 2},
		{"Additional inputs with empty request should have no effect [service=1&method=2]", "/echo/query?service=1&method=2", 0, "Additional", "", 1, 2},
		{"Additional inputs with full request should have no effect [service=1&method=2]", "/echo/query?service=1&method=2", 212, "Additional", "Hallo", 1, 2},

		{"empty json should be empty [service=2&method=1]", "/echo/query?service=2&method=1", 0, "", "", 2, 1},
		{"message in json should be message [service=2&method=1]", "/echo/query?service=2&method=1", 0, "", "Message", 2, 1},
		{"flag in json should be flag [service=2&method=1]", "/echo/query?service=2&method=1", 300, "", "", 2, 1},
		{"message and flag in json should be message and flag json [service=2&method=1]", "/echo/query?service=2&method=1", 212, "", "Hello", 2, 1},
		{"Additional inputs with empty request should have no effect [service=2&method=1]", "/echo/query?service=2&method=1", 0, "Additional", "", 2, 1},
		{"Additional inputs with full request should have no effect [service=2&method=1]", "/echo/query?service=2&method=1", 212, "Additional", "Hallo", 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := json.Marshal(map[string]interface{}{"Message": tt.sendMessage, "Flag": tt.wantFlag, "Additional": tt.additionalInputs})
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

			responseMessage := responseMessage(tt.service, tt.method, tt.sendMessage)

			if sb["Message"] != responseMessage {
				t.Errorf("got %s, want %s", sb["Message"], responseMessage)
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

func responseMessage(service int32, method int32, message string) string {
	return fmt.Sprintf("Parsed Message from Service %d and Method %d: %s", service, method, message)
}
