package main

import (
	"FSDS/fsds"
	"fmt"
	"os"
)

func main() {
	for _, dir := range os.Args[1:] {
		fsds.GetDirSizeMap(dir)
		fmt.Println()
	}
}
