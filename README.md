# go-btfs (Shared Library)

## Why do we need a Shared Library repo besides the official go-btfs repo?

1. Bittorrent does not officially maintains experimental mobile builds
2. BitTorrent does not officially maintains experimental wasm builds
   I've decided to go and fork go-btfs for experimental purposes and to release the necesary shared libraries for iOS and Android for the dCloud app

## Repository Organization & Build Instructions

Due to compatibility issues between desktop and mobile environments we have to keep both separated.

### Mobile Build

Location:

> Mobile root make folder:
>
> root/
>
> |---cmd
>
> |------btfs_lib

btfs_lib contains all original btfs folder files with slight differences:

`main.go` This includes the following changes:

> btfs main differences report
> Produced: 15/12/24 19:04:03
>
> Mode:  Differences
> Left file: main\_lib.go
> Right file: main.go
>
>
> | 4   | **import** **"C"**                                                          | +- |    |                      |
> | --- | --------------------------------------------------------------------------- | -- | -- | -------------------- |
> | 15  | **"strconv"**                                                               | +- |    |                      |
> | 16  | **"strings"**                                                               |    |    |                      |
> | 75  | **//** os.Exit(mainRet())                                                   | <> | 71 | os.Exit(mainRet())   |
> | 78  | **//export** **mainC**                                                      | <> |    |                      |
> | 79  | **func** **mainC(in** **\*C.char)** **\*C.char** **{**                      |    |    |                      |
> | 80  | **args** **:=** **strings.Split(C.GoString(in),** **"** **")**              |    |    |                      |
> | 81  | **args** **=** **append([]string{"btfs"},** **args...)**                    |    |    |                      |
> | 82  | **fmt.Println("args:",** **args)**                                          |    |    |                      |
> | 83  | **exitCode** **:=** **mainRet(args)**                                       |    |    |                      |
> | 84  | **return** **C.CString("exit** **code:"** **+** **strconv.Itoa(exitCode))** |    |    |                      |
> | 85  | **}**                                                                       |    |    |                      |
> | 86  | func mainRet(**args** **[]string**) int {                                   |    | 74 | func mainRet() int { |
> | 184 | **os.Args** **=** **args**                                                  | +- |    |                      |

`Makefile` This file specifies the different mobile builds requirements.

Summary:

* For cgo shared libraries it is important to explicitly export the main function as `mainC`, thus we need to `import C` at the beggining of our main file.

Stable builds:

* ios-arm64
* ios-x86_64
* ios (fat binary=arm64 + x86_64)
* android-armv7a
* android-arm64
* android-x86_64
* android (fat binary=armv7a + arm64 + x86_64)

### Desktop Build

Location:

> Deskptop experimental root make folder:
>
> root/
>
> |--cmd
>
> |------btfs

This folder contains the original cmd/btfs files, no additional changes are made besides a specific Makefile to compile for different OS and wasm.

`Makefile` This file specifies the different desktop builds requirements.

Stable builds:

* linux-amd64
* linux-i386
* windows-amd64
* windows-i386
* wasm
* darwin-arm64-macos

### How to Build

Pre setup:

If you have old btfs-sharedlib go packages versions make sure to clean your go cache and ensure dependencies are updated by:

* `go clean -modcache`  // Cleans the cached packages
* `go mod tidy`  //matches the go.mod file with dependencies required in the source files

Build your desired shared library or binary file by:

1. Clone the https://github.com/simbadMarino/btfs-sharedLib repository
2. `cd cmd/btfs/` for Desktop and wasm | `cd cmd/btfs_lib` for mobile
3. `make YOUR_BUILD` where "YOUR_BUILD" is any of the builds described in the above section

## go-btfs Reference

For further details about go-btfs please go to: https://github.com/bittorrent/go-btfs

## License 

[MIT](./LICENSE)
