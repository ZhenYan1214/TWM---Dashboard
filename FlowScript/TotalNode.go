package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	url := "https://rest-mainnet.onflow.org/v1/scripts?blocks_height=final"
	method := "POST"

	payload := `{
		"script": "aW1wb3J0IEZsb3dJRFRhYmxlU3Rha2luZyBmcm9tIDB4ODYyNGI1MmY5ZGRjZDA0YQ0KDQovLyBUaGlzIHNjcmlwdCBnZXRzIGFsbCB0aGUgaW5mbyBhYm91dCBhIG5vZGUgYW5kIHJldHVybnMgaXQNCmFjY2VzcyhhbGwpIGZ1biBtYWluKG5vZGVJRDogU3RyaW5nKTogRmxvd0lEVGFibGVTdGFraW5nLk5vZGVJbmZvIHsNCiAgICByZXR1cm4gRmxvd0lEVGFibGVTdGFraW5nLk5vZGVJbmZvKG5vZGVJRDogbm9kZUlEKQ0KfQ==",
                "arguments": [
                "eyJ0eXBlIjoiU3RyaW5nIiwidmFsdWUiOiJlOGY0YmQ2NDlkMDhlY2I1YWZiNzAyM2EwYzVlOGJiMTBjZTU2NjU5Mzk5NjY1ZGE4Y2M5ZDUxN2U3OTgyZjkyIn0="
                ]
	}`

	body1 := strings.NewReader(payload)
	req, err := http.NewRequest(method, url, body1)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	cleaned := strings.Trim(string(body), "\"")
	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}

	// 用 RawMessage 動態解析
	var outer struct {
		Value struct {
			Fields []struct {
				Name  string          `json:"name"`
				Value json.RawMessage `json:"value"`
			} `json:"fields"`
		} `json:"value"`
	}
	if err := json.Unmarshal(decoded, &outer); err != nil {
		fmt.Println("Error parsing decoded json:", err)
		return
	}

	var id, tokensRewarded string
	for _, field := range outer.Value.Fields {
		if field.Name == "id" {
			// 解析 value.value 為 string
			var v struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(field.Value, &v); err == nil {
				id = v.Value
			}
		}
		if field.Name == "tokensRewarded" {
			var v struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(field.Value, &v); err == nil {
				tokensRewarded = v.Value
			}
		}
	}
	fmt.Printf("id: %s, tokensRewarded: %s\n", id, tokensRewarded)
}
