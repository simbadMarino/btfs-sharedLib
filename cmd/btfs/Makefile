# go/Makefile
#ANDROID_OUT=./bin/jniLibs
IOS_OUT=./bin/ios
#ANDROID_SDK=$(HOME)/Library/Android/sdk
#NDK_BIN=$(ANDROID_SDK)/ndk/21.4.7075529/toolchains/llvm/prebuilt/darwin-x86_64/bin
WASM_OUT=./bin/wasmBinary
WINDOWS_OUT=./bin/winBin
LINUX_OUT=./bin/linuxBin

# ios-arm64:
# 	CGO_ENABLED=1 \
# 	GOOS=ios \
# 	GOARCH=arm64 \
# 	SDK=iphoneos \
# 	SDK_PATH=`xcrun --sdk iphoneos --show-sdk-path` \
# 	CARCH="arm64" \
# 	CC=$(PWD)/clangwrap.sh \
# 	CGO_CFLAGS="-fembed-bitcode" \
# 	go build -v -buildmode=c-archive -ldflags="-s -w" -gcflags=all="-l -B" -tags ios -o $(IOS_OUT)/arm64.a .
#
#
# ios-x86_64:
# 	CGO_ENABLED=1 \
# 	GOOS=darwin \
# 	GOARCH=amd64 \
# 	CARCH="x86_64" \
# 	SDK=iphonesimulator \
# 	CGO_CFLAGS="-Wno-undef-prefix" \
# 	CC=$(PWD)/clangwrap.sh \
# 	go build -v -buildmode=c-archive -ldflags="-s -w" -gcflags=all="-l -B" -tags ios -o $(IOS_OUT)/x86_64.a .
#
#
# ios: ios-arm64 ios-x86_64
# 	lipo $(IOS_OUT)/x86_64.a $(IOS_OUT)/arm64.a -create -output $(IOS_OUT)/btfs.a
# 	cp $(IOS_OUT)/arm64.h $(IOS_OUT)/btfs.h
#
#
# android-armv7a:
# 	CGO_ENABLED=1 \
# 	GOOS=android \
# 	GOARCH=arm \
# 	GOARM=6 \
# 	CGO_CFLAGS="-fembed-bitcode" \
# 	CC=$(NDK_BIN)/armv7a-linux-androideabi21-clang \
# 	CXX=$(NDK_BIN)/armv7a-linux-androideabi21-clang++ \
# 	go build main_lib.go -v -buildmode=c-shared -trimpath -ldflags="-s -w" -gcflags=all="-l -B"  -o $(ANDROID_OUT)/armeabi-v7a/libbtfs.so .
#
#
# android-arm64:
# 	CGO_ENABLED=1 \
# 	GOOS=android \
# 	GOARCH=arm64 \
# 	CGO_CFLAGS="-fembed-bitcode" \
# 	CC=$(NDK_BIN)/aarch64-linux-android21-clang \
# 	CXX=$(NDK_BIN)/aarch64-linux-android21-clang++ \
# 	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(ANDROID_OUT)/arm64-v8a/libbtfs.so .
#
#
#
# android-x86_64:
# 	CGO_ENABLED=1 \
# 	GOOS=android \
# 	GOARCH=amd64 \
# 	CC=$(NDK_BIN)/x86_64-linux-android21-clang \
# 	CXX=$(NDK_BIN)/x86_64-linux-android21-clang++ \
# 	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w" -gcflags=all="-l -B"  -o $(ANDROID_OUT)/x86_64/libbtfs.so .
#
# android: android-armv7a android-arm64 android-x86_64
linux-amd64:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	CC="x86_64-linux-musl-gcc" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(LINUX_OUT)/linux_amd64 .

linux-i386:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=386 \
	CC="i486-linux-musl-gcc" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(LINUX_OUT)/linux_386 .

# linux-arm64:
# 	CGO_ENABLED=1 \
# 	GOOS=linux \
# 	GOARCH=arm64 \
# 	CC="aarch64-linux-musl-gcc" \
# 	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(LINUX_OUT)/linux_arm64 .

windows-amd64:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC="x86_64-w64-mingw32-gcc" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(WINDOWS_OUT)/btfs-amd64.exe .

windows-i386:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=386 \
	CC="i686-w64-mingw32-gcc" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(WINDOWS_OUT)/btfs-386.exe .

wasm:
	CGO_ENABLED=1 \
	GOOS=js \
	GOARCH=wasm \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -tags="nofuse purego" -o $(WASM_OUT)/btfs.wasm .


darwin_m1_macos:
	CGO_ENABLED=1 \
	GOOS=ios \
	GOARCH=arm64 \
	CARCH="arm64" \
	SDK=macos \
	CGO_CFLAGS="-fembed-bitcode" \
	go build -v -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -tags ios -o $(IOS_OUT)/btfs_macOS .

all:linux-amd64 linux-i386 windows-amd64 windows-i386 darwin_m1_macos
