set GOOS=windows
go build .
set GOOS=linux
go build .

pause