package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestNacosServiceDiscovery(t *testing.T) {

	resp, err := http.Get("http://127.0.0.1:8848/nacos/v1/ns/instance/list?serviceName=ExampleService")

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

	if data["ip"].(string) != "127.0.0.1" {
		t.Errorf("Wrong IP addr for instance, given: %s", data["ip"].(string))
	}

	if data["port"].(float64) != float64(8888) {
		t.Errorf("Wrong port for instance, given: %f", data["port"].(float64))
	}

}
