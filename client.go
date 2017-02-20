package congress

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const API_URL = "https://api.propublica.org/congress/v1"

type Client struct {
	httpClient *http.Client
	apiKey     string
}

type errorBody struct {
	Message string `json:"message"`
}

func NewClient(apiKey string) *Client {
	client := &Client{
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
	}

	return client
}

func (c *Client) Get(path string, queryParams map[string][]string, result interface{}) error {
	uri := fmt.Sprintf("%s/%s.json", API_URL, path)
	queryString := url.Values{}
	for key, values := range queryParams {
		for _, value := range values {
			queryString.Add(key, value)
		}
	}

	if len(queryString) > 0 {
		uri += "?" + queryString.Encode()
	}

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating request object: %s", err.Error()))
	}

	request.Header.Set("X-API-Key", c.apiKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("Error during request execution: %s", err.Error()))
	}

	bodyDecoder := json.NewDecoder(response.Body)
	if (response.StatusCode / 100) != 2 {
		body := errorBody{}
		err = bodyDecoder.Decode(&body)
		if err != nil {
			return errors.New(fmt.Sprintf("Error decoding response body: %s", err.Error()))
		}

		return errors.New(fmt.Sprintf("HTTP Error %s: %s", response.Status, body.Message))
	}

	err = bodyDecoder.Decode(result)
	if err != nil {
		return errors.New(fmt.Sprintf("Error decoding response body: %s", err.Error()))
	}

	return nil
}
