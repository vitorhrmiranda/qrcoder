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

	qrcode, mimetype, err := qrcoder.CreateQRCode(body, ext, size)
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
