package main

import "gowatts/server"

func main() {
	site := server.New()
	site.Start()
}
