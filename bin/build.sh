env GOOS=linux GOARCH=386 go build -o build/chappy-linux-386
env GOOS=linux GOARCH=amd64 go build -o build/chappy-linux-amd64
env GOOS=linux GOARCH=arm go build -o build/chappy-linux-arm
env GOOS=linux GOARCH=arm64 go build -o build/chappy-linux-arm64

env GOOS=netbsd GOARCH=386 go build -o build/chappy-netbsd-386
env GOOS=netbsd GOARCH=amd64 go build -o build/chappy-netbsd-amd64
env GOOS=netbsd GOARCH=arm go build -o build/chappy-netbsd-arm

env GOOS=darwin GOARCH=386 go build -o build/chappy-macos-386
env GOOS=darwin GOARCH=amd64 go build -o build/chappy-macos-amd64
env GOOS=darwin GOARCH=arm go build -o build/chappy-macos-arm
env GOOS=darwin GOARCH=arm64 go build -o build/chappy-macos-arm64

env GOOS=windows GOARCH=386 go build -o build/chappy-windows-386
mv build/chappy-windows-386 build/chappy-windows-386.exe
env GOOS=windows GOARCH=amd64 go build -o build/chappy-windows-amd64
mv build/chappy-windows-amd64 build/chappy-windows-amd64.exe
