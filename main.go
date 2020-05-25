package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
	"math/rand"
	"bufio"
    "fmt"
    "strings"
)

func get_list_proxy() []string{
	file, err := os.Open("listproxy.txt")
 
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
 
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var txtlines []string
 
    for scanner.Scan() {
        txtlines = append(txtlines, strings.TrimSuffix(scanner.Text(),"\n"))
    }
 
    file.Close()
 
    for _, eachline := range txtlines {
        fmt.Println(eachline)
    }
    return txtlines
}

func select_proxy() string{
	listproxies := get_list_proxy()
	rand.Seed(time.Now().Unix())
	return listproxies[rand.Intn(len(listproxies))]
}

func main() {
	// if len(os.Args) <= 2 {
	// 	log.Fatal("usage: portfw local:port remote:port")
	// }
	localAddrString := "0.0.0.0:8181"

	localAddr, err := net.ResolveTCPAddr("tcp", localAddrString)
	if localAddr == nil {
		log.Fatalf("net.ResolveTCPAddr failed: %s", err)
	}
	local, err := net.ListenTCP("tcp", localAddr)
	if local == nil {
		log.Fatalf("portfw: %s", err)
	}
	log.Printf("portfw listen on %s", localAddr)

	for {
		conn, err := local.Accept()
		if conn == nil {
			log.Printf("accept failed: %s", err)
			continue
		}
		remoteAddrString := select_proxy()
		go forward(conn, remoteAddrString)
	}
}

/// forward requests to other host.
func forward(local net.Conn, remoteAddr string) {

	remote, err := net.DialTimeout("tcp", remoteAddr, time.Duration(5*time.Second))
	if remote == nil {
		log.Printf("remote dial failed: %v\n", err)
		local.Close()
		return
	}
	go func() {
		defer local.Close()
		io.Copy(local, remote)
	}()
	go func() {
		defer remote.Close()
		io.Copy(remote, local)
	}()
	log.Printf("forward %s to %s", local.RemoteAddr(), remoteAddr)
}
