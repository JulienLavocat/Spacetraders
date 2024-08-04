package main

import "os"

func main() {
	for _, system := range os.Args[1:] {
		println(system)
	}
}
