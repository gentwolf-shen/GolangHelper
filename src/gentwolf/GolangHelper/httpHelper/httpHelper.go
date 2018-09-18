package httpHelper

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func UrlValueToMap(urlValues url.Values) map[string]string {
	params := make(map[string]string, len(urlValues))
	for k, v := range urlValues {
		params[k] = v[0]
	}
	return params
}

func MapToUrlValue(data map[string]string) url.Values {
	params := url.Values{}

	for k, v := range data {
		params.Add(k, v)
	}
	return params
}

func Get(hostUrl string, params url.Values, headers map[string]string) ([]byte, error) {
	return httpRequest("GET", hostUrl, params, headers, false, nil)
}

func Delete(hostUrl string, params url.Values, headers map[string]string) ([]byte, error) {
	return httpRequest("DELETE", hostUrl, params, headers, false, nil)
}

func Put(hostUrl string, params url.Values, headers map[string]string) ([]byte, error) {
	return httpRequest("PUT", hostUrl, params, headers, false, nil)
}

func Post(hostUrl string, params url.Values, headers map[string]string) ([]byte, error) {
	return httpRequest("POST", hostUrl, params, headers, false, nil)
}

func PutToBody(hostUrl string, body []byte, headers map[string]string) ([]byte, error) {
	return httpRequest("PUT", hostUrl, nil, headers, true, body)
}

func PostToBody(hostUrl string, body []byte, headers map[string]string) ([]byte, error) {
	return httpRequest("POST", hostUrl, nil, headers, true, body)
}

func DeleteToBody(hostUrl string, body []byte, headers map[string]string) ([]byte, error) {
	return httpRequest("DELETEBODY", hostUrl, nil, headers, true, body)
}

func httpRequest(method string, hostUrl string, params url.Values, headers map[string]string, isPostToBody bool, body []byte) ([]byte, error) {
	var request *http.Request
	var err error
	if method == "GET" || method == "DELETE" {
		if !strings.Contains(hostUrl, "?") {
			hostUrl += "?"
		}

		hostUrl += params.Encode()
		request, err = http.NewRequest(method, hostUrl, nil)
	} else {
		if isPostToBody {
			if method == "DELETEBODY" {
				method = "DELETE"
			}

			request, err = http.NewRequest(method, hostUrl, bytes.NewReader(body))
			if headers == nil || len(headers) == 0 {
				request.Header.Add("Content-Type", "application/json;charset=utf-8")
			}
		} else {
			request, err = http.NewRequest(method, hostUrl, strings.NewReader(params.Encode()))
			if err == nil {
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			}
		}
	}

	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var responseBody []byte
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(response.Body)
		responseBody, _ = ioutil.ReadAll(reader)
		reader.Close()
	default:
		responseBody, _ = ioutil.ReadAll(response.Body)
	}

	if response.StatusCode == 200 {
		return responseBody, nil
	}

	return responseBody, errors.New(strconv.Itoa(response.StatusCode))
}

func PostWithFile(url string, param url.Values, files map[string]string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range param {
		for _, item := range v {
			_ = writer.WriteField(k, item)
		}
	}

	for k, v := range files {
		file, err := os.Open(v)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile(k, v)
		if err == nil {
			_, err = io.Copy(part, file)
		}

		file.Close()
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if nil != err {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func Download(filename string, url string) error {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}
