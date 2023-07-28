package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestIntegration(t *testing.T) {
	var tests = []struct {
		name             string
		wantFlag         int32
		additionalInputs string
		sendMessage      string
		wantService      int32
		wantMethod       int32
		token            string
		sendService      int32
		sendMethod       int32
	}{
		{"empty json should be empty [default]()", 0, "", "", 1, 1, "", 0, 0},
		{"message in json should be message [default]()", 0, "", "Message", 1, 1, "", 0, 0},
		{"flag in json should be flag [default]()", 300, "", "", 1, 1, "", 0, 0},
		{"message and flag in json should be message and flag json [default]()", 211, "", "Hello", 1, 1, "", 0, 0},
		{"Additional inputs with empty request should have no effect [default]()", 0, "Additional", "", 1, 1, "", 0, 0},
		{"Additional inputs with full request should have no effect [default]()", 212, "Additional", "Hallo", 1, 1, "", 0, 0},

		{"empty json should be empty [service=1]()", 0, "", "", 1, 1, "", 1, 0},
		{"message in json should be message [service=1]()", 0, "", "Message", 1, 1, "", 1, 0},
		{"flag in json should be flag [service=1]()", 300, "", "", 1, 1, "", 1, 0},
		{"message and flag in json should be message and flag json [service=1]()", 211, "", "Hello", 1, 1, "", 1, 0},
		{"Additional inputs with empty request should have no effect [service=1]()", 0, "Additional", "", 1, 1, "", 1, 0},
		{"Additional inputs with full request should have no effect [service=1]()", 212, "Additional", "Hallo", 1, 1, "", 1, 0},

		{"empty json should be empty [service=2]()", 0, "", "", 2, 1, "", 2, 0},
		{"message in json should be message [service=2]()", 0, "", "Message", 2, 1, "", 2, 0},
		{"flag in json should be flag [service=2]()", 300, "", "", 2, 1, "", 2, 0},
		{"message and flag in json should be message and flag json [service=2]()", 211, "", "Hello", 2, 1, "", 2, 0},
		{"Additional inputs with empty request should have no effect [service=2]()", 0, "Additional", "", 2, 1, "", 2, 0},
		{"Additional inputs with full request should have no effect [service=2]()", 212, "Additional", "Hallo", 2, 1, "", 2, 0},

		{"empty json should be empty [method=1]()", 0, "", "", 1, 1, "", 0, 1},
		{"message in json should be message [method=1]()", 0, "", "Message", 1, 1, "", 0, 1},
		{"flag in json should be flag [method=1]()", 300, "", "", 1, 1, "", 0, 1},
		{"message and flag in json should be message and flag json [method=1]()", 211, "", "Hello", 1, 1, "", 0, 1},
		{"Additional inputs with empty request should have no effect [method=1]()", 0, "Additional", "", 1, 1, "", 0, 1},
		{"Additional inputs with full request should have no effect [method=1]()", 212, "Additional", "Hallo", 1, 1, "", 0, 1},

		{"empty json should be empty [method=2]()", 0, "", "", 1, 2, "", 0, 2},
		{"message in json should be message [method=2]()", 0, "", "Message", 1, 2, "", 0, 2},
		{"flag in json should be flag [method=2]()", 300, "", "", 1, 2, "", 0, 2},
		{"message and flag in json should be message and flag json [method=2]()", 211, "", "Hello", 1, 2, "", 0, 2},
		{"Additional inputs with empty request should have no effect [method=2]()", 0, "Additional", "", 1, 2, "", 0, 2},
		{"Additional inputs with full request should have no effect [method=2]()", 212, "Additional", "Hallo", 1, 2, "", 0, 2},

		{"empty json should be empty [service=1&method=1]()", 0, "", "", 1, 1, "", 1, 1},
		{"message in json should be message [service=1&method=1]()", 0, "", "Message", 1, 1, "", 1, 1},
		{"flag in json should be flag [service=1&method=1]()", 300, "", "", 1, 1, "", 1, 1},
		{"message and flag in json should be message and flag json [service=1&method=1]()", 211, "", "Hello", 1, 1, "", 1, 1},
		{"Additional inputs with empty request should have no effect [service=1&method=1]()", 0, "Additional", "", 1, 1, "", 1, 1},
		{"Additional inputs with full request should have no effect [service=1&method=1]()", 212, "Additional", "Hallo", 1, 1, "", 1, 1},

		{"empty json should be empty [service=1&method=2]()", 0, "", "", 1, 2, "", 1, 2},
		{"message in json should be message [service=1&method=2]()", 0, "", "Message", 1, 2, "", 1, 2},
		{"flag in json should be flag [service=1&method=2]()", 300, "", "", 1, 2, "", 1, 2},
		{"message and flag in json should be message and flag json [service=1&method=2]()", 212, "", "Hello", 1, 2, "", 1, 2},
		{"Additional inputs with empty request should have no effect [service=1&method=2]()", 0, "Additional", "", 1, 2, "", 1, 2},
		{"Additional inputs with full request should have no effect [service=1&method=2]()", 212, "Additional", "Hallo", 1, 2, "", 1, 2},

		{"empty json should be empty [service=2&method=1]()", 0, "", "", 2, 1, "", 2, 1},
		{"message in json should be message [service=2&method=1]()", 0, "", "Message", 2, 1, "", 2, 1},
		{"flag in json should be flag [service=2&method=1]()", 300, "", "", 2, 1, "", 2, 1},
		{"message and flag in json should be message and flag json [service=2&method=1]()", 212, "", "Hello", 2, 1, "", 2, 1},
		{"Additional inputs with empty request should have no effect [service=2&method=1]()", 0, "Additional", "", 2, 1, "", 2, 1},
		{"Additional inputs with full request should have no effect [service=2&method=1]()", 212, "Additional", "Hallo", 2, 1, "", 2, 1},

		{"empty json should be empty [default](WrongToken)", 0, "", "", 1, 1, "WrongToken", 0, 0},
		{"message in json should be message [default](WrongToken)", 0, "", "Message", 1, 1, "WrongToken", 0, 0},
		{"flag in json should be flag [default](WrongToken)", 300, "", "", 1, 1, "WrongToken", 0, 0},
		{"message and flag in json should be message and flag json [default](WrongToken)", 211, "", "Hello", 1, 1, "WrongToken", 0, 0},
		{"Additional inputs with empty request should have no effect [default](WrongToken)", 0, "Additional", "", 1, 1, "WrongToken", 0, 0},
		{"Additional inputs with full request should have no effect [default](WrongToken)", 212, "Additional", "Hallo", 1, 1, "WrongToken", 0, 0},

		{"empty json should be empty [service=1](WrongToken)", 0, "", "", 1, 1, "WrongToken", 1, 0},
		{"message in json should be message [service=1](WrongToken)", 0, "", "Message", 1, 1, "WrongToken", 1, 0},
		{"flag in json should be flag [service=1](WrongToken)", 300, "", "", 1, 1, "WrongToken", 1, 0},
		{"message and flag in json should be message and flag json [service=1](WrongToken)", 211, "", "Hello", 1, 1, "WrongToken", 1, 0},
		{"Additional inputs with empty request should have no effect [service=1](WrongToken)", 0, "Additional", "", 1, 1, "WrongToken", 1, 0},
		{"Additional inputs with full request should have no effect [service=1](WrongToken)", 212, "Additional", "Hallo", 1, 1, "WrongToken", 1, 0},

		{"empty json should be empty [service=2](WrongToken)", 0, "", "", 2, 1, "WrongToken", 2, 0},
		{"message in json should be message [service=2](WrongToken)", 0, "", "Message", 2, 1, "WrongToken", 2, 0},
		{"flag in json should be flag [service=2](WrongToken)", 300, "", "", 2, 1, "WrongToken", 2, 0},
		{"message and flag in json should be message and flag json [service=2](WrongToken)", 211, "", "Hello", 2, 1, "WrongToken", 2, 0},
		{"Additional inputs with empty request should have no effect [service=2](WrongToken)", 0, "Additional", "", 2, 1, "WrongToken", 2, 0},
		{"Additional inputs with full request should have no effect [service=2](WrongToken)", 212, "Additional", "Hallo", 2, 1, "WrongToken", 2, 0},

		{"empty json should be empty [method=1](WrongToken)", 0, "", "", 1, 1, "WrongToken", 0, 1},
		{"message in json should be message [method=1](WrongToken)", 0, "", "Message", 1, 1, "WrongToken", 0, 1},
		{"flag in json should be flag [method=1](WrongToken)", 300, "", "", 1, 1, "WrongToken", 0, 1},
		{"message and flag in json should be message and flag json [method=1](WrongToken)", 211, "", "Hello", 1, 1, "WrongToken", 0, 1},
		{"Additional inputs with empty request should have no effect [method=1](WrongToken)", 0, "Additional", "", 1, 1, "WrongToken", 0, 1},
		{"Additional inputs with full request should have no effect [method=1](WrongToken)", 212, "Additional", "Hallo", 1, 1, "WrongToken", 0, 1},

		{"empty json should be empty [method=2](WrongToken)", 0, "", "", 1, 2, "WrongToken", 0, 2},
		{"message in json should be message [method=2](WrongToken)", 0, "", "Message", 1, 2, "WrongToken", 0, 2},
		{"flag in json should be flag [method=2](WrongToken)", 300, "", "", 1, 2, "WrongToken", 0, 2},
		{"message and flag in json should be message and flag json [method=2](WrongToken)", 211, "", "Hello", 1, 2, "WrongToken", 0, 2},
		{"Additional inputs with empty request should have no effect [method=2](WrongToken)", 0, "Additional", "", 1, 2, "WrongToken", 0, 2},
		{"Additional inputs with full request should have no effect [method=2](WrongToken)", 212, "Additional", "Hallo", 1, 2, "WrongToken", 0, 2},

		{"empty json should be empty [service=1&method=1](WrongToken)", 0, "", "", 1, 1, "WrongToken", 1, 1},
		{"message in json should be message [service=1&method=1](WrongToken)", 0, "", "Message", 1, 1, "WrongToken", 1, 1},
		{"flag in json should be flag [service=1&method=1](WrongToken)", 300, "", "", 1, 1, "WrongToken", 1, 1},
		{"message and flag in json should be message and flag json [service=1&method=1](WrongToken)", 211, "", "Hello", 1, 1, "WrongToken", 1, 1},
		{"Additional inputs with empty request should have no effect [service=1&method=1](WrongToken)", 0, "Additional", "", 1, 1, "WrongToken", 1, 1},
		{"Additional inputs with full request should have no effect [service=1&method=1](WrongToken)", 212, "Additional", "Hallo", 1, 1, "WrongToken", 1, 1},

		{"empty json should be empty [service=1&method=2](WrongToken)", 0, "", "", 1, 2, "WrongToken", 1, 2},
		{"message in json should be message [service=1&method=2](WrongToken)", 0, "", "Message", 1, 2, "WrongToken", 1, 2},
		{"flag in json should be flag [service=1&method=2](WrongToken)", 300, "", "", 1, 2, "WrongToken", 1, 2},
		{"message and flag in json should be message and flag json [service=1&method=2](WrongToken)", 212, "", "Hello", 1, 2, "WrongToken", 1, 2},
		{"Additional inputs with empty request should have no effect [service=1&method=2](WrongToken)", 0, "Additional", "", 1, 2, "WrongToken", 1, 2},
		{"Additional inputs with full request should have no effect [service=1&method=2](WrongToken)", 212, "Additional", "Hallo", 1, 2, "WrongToken", 1, 2},

		{"empty json should be empty [service=2&method=1](WrongToken)", 0, "", "", 2, 1, "WrongToken", 2, 1},
		{"message in json should be message [service=2&method=1](WrongToken)", 0, "", "Message", 2, 1, "WrongToken", 2, 1},
		{"flag in json should be flag [service=2&method=1](WrongToken)", 300, "", "", 2, 1, "WrongToken", 2, 1},
		{"message and flag in json should be message and flag json [service=2&method=1](WrongToken)", 212, "", "Hello", 2, 1, "WrongToken", 2, 1},
		{"Additional inputs with empty request should have no effect [service=2&method=1](WrongToken)", 0, "Additional", "", 2, 1, "WrongToken", 2, 1},
		{"Additional inputs with full request should have no effect [service=2&method=1](WrongToken)", 212, "Additional", "Hallo", 2, 1, "WrongToken", 2, 1},

		{"empty json should be empty [default](token)", 0, "", "", 1, 1, "token", 0, 0},
		{"message in json should be message [default](token)", 0, "", "Message", 1, 1, "token", 0, 0},
		{"flag in json should be flag [default](token)", 300, "", "", 1, 1, "token", 0, 0},
		{"message and flag in json should be message and flag json [default](token)", 211, "", "Hello", 1, 1, "token", 0, 0},
		{"Additional inputs with empty request should have no effect [default](token)", 0, "Additional", "", 1, 1, "token", 0, 0},
		{"Additional inputs with full request should have no effect [default](token)", 212, "Additional", "Hallo", 1, 1, "token", 0, 0},

		{"empty json should be empty [service=1](token)", 0, "", "", 1, 1, "token", 1, 0},
		{"message in json should be message [service=1](token)", 0, "", "Message", 1, 1, "token", 1, 0},
		{"flag in json should be flag [service=1](token)", 300, "", "", 1, 1, "token", 1, 0},
		{"message and flag in json should be message and flag json [service=1](token)", 211, "", "Hello", 1, 1, "token", 1, 0},
		{"Additional inputs with empty request should have no effect [service=1](token)", 0, "Additional", "", 1, 1, "token", 1, 0},
		{"Additional inputs with full request should have no effect [service=1](token)", 212, "Additional", "Hallo", 1, 1, "token", 1, 0},

		{"empty json should be empty [service=2](token)", 0, "", "", 2, 1, "token", 2, 0},
		{"message in json should be message [service=2](token)", 0, "", "Message", 2, 1, "token", 2, 0},
		{"flag in json should be flag [service=2](token)", 300, "", "", 2, 1, "token", 2, 0},
		{"message and flag in json should be message and flag json [service=2](token)", 211, "", "Hello", 2, 1, "token", 2, 0},
		{"Additional inputs with empty request should have no effect [service=2](token)", 0, "Additional", "", 2, 1, "token", 2, 0},
		{"Additional inputs with full request should have no effect [service=2](token)", 212, "Additional", "Hallo", 2, 1, "token", 2, 0},

		{"empty json should be empty [method=1](token)", 0, "", "", 1, 1, "token", 0, 1},
		{"message in json should be message [method=1](token)", 0, "", "Message", 1, 1, "token", 0, 1},
		{"flag in json should be flag [method=1](token)", 300, "", "", 1, 1, "token", 0, 1},
		{"message and flag in json should be message and flag json [method=1](token)", 211, "", "Hello", 1, 1, "token", 0, 1},
		{"Additional inputs with empty request should have no effect [method=1](token)", 0, "Additional", "", 1, 1, "token", 0, 1},
		{"Additional inputs with full request should have no effect [method=1](token)", 212, "Additional", "Hallo", 1, 1, "token", 0, 1},

		{"empty json should be empty [method=2](token)", 0, "", "", 1, 2, "token", 0, 2},
		{"message in json should be message [method=2](token)", 0, "", "Message", 1, 2, "token", 0, 2},
		{"flag in json should be flag [method=2](token)", 300, "", "", 1, 2, "token", 0, 2},
		{"message and flag in json should be message and flag json [method=2](token)", 211, "", "Hello", 1, 2, "token", 0, 2},
		{"Additional inputs with empty request should have no effect [method=2](token)", 0, "Additional", "", 1, 2, "token", 0, 2},
		{"Additional inputs with full request should have no effect [method=2](token)", 212, "Additional", "Hallo", 1, 2, "token", 0, 2},

		{"empty json should be empty [service=1&method=1](token)", 0, "", "", 1, 1, "token", 1, 1},
		{"message in json should be message [service=1&method=1](token)", 0, "", "Message", 1, 1, "token", 1, 1},
		{"flag in json should be flag [service=1&method=1](token)", 300, "", "", 1, 1, "token", 1, 1},
		{"message and flag in json should be message and flag json [service=1&method=1](token)", 211, "", "Hello", 1, 1, "token", 1, 1},
		{"Additional inputs with empty request should have no effect [service=1&method=1](token)", 0, "Additional", "", 1, 1, "token", 1, 1},
		{"Additional inputs with full request should have no effect [service=1&method=1](token)", 212, "Additional", "Hallo", 1, 1, "token", 1, 1},

		{"empty json should be empty [service=1&method=2](token)", 0, "", "", 1, 2, "token", 1, 2},
		{"message in json should be message [service=1&method=2](token)", 0, "", "Message", 1, 2, "token", 1, 2},
		{"flag in json should be flag [service=1&method=2](token)", 300, "", "", 1, 2, "token", 1, 2},
		{"message and flag in json should be message and flag json [service=1&method=2](token)", 212, "", "Hello", 1, 2, "token", 1, 2},
		{"Additional inputs with empty request should have no effect [service=1&method=2](token)", 0, "Additional", "", 1, 2, "token", 1, 2},
		{"Additional inputs with full request should have no effect [service=1&method=2](token)", 212, "Additional", "Hallo", 1, 2, "token", 1, 2},

		{"empty json should be empty [service=2&method=1](token)", 0, "", "", 2, 1, "token", 2, 1},
		{"message in json should be message [service=2&method=1](token)", 0, "", "Message", 2, 1, "token", 2, 1},
		{"flag in json should be flag [service=2&method=1](token)", 300, "", "", 2, 1, "token", 2, 1},
		{"message and flag in json should be message and flag json [service=2&method=1](token)", 212, "", "Hello", 2, 1, "token", 2, 1},
		{"Additional inputs with empty request should have no effect [service=2&method=1](token)", 0, "Additional", "", 2, 1, "token", 2, 1},
		{"Additional inputs with full request should have no effect [service=2&method=1](token)", 212, "Additional", "Hallo", 2, 1, "token", 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := json.Marshal(map[string]interface{}{
				"Message":    tt.sendMessage,
				"Flag":       tt.wantFlag,
				"Additional": tt.additionalInputs,
			})

			if err != nil {
				t.Fatalf("Error from creating json data %s", err.Error())
			}

			url := URL(tt.sendService, tt.sendMethod, tt.token)

			t.Log(url)

			reqBody := bytes.NewReader(data)
			req, err := http.NewRequest(http.MethodGet, url, reqBody)
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

			responseMessage, keyAuthErr := responseMessage(tt.wantService, tt.wantMethod, tt.sendMessage, tt.token)

			if keyAuthErr {
				if string(body) != responseMessage {
					t.Errorf("got %s, want %s", string(body), responseMessage)
				}
			} else {
				if sb["Message"] != responseMessage {
					t.Errorf("got %s, want %s", sb["Message"], responseMessage)
				}
				if int32(sb["Flag"].(float64)) != tt.wantFlag {
					t.Errorf("got %f, want %d", sb["Flag"].(float64), tt.wantFlag)
				}
				if len(sb) != 2 {
					t.Errorf("got %d, want %d", len(sb), 2)
				}
			}
		})
	}
}

func responseMessage(service int32, method int32, message string, token string) (string, bool) {
	if service == 2 {
		if token == "" {
			return "Error by keyauth: missing or malformed API Key", true
		}

		if token != "token" {
			return fmt.Sprintf("Error by keyauth unauthorized: %s", token), true
		}
	}

	return fmt.Sprintf("Parsed Message from Service %d and Method %d: %s", service, method, message), false
}

func URL(service int32, method int32, token string) string {
	values := url.Values{}

	if service != 0 {
		values.Add("service", strconv.FormatInt(int64(service), 10))
	}
	if method != 0 {
		values.Add("method", strconv.FormatInt(int64(method), 10))
	}
	if token != "" {
		values.Add("token", token)
	}

	url := url.URL{
		Scheme:   "http",
		RawQuery: values.Encode(),
		Path:     "echo/query",
		Host:     "localhost:8080",
	}

	return url.String()
}
