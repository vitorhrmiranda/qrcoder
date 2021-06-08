package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vitorhrmiranda/qrcode/qrcoder"
)

type body struct {
	Content      string `json:"content"`
	Extension    string `json:"extension"`
	Size         int    `json:"size,omitempty"`
	Base64Encode bool   `json:"base64_encode,omitempty"`
}

func extension(path string, b *body) {
	re := regexp.MustCompile(`(\.)(\w{3,4})($|\?)`)
	f := re.FindStringSubmatch(path)
	if len(f) > 0 {
		b.Extension = f[len(f)-2]
	}

	//default extension
	if len(b.Extension) == 0 {
		b.Extension = qrcoder.SVG
	}
}

func Handler(request events.ALBTargetGroupRequest) (response events.ALBTargetGroupResponse, err error) {
	response = events.ALBTargetGroupResponse{StatusCode: http.StatusInternalServerError}

	if request.HTTPMethod != http.MethodPost {
		return response, fmt.Errorf("method %s not allowed. Use POST", request.HTTPMethod)
	}

	var rb = new(body)
	if err = json.Unmarshal([]byte(request.Body), rb); err != nil {
		return response, fmt.Errorf("invalid parameters")
	}

	var cacheFrom = time.Now().Format(http.TimeFormat)
	var cacheUntil = time.Now().AddDate(1, 0, 0).Format(http.TimeFormat)

	if len(rb.Content) == 0 {
		return response, fmt.Errorf("the 'content' value parameter is required")
	}

	if rb.Size == 0 {
		rb.Size = 500
	}

	extension(request.Path, rb)

	qrcode, mimetype, err := qrcoder.CreateQRCode(rb.Content, rb.Extension, rb.Size)
	if err != nil {
		return response, fmt.Errorf("create qrcode: %w", err)
	}

	var responseBody string = string(qrcode)
	if rb.Base64Encode {
		responseBody = base64.StdEncoding.EncodeToString(qrcode)
	}

	return events.ALBTargetGroupResponse{
		Headers: map[string]string{
			"Content-Type":  mimetype,
			"Last-Modified": cacheFrom,
			"Expires":       cacheUntil,
		},
		Body:            responseBody,
		StatusCode:      http.StatusOK,
		IsBase64Encoded: rb.Base64Encode,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
