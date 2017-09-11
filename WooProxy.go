package main

import (
	"WooProxy/ssl"
	//"WooProxy/sockets"
)

func main() {
	ssl.NewSSL(":8080")
	//sockets.NewSockets()
}
