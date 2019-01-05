@echo off
mkdir out
echo Compiling for Linux ARMv7
set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -o out/linux-armv7 -ldflags="-s -w" .
echo Compiling for Windows x64
set GOOS=windows
set GOARCH=amd64
go build -o out/windows-x64.exe .

