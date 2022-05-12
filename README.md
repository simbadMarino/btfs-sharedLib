# btfs-sharedLib
BTFS shared libraries implementation for iOS and Android
This code is a modified fork from the official bittorrent/go-btfs code. Some dependencies were adapted/removed so it can run on iOS. Android support is pending to be tested.

Steps to build a shared library for iOS:

1. git clone this repository
2. cd into cmd/btfs
3. make ios


Steps to add shared library into iOS for testing:

1. Open btfs-ios-app in Xcode
2. Copy btfs.a and bfs.h into btfs/btfs-test using Xcode drag and drop function.
3. Run app
4. write "init" in text input field and push the "Do cli call" button, this will initialize your btfs node locally on your iphone
5. write btfs --chain-id 199 in the text input field and then push the "Do cli call" button, this will start running your btfs daemon (make sure you fill up your BTTC address so it can start properly
