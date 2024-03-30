# go/Makefile
ANDROID_OUT=./bin/jniLibs
IOS_OUT=./bin/ios
ANDROID_SDK=$(HOME)/Library/Android/sdk
NDK_BIN=$(ANDROID_SDK)/ndk/23.1.7779620/toolchains/llvm/prebuilt/darwin-x86_64/bin

# README
# CGO_CFLAGS=-v  //This flag can be used to obtain more building information on C and C++ related building processes.
# CGO_LDFLAGS=-v //This flag can be used to obtain more building information on C and C++ related building processes.
# -lresol flag is used starting from go 1.20 to properly link c-archive libraries.
# anacrolix/torrrent repo needs to be compiled using C++, we need to add CXX flag to avoid Apple Clang going for the default macos compiler and thus having a btfs.a binary file with macOS static links
# END OF README

ios-arm64:
	CGO_ENABLED=1 \
	GOOS=ios \
	GOARCH=arm64 \
	CC=$(PWD)/clangwrap.sh \
	CXX=$(PWD)/clangwrap.sh \
	CGO_CFLAGS="-fembed-bitcode" \
	CGO_LDFLAGS=-lresolv \
	go build -v -buildmode=c-archive -ldflags="-s -w" -gcflags=all="-l -B" -tags ios -o $(IOS_OUT)/btfs.a .


ios-x86_64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	CARCH="x86_64" \
	SDK=iphonesimulator \
	CGO_CFLAGS="-Wno-undef-prefix" \
	CC=$(PWD)/clangwrap.sh \
	CXX=$(PWD)/clangwrap.sh \
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
