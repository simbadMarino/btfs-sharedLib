# go/Makefile
ANDROID_OUT=./jniLibs
ANDROID_SDK=$(HOME)/Library/Android/sdk
NDK_BIN=$(ANDROID_SDK)/ndk/21.1.6352462/toolchains/llvm/prebuilt/darwin-x86_64/bin
WASM_OUT=./wasmBinary
WINDOWS_OUT=./winBin

ios-arm64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=arm64 \
	SDK=iphoneos \
	SDK_PATH=`xcrun --sdk iphoneos --show-sdk-path` \
	CARCH="arm64" \
	CC=$(PWD)/clangwrap.sh \
	CGO_CFLAGS="-fembed-bitcode" \
	go build -v -buildmode=c-archive -ldflags="-s -w" -tags ios -o $(IOS_OUT)/arm64.a .

ios-x86_64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	CARCH="x86_64" \
	SDK=iphonesimulator \
	CGO_CFLAGS="-Wno-undef-prefix" \
	CC=$(PWD)/clangwrap.sh \
	go build -v -buildmode=c-archive -ldflags="-s -w" -tags ios -o $(IOS_OUT)/x86_64.a .


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
	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w"   -o $(ANDROID_OUT)/armeabi-v7a/libbtfs.so .


android-arm64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm64 \
	CGO_CFLAGS="-fembed-bitcode" \
	CC=$(NDK_BIN)/aarch64-linux-android21-clang \
	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w"  -o $(ANDROID_OUT)/arm64-v8a/libbtfs.so .


android-x86:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=386 \
	CC=$(NDK_BIN)/i686-linux-android21-clang \
	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w"  -o $(ANDROID_OUT)/x86/libbtfs.so .

android-x86_64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=amd64 \
	CC=$(NDK_BIN)/x86_64-linux-android21-clang \
	go build -v -buildmode=c-shared -trimpath -ldflags="-s -w"  -o $(ANDROID_OUT)/x86_64/libbtfs.so .

android: android-armv7a android-arm64 android-x86 android-x86_64

windows:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC="x86_64-w64-mingw32-gcc" \
	go build -buildmode=c-shared -o $(WINDOWS_OUT)/btfs-amd64.exe .

windows-amd64:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=386 \
	go build -buildmode=c-shared -o $(WINDOWS_OUT)/btfs-386.exe .

wasm:
	CGO_ENABLED=1 \
	GOOS=js \
	GOARCH=wasm \
	go build -o $(WASM_OUT)/btfs.wasm .