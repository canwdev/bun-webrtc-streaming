@echo off
echo Copying frontend files...
xcopy /Y /Q "..\public\*" "public\" >nul
echo Building...
set CGO_ENABLED=0
go build -trimpath -ldflags="-s -w" -o bin\webrtc-server.exe .
echo Built: bin\webrtc-server.exe
