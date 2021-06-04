setup:
	mkdir -p bin
	go build -o bin/qr
	sls deploy
