package pica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jeremaihloo/funny/langs"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func CreateHttpRequest(req *ApiRequest, runner *ApiRunner) (httpReq *http.Request, err error) {
	var bodyParams map[string]langs.Value
	if req.Method != "GET" {
		bodyParams = runner.vm.Lookup(strings.ToLower(req.Method)).(map[string]langs.Value)
	}
	switch req.Headers["Content-Type"][0] {
	case "application/x-www-form-urlencoded":
		httpReq, err = createFormUrlEncodedRequest(req, runner, bodyParams)
		if err != nil {
			return nil, err
		}
		break
	case "multipart/form-data":
		httpReq, err = createFormDataRequest(req, runner, bodyParams)
		if err != nil {
			return nil, err
		}
		break
	case "application/json":
		httpReq, err = createJsonRequest(req, runner, bodyParams)
		if err != nil {
			return nil, err
		}
		break
	default:
		if req.Method == "GET" {
			targetUrl, err := getTargetUrl(req, runner)
			if err != nil {
				return nil, err
			}
			httpReq, err = http.NewRequest(req.Method, targetUrl, nil)
			if err != nil {
				return nil, err
			}
		}
		return nil, errors.New("unknow http")

	}
	httpReq.Header = req.Headers
	return
}

func getValue(val langs.Value) string {
	switch val := val.(type) {
	case int:
		return string(val)
	case string:
		return val
	default:
		panic(fmt.Errorf("unsupport type [%s], only support [int][string]", langs.Typing(val)))
	}
}

func getTargetUrl(req *ApiRequest, runner *ApiRunner) (string, error) {
	baseUrl := runner.vm.Lookup("baseUrl").(string)
	targetUrl, query, err := CompileUrl(baseUrl+req.Url, runner.vm)
	req.Query = query
	return targetUrl, err
}

func createFormUrlEncodedRequest(req *ApiRequest, runner *ApiRunner, bodyParams map[string]langs.Value) (*http.Request, error) {
	v := url.Values{}
	for key, val := range bodyParams {
		v.Set(key, getValue(val))
	}
	u := ioutil.NopCloser(strings.NewReader(v.Encode()))
	targetUrl, err := getTargetUrl(req, runner)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(req.Method, targetUrl, u)
}

func createFormDataRequest(req *ApiRequest, runner *ApiRunner, bodyParams map[string]langs.Value) (*http.Request, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	for key, val := range bodyParams {
		v := getValue(val)
		if strings.HasPrefix(v, "@") {
			fullFileName := v[1:]
			_, filename := filepath.Split(fullFileName)
			formFile, err := writer.CreateFormFile(key, filename)
			if err != nil {
				return nil, errors.New("Create form file failed: %s\n")
			}

			srcFile, err := os.Open(fullFileName)
			if err != nil {
				return nil, errors.New("%Open source file failed: s\n")
			}
			defer srcFile.Close()
			_, err = io.Copy(formFile, srcFile)
			if err != nil {
				return nil, errors.New("Write to form file falied: %s\n")
			}
		} else {
			writer.WriteField(key, v)
		}
	}
	writer.Close()
	targetUrl, err := getTargetUrl(req, runner)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(req.Method, targetUrl, buf)
}

func createJsonRequest(req *ApiRequest, runner *ApiRunner, bodyParams map[string]langs.Value) (*http.Request, error) {
	targetUrl, err := getTargetUrl(req, runner)
	if err != nil {
		return nil, err
	}
	jsonContent, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(req.Method, targetUrl, bytes.NewBuffer(jsonContent))
}