package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vitorhrmiranda/qrcode/qrcoder"
)

func buildRequest() (request *events.ALBTargetGroupRequest) {
	request = new(events.ALBTargetGroupRequest)

	dat, err := ioutil.ReadFile("./request.json")
	if err != nil {
		return request
	}

	json.Unmarshal(dat, request)

	return request
}

func TestHandler(t *testing.T) {
	request := buildRequest()

	t.Run("when method isnt POST", func(t *testing.T) {
		r := *request

		r.HTTPMethod = http.MethodGet

		response, err := Handler(r)
		if response.StatusCode != http.StatusInternalServerError || err == nil {
			t.Errorf("Handler() \n response = %#v", response)
		}
	})

	t.Run("when has invalid body", func(t *testing.T) {
		r := *request

		r.Body = ``

		response, err := Handler(r)
		if response.StatusCode != http.StatusInternalServerError || err == nil {
			t.Errorf("Handler() \n response = %#v", response)
		}
	})

	t.Run("when content value is empty", func(t *testing.T) {
		r := *request

		r.Body = `{"content": ""}`

		response, err := Handler(r)
		if response.StatusCode != http.StatusInternalServerError || err == nil {
			t.Errorf("Handler() \n response = %#v", response)
		}
	})

	t.Run("when extension value is empty", func(t *testing.T) {
		r := *request

		r.Body = `{"content": "enjoei"}`

		response, err := Handler(r)
		if err != nil || response.Headers["Content-Type"] != qrcoder.MIME_SVG {
			t.Errorf("Handler() \n response = %#v", response)
		}
	})

	t.Run("when use invalid extension", func(t *testing.T) {
		r := *request

		r.Body = `{"content": "enjoei", "extension": "inv"}`

		response, err := Handler(r)
		if response.StatusCode != http.StatusInternalServerError || err == nil {
			t.Errorf("Handler() \n response = %#v", response)
		}
	})

	t.Run("when use base64 encode", func(t *testing.T) {
		r := *request

		r.Body = `{"content": "enjoei", "base64_encode": true}`

		response, err := Handler(r)
		if !response.IsBase64Encoded || err != nil {
			t.Errorf("Handler() \n response = %#v", response)
		}
	})

	t.Run("when describe extension in url path", func(t *testing.T) {
		r := *request

		r.Body = `{"content": "enjoei", "size": 500}`
		r.Path = "/generate.pdf"

		response, err := Handler(r)
		if err != nil || response.Headers["Content-Type"] != qrcoder.MIME_PDF {
			t.Errorf("Handler() \n response = %#v", response)
		}
	})
}
