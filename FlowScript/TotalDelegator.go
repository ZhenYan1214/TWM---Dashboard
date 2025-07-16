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

	ids := []int{2, 3, 4}
	for _, delegatorID := range ids {
		// 產生 delegatorID 的 base64
		arg := fmt.Sprintf(`{"type":"UInt32","value":"%d"}`, delegatorID)
		argBase64 := base64.StdEncoding.EncodeToString([]byte(arg))

		payload := fmt.Sprintf(`{
            "script": "aW1wb3J0IEZsb3dJRFRhYmxlU3Rha2luZyBmcm9tIDB4ODYyNGI1MmY5ZGRjZDA0YQ0KDQovLyBUaGlzIHNjcmlwdCByZXR1cm5zIGFsbCB0aGUgaW5mbyBhc3NvY2lhdGVkIHdpdGggYSBkZWxlZ2F0b3INCg0KYWNjZXNzKGFsbCkgZnVuIG1haW4obm9kZUlEOiBTdHJpbmcsIGRlbGVnYXRvcklEOiBVSW50MzIpOiBGbG93SURUYWJsZVN0YWtpbmcuRGVsZWdhdG9ySW5mbyB7DQogICAgcmV0dXJuIEZsb3dJRFRhYmxlU3Rha2luZy5EZWxlZ2F0b3JJbmZvKG5vZGVJRDogbm9kZUlELCBkZWxlZ2F0b3JJRDogZGVsZWdhdG9ySUQpDQp9DQo=",
            "arguments": [
                "eyJ0eXBlIjoiU3RyaW5nIiwidmFsdWUiOiJlOGY0YmQ2NDlkMDhlY2I1YWZiNzAyM2EwYzVlOGJiMTBjZTU2NjU5Mzk5NjY1ZGE4Y2M5ZDUxN2U3OTgyZjkyIn0=",
                "%s"
            ]
        }`, argBase64)

		body1 := strings.NewReader(payload)
		req, err := http.NewRequest(method, url, body1)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		// 先去掉最外層的雙引號
		cleaned := strings.Trim(string(body), "\"")
		decoded, err := base64.StdEncoding.DecodeString(cleaned)
		if err != nil {
			fmt.Println("Error decoding base64:", err)
			continue
		}

		// 解析巢狀 JSON
		var outer struct {
			Value struct {
				Fields []struct {
					Name  string `json:"name"`
					Value struct {
						Value string `json:"value"`
					} `json:"value"`
				} `json:"fields"`
			} `json:"value"`
		}
		if err := json.Unmarshal(decoded, &outer); err != nil {
			fmt.Println("Error parsing decoded json:", err)
			continue
		}

		var id, tokensStaked string
		for _, field := range outer.Value.Fields {
			if field.Name == "id" {
				id = field.Value.Value
			}
			if field.Name == "tokensStaked" {
				tokensStaked = field.Value.Value
			}
		}

		fmt.Printf("DelegatorID: %s, tokensStaked: %s\n", id, tokensStaked)
	}
}
