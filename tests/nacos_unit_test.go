package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestNacosServiceDiscovery(t *testing.T) {
	var tests = []struct {
		name     string
		query    string
		wantIP   string
		wantPort float64
	}{
		{"Testing for ExampleService1", "ExampleService1", "127.0.0.1", 8888},
		{"Testing for ExampleService2", "ExampleService2", "127.0.0.1", 8885},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8848/nacos/v1/ns/instance/list?serviceName=%s", tt.query))

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

			t.Log(sb)

			data := sb["hosts"].([]interface{})[0].(map[string]interface{})

			if data["ip"].(string) != tt.wantIP {
				t.Errorf("Wrong IP addr for instance, given: %s", data["ip"].(string))
			}

			if data["port"].(float64) != tt.wantPort {
				t.Errorf("Wrong port for instance, given: %f", data["port"].(float64))
			}
		})
	}

}
