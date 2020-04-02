package main

import (
	"fmt"
	"gowatts/server"
)

func main() {
	site := server.New()
	fmt.Println(site.Start())
}
