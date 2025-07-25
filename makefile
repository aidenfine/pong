APP_EXEC=pong
MAIN_PATH=cmd/pong/main.go
build:
	GOARCH=amd64 GOOS=darwin go build -o ${APP_EXEC}-darwin ${MAIN_PATH}
	GOARCH=amd64 GOOS=linux go build -o ${APP_EXEC}-linux ${MAIN_PATH}
	GOARCH=amd64 GOOS=windows go build -o ${APP_EXEC}-windows ${MAIN_PATH}
run: build
	./${APP_EXEC}-darwin

clean:
	go clean
	rm -f ${APP_EXEC}-darwin
	rm -f ${APP_EXEC}-linux
	rm -f ${APP_EXEC}-windows