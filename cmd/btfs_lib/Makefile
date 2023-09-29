# go/Makefile
ANDROID_OUT=./bin/jniLibs
IOS_OUT=./bin/ios
ANDROID_SDK=$(HOME)/Library/Android/sdk
NDK_BIN=$(ANDROID_SDK)/ndk/21.4.7075529/toolchains/llvm/prebuilt/darwin-x86_64/bin


ios-arm64:
	CGO_ENABLED=1 \
	GOOS=ios \
	GOARCH=arm64 \
	SDK=iphoneos \
	SDK_PATH=`xcrun --sdk iphoneos --show-sdk-path` \
	CARCH="arm64" \
	CC=$(PWD)/clangwrap.sh \
	CGO_CFLAGS="-fembed-bitcode" \
	go build -v -buildmode=c-archive -ldflags="-s -w" -gcflags=all="-l -B" -tags ios -o $(IOS_OUT)/arm64.a .


ios-x86_64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	CARCH="x86_64" \
	SDK=iphonesimulator \
	CGO_CFLAGS="-Wno-undef-prefix" \
	CC=$(PWD)/clangwrap.sh \
	go build -v -buildmode=c-archive -ldflags="-s -w" -gcflags=all="-l -B" -tags ios -o $(IOS_OUT)/x86_64.a .


ios: ios-arm64 ios-x86_64
	lipo $(IOS_OUT)/x86_64.a $(IOS_OUT)/arm64.a -create -output $(IOS_OUT)/btfs.a
	cp $(IOS_OUT)/arm64.h $(IOS_OUT)/btfs.h


android-armv7a:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm \
	GOARM=6 \
	CGO_CFLAGS="-fembed-bitcode" \
	CC=$(NDK_BIN)/armv7a-linux-androideabi21-clang \
	CXX=$(NDK_BIN)/armv7a-linux-androideabi21-clang++ \
	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w" -gcflags=all="-l -B"  -o $(ANDROID_OUT)/armeabi-v7a/libbtfs.so .


android-arm64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm64 \
	CGO_CFLAGS="-fembed-bitcode" \
	CC=$(NDK_BIN)/aarch64-linux-android21-clang \
	CXX=$(NDK_BIN)/aarch64-linux-android21-clang++ \
	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w" -gcflags=all="-l -B" -o $(ANDROID_OUT)/arm64-v8a/libbtfs.so .



android-x86_64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=amd64 \
	CC=$(NDK_BIN)/x86_64-linux-android21-clang \
	CXX=$(NDK_BIN)/x86_64-linux-android21-clang++ \
	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w" -gcflags=all="-l -B"  -o $(ANDROID_OUT)/x86_64/libbtfs.so .

android: android-armv7a android-arm64 android-x86_64

all: ios-arm64 android-armv7a android-arm64 android-x86_64
