$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -ldflags="-s -w" -o tftp/tftp-$env:GOOS-$env:GOARCH .\tftp\main.go