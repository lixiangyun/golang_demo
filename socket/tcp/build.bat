set GOOS=windows
go build .
set GOOS=linux
go build .
set GOARCH=arm64
go build -o tcp_arm64

pause