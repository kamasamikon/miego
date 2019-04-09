package main

import (
	"fmt"
	_ "github.com/kamasamikon/miego/conf"
	_ "github.com/kamasamikon/miego/klog"
	_ "github.com/kamasamikon/miego/misc"
	_ "github.com/kamasamikon/miego/msa"
	_ "github.com/kamasamikon/miego/msb"
	_ "github.com/kamasamikon/miego/node"
)

func main() {
	fmt.Println("vim-go")
}
