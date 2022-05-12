# go/Makefile

ios-arm64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=arm64 \
	SDK=iphoneos \
	SDKROOT=/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk \
	CC=$(PWD)/clangwrap.sh \
	CGO_CFLAGS="-fembed-bitcode -Wno-undef-prefix" \
	go build -buildmode=c-archive -tags ios -o $(IOS_OUT)/arm64.a .

ios-x86_64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	SDK=iphonesimulator \
	CGO_CFLAGS="-Wno-undef-prefix" \
	CC=$(PWD)/clangwrap.sh \
	go build -buildmode=c-archive -tags ios -o $(IOS_OUT)/x86_64.a .

ios: ios-arm64 ios-x86_64
	lipo $(IOS_OUT)/x86_64.a $(IOS_OUT)/arm64.a -create -output $(IOS_OUT)/btfs.a
	cp $(IOS_OUT)/arm64.h $(IOS_OUT)/btfs.h