cls
set GOPATH=N:\360тфел\workspace\golang\crash_analyze2\
cd bin
del main.exe
go build ../main/main.go
xcopy /y /f "main.exe"  "../bin-server/"
xcopy /y /f "main.exe"  "../bin-client/"
cd ..