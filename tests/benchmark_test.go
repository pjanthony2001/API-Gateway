package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func BenchmarkIntegration(b *testing.B) {

	data, err := json.Marshal(
		map[string]interface{}{
			"Message": "Message",
			"Flag":    340,
		})

	if err != nil {
		b.Fatalf("Error from creating json data %s", err.Error())
	}

	url := URL(0, 0, "")

	reqBody := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodGet, url, reqBody)
	if err != nil {
		b.Fatalf("Error from creating request %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	cli := http.DefaultClient

	resp, err := cli.Do(req)
	if err != nil {
		b.Fatalf("Error from making request %s", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		b.Fatalf("Error from reading response %s", err.Error())
	}

	var sb map[string]interface{}
	json.Unmarshal(body, &sb)

	responseMessage, _ := responseMessage(1, 1, "Message", "token")

	if responseMessage != sb["Message"] {
		b.Fatalf("Error with response: got %s, want %s", sb["Message"], responseMessage)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli.Do(req)
	}

}
