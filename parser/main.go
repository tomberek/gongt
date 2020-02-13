package main

import (
	"fmt"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"github.com/tomberek/gongt/parser/grp"
)

func main() {
	fmt.Println("vim-go")

	fd, _ := os.Open(os.Args[1])
	var ngt *grp.Ngt
	stream := kaitai.NewStream(fd)
	ngt.Read(stream, nil, ngt)
	fmt.Printf("%+v\n", ngt)
}
