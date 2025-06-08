package fetchRabbit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Define the response struct (same as yours)
type Result struct {
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		CacheControl   string `json:"Cache-Control"`
		Connection     string `json:"Connection"`
		DNT            string `json:"DNT"`
		Host           string `json:"Host"`
		Pragma         string `json:"Pragma"`
		SecFetchDest   string `json:"Sec-Fetch-Dest"`
		SecFetchMode   string `json:"Sec-Fetch-Mode"`
		SecFetchSite   string `json:"Sec-Fetch-Site"`
		SecGPC         string `json:"Sec-GPC"`
		TE             string `json:"TE"`
		UserAgent      string `json:"User-Agent"`
		Referer        string `json:"Referer"`
	} `json:"headers"`
	Provider string `json:"provider"`
	Servers  []interface{} `json:"servers"`
	URL      []struct {
		Lang string `json:"lang"`
		Link string `json:"link"`
		Type string `json:"type"`
	} `json:"url"`
	Tracks []struct {
		Lang string `json:"lang"`
		URL  string `json:"url"`
	} `json:"tracks"`
}

// FetchRabbit fetches and parses the API response with custom headers
func FetchRabbit(id string, ss *string, ep *string) *Result {
	url := fmt.Sprintf("https://api.vidjoy.pro/rabbit/fetch/%s", id)
	if ss != nil {
		url = fmt.Sprintf("%s?ss=%s", url, *ss)
	}
	if ep != nil {
		url = fmt.Sprintf("%s&ep=%s", url, *ep)
	}

	// Create a custom request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request Creation Error:", err)
		return nil
	}

	// Set custom headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("x-api-key", "")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request Error:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Error:", err)
		return nil
	}

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("JSON Parse Error:", err)
		fmt.Println("Raw Body:", string(body))
		return nil
	}

	return &result
}
