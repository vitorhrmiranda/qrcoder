package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vitorhrmiranda/qrcode/qrcoder"
)

// Mimetypes enableds
const (
	MIME_SVG string = "image/svg+xml"
	MIME_PNG string = "image/png"
	MIME_PDF string = "application/pdf"
)

// Allowed extensions
const (
	SVG string = "svg"
	PNG string = "png"
	PDF string = "pdf"
)

// CreateQRCode is a factory to gerenate QRCode
func CreateQRCode(content string, extension string, size int) (qrcode []byte, mimetype string, err error) {
	var qrEncoder qrcoder.QRCoder

	switch mimetype {
	case SVG:
		qrEncoder = qrcoder.NewCoderSVG(content)
		qrcode, err = qrEncoder.Encode()
		mimetype = MIME_SVG

	case PNG:
		qrEncoder = qrcoder.NewCoderPNG(content, size)
		qrcode, err = qrEncoder.Encode()
		mimetype = MIME_PNG
	}

	return qrcode, mimetype, err
}

func Handler(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var statusCode int = 200
	var body string = request.QueryStringParameters["body"]
	var cacheFrom = time.Now().Format(http.TimeFormat)
	var cacheUntil = time.Now().AddDate(1, 0, 0).Format(http.TimeFormat)

	response = events.APIGatewayProxyResponse{}

	if len(body) == 0 {
		return response, fmt.Errorf("the 'body' query string parameter is required")
	}

	qrcode, mimetype, err := CreateQRCode(body, PNG, 500)
	if err != nil {
		return response, fmt.Errorf("create qrcode: %w", err)
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":  mimetype,
			"Last-Modified": cacheFrom,
			"Expires":       cacheUntil,
		},
		Body:            string(qrcode),
		StatusCode:      statusCode,
		IsBase64Encoded: true,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
