package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func listRoles(claimMappingURL string) ([]interface{}, error) {
	var resp *http.Response
	var responseBody []byte
	method := "GET"
	var emptyResponseBody []interface{}
	URL := claimMappingURL + "/list/roles"

	request, err := http.NewRequest(method, URL, strings.NewReader(""))

	resp, err = http.DefaultClient.Do(request)
	if err == nil {
		responseBody, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode <= 300 {
			var f interface{}
			json.Unmarshal(responseBody, &f)
			switch f.(type) {
			case []interface{}:
				return f.([]interface{}), nil
			}

			m := f.([]interface{})

			return m, nil
		} else {
			err = fmt.Errorf("invalid Status code (%v)", resp.StatusCode)
			return emptyResponseBody, err
		}
	} else {
		return emptyResponseBody, err
	}
}
