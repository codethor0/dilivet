package main

import (
	"flag"
	"fmt"
)

var version = "dev"

func main() {
	v := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *v {
		fmt.Println(version)
		return
	}
	fmt.Println("DiliVet - ML-DSA vetting tool. Use -version to print version.")
}
