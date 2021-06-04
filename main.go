package main

import (
	"encoding/base64"
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

	switch extension {
	case SVG:
		qrEncoder = qrcoder.NewCoderSVG(content)
		qrcode, err = qrEncoder.Encode()
		mimetype = MIME_SVG

	case PNG:
		qrEncoder = qrcoder.NewCoderPNG(content, size)
		qrcode, err = qrEncoder.Encode()
		mimetype = MIME_PNG

	default:
		err = fmt.Errorf("extension not recognized")
	}

	return qrcode, mimetype, err
}

func Handler(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var size int

	var body string = request.QueryStringParameters["body"]
	var ext string = request.QueryStringParameters["ext"]

	var cacheFrom = time.Now().Format(http.TimeFormat)
	var cacheUntil = time.Now().AddDate(1, 0, 0).Format(http.TimeFormat)

	response = events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}

	if len(body) == 0 {
		return response, fmt.Errorf("the 'body' query string parameter is required")
	}

	if _, err := fmt.Scanf("%d", &size); err != nil {
		size = 500
	}

	qrcode, mimetype, err := CreateQRCode(body, ext, size)
	if err != nil {
		return response, fmt.Errorf("create qrcode: %w", err)
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":  mimetype,
			"Last-Modified": cacheFrom,
			"Expires":       cacheUntil,
		},
		Body:            base64.StdEncoding.EncodeToString(qrcode),
		StatusCode:      http.StatusOK,
		IsBase64Encoded: true,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
