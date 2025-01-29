//go:build (darwin || linux) && !arm
//go:build !windows
//go:build cgo
//go:build !no_protobuf
//go:build go1.18
//go:build darwin || linux
// +build darwin linux

//go:generate mockgen -source=myfile.go
//go:generate protoc --go_out=. myproto.proto
//go:generate stringer -type=MyEnumType
//go:generate command
//go:linkname
package main

import "fmt"

/*
#include <stdint.h>
uint8_t shift(uint8_t x, int y) {
  return x << y;
}
*/
import "C"

func main() {
	in, n := uint8(0x10), 2
	out := C.shift(C.uint8_t(in), C.int(2))
	fmt.Printf("Shifted %#x by %d, got %#x\n", in, n, out)
}
