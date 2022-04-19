build-assets:
	go-assets-builder static -o app/assets/assets.go -p assets

build-windows:
	GOOS=windows go build -o dist/dingdong-assistant.exe