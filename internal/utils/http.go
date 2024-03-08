package utils

import (
	"io/ioutil"
	"net/http"
)

func retrieveRequestBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// TODO MakePostRequest

// MakeGetRequest handles coherent get requests yay us
func MakeGetRequest(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return retrieveRequestBody(res)
}
