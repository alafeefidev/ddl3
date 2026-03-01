package httputils

import (
	"io"
	"net/http"

	"github.com/alafeefidev/ddl3/internal/utils"
)

type HttpClient struct {
	BaseUrl    string
	Headers    map[string]string
	HttpClient *http.Client
}

func (h *HttpClient) ReqContentGet(path string) ([]byte, error) {
	//TODO add context timeout and such from config struct

	if _, err := utils.IsCorrectUrl(h.BaseUrl); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, h.BaseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range h.Headers {
		req.Header.Set(k, v)
	}

	resp, err := h.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func NewHttpClient(baseUrl string) *HttpClient {
	return &HttpClient{
		BaseUrl:    baseUrl,
		Headers:    make(map[string]string),
		HttpClient: &http.Client{},
	}
}
