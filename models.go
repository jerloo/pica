package pica

import (
	"fmt"
	"github.com/jeremaihloo/funny/langs"
	"github.com/mitchellh/go-homedir"
	"net/http"
	"os"
	"path/filepath"
)

type ApiRequest struct {
	Headers     http.Header
	Method      string
	Url         string
	Query       Query
	Name        string
	Description string
	Body        []byte
	lines       langs.Block
}

type ApiResponse struct {
	Headers http.Header
	Body    []byte
	Status  int
	lines   langs.Block

	saveLines langs.Block
}

type ApiItem struct {
	Request  *ApiRequest
	Response *ApiResponse
}

type PicaContext struct {
	Name        string
	Description string
	Author      string
	Version     string

	PicaVersion      string
	MaxArrayShowRows int

	ApiItems     []*ApiItem
	Headers      map[string]langs.Value
	VersionNotes *VersionNote
}

var DefaultHeaders = map[string]langs.Value{
	"Accept":          "* /*",
	"Accept-Language": "en-US,en;q=0.8",
	"Cache-Control":   "max-age=0",
	"User-Agent":      fmt.Sprintf("Pica Api Test Client/%s https://github.com/jeremaihloo/pica", Version),
	"Connection":      "keep-alive",
	"Referer":         "http://www.baidu.com/",
	"Content-Type":    "application/json",
}

var DefaultInitScope = map[string]langs.Value{
	"headers": DefaultHeaders,
}

var PROFILE_HOME = ""

func init() {
	PROFILE_HOME, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	PROFILE_HOME = filepath.Join(PROFILE_HOME, ".pica")
	_, err = os.Stat(PROFILE_HOME)
	if err != nil {
		err = os.Mkdir(PROFILE_HOME, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}