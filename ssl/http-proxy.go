package ssl
import (
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"crypto/tls"
	"runtime"
)

func NewSSL(port string)  {
	ssl_proxy(port)
}

func ssl_proxy(port string) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cert, err := tls.LoadX509KeyPair("keys/cert.pem", "keys/key.pem")
	if err != nil {
		log.Fatalf("error in tls.LoadX509KeyPair: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	lis, err := tls.Listen("tcp", port, &config)
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleRequest(conn)
	}
}
func handleRequest(conn net.Conn) {
	if conn == nil {
		return
	}
	defer conn.Close()
	var buf [1024]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		log.Println(err)
		return
	}
	method := (strings.Split(string(buf[:])," "))[0]
	host := (strings.Split(string(buf[:])," "))[1]
	uRl, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}
	var address string
	// www.baidu.com is https://www.baidu.com, if use www.baidu.com not work well
	if uRl.Opaque == "443" {
		address = uRl.Scheme + ":443"
	} else {
		if strings.Index(uRl.Host, ":") == -1 {
			address = uRl.Host + ":80"
		} else {
			address = uRl.Host
		}
	}
	// request
	fmt.Println("访问地址: ",address)
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Fprint(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(buf[:n])
	}
	// convert data
	go io.Copy(server, conn)
	io.Copy(conn, server)
}
