//go:generate sh -c "protoc -I . -I \"$(go list -f '{{ .Dir }}' -m github.com/bittorrent/protobuf)/protobuf\" --gogofaster_out=. swap.proto"

package pb
