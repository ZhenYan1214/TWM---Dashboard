package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	url := "https://rest-mainnet.onflow.org/v1/scripts?blocks_height=final"
	method := "POST"

	var payload string
	payload = `{
		"script": "aW1wb3J0IEZsb3dJRFRhYmxlU3Rha2luZyBmcm9tIDB4ODYyNGI1MmY5ZGRjZDA0YQ0KDQovLyBUaGlzIHNjcmlwdCByZXR1cm5zIGFsbCB0aGUgaW5mbyBhc3NvY2lhdGVkIHdpdGggYSBkZWxlZ2F0b3INCg0KYWNjZXNzKGFsbCkgZnVuIG1haW4obm9kZUlEOiBTdHJpbmcsIGRlbGVnYXRvcklEOiBVSW50MzIpOiBGbG93SURUYWJsZVN0YWtpbmcuRGVsZWdhdG9ySW5mbyB7DQogICAgcmV0dXJuIEZsb3dJRFRhYmxlU3Rha2luZy5EZWxlZ2F0b3JJbmZvKG5vZGVJRDogbm9kZUlELCBkZWxlZ2F0b3JJRDogZGVsZWdhdG9ySUQpDQp9DQo=",
		"arguments": [
		"eyJ0eXBlIjoiU3RyaW5nIiwidmFsdWUiOiJlOGY0YmQ2NDlkMDhlY2I1YWZiNzAyM2EwYzVlOGJiMTBjZTU2NjU5Mzk5NjY1ZGE4Y2M5ZDUxN2U3OTgyZjkyIn0=",
		"eyJ0eXBlIjoiVUludDMyIiwidmFsdWUiOiIyIn0="
		]
	}`

	//_ = payload1
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

	//fmt.Println("Response Status:", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	response_body := string(body)
	//fmt.Println("Response Body:", response_body)
	fmt.Println(response_body)

}
