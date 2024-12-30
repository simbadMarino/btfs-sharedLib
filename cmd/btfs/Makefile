# go/Makefile
IOS_OUT=./bin/ios
WASM_OUT=./bin/wasmBinary
WINDOWS_OUT=./bin/winBin
LINUX_OUT=./bin/linuxBin

linux-amd64:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	CC="x86_64-linux-musl-gcc" \
	CXX="x86_64-linux-musl-g++" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(LINUX_OUT)/linux_amd64 .

linux-i386:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=386 \
	CC="i486-linux-musl-gcc" \
	CXX="i486-linux-musl-g++" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(LINUX_OUT)/linux_386 .


windows-amd64:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC="x86_64-w64-mingw32-gcc" \
	CXX="x86_64-w64-mingw32-g++" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(WINDOWS_OUT)/btfs-amd64.exe .

windows-i386:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=386 \
	CC="i686-w64-mingw32-gcc" \
	CXX="i686-w64-mingw32-g++" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(WINDOWS_OUT)/btfs-386.exe .

wasm:
	CGO_ENABLED=1 \
	GOOS=js \
	GOARCH=wasm \
	go build -o $(WASM_OUT)/btfs.wasm .


darwin-arm64-macos:
	CGO_ENABLED=1 \
	GOOS=ios \
	GOARCH=arm64 \
	CARCH="arm64" \
	SDK=macos \
	go build -x -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -tags ios -o $(IOS_OUT)/btfs_macOS .

all:linux-amd64 linux-i386 windows-amd64 windows-i386 darwin-arm64-macos
