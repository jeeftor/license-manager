//go:build (darwin || linux) && !arm
//go:build !windows
//go:build cgo
//go:build !no_protobuf
//go:build go1.18

//go:generate mockgen -source=myfile.go
//go:generate protoc --go_out=. myproto.proto
//go:generate stringer -type=MyEnumType

package main2

import "fmt"

func greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func main() {
	fmt.Println(greet("World"))
}
