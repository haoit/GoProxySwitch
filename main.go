package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
	"bufio"
    "strings"
	"math/rand"
)

var list_proxy_used []string;
var current_item int = 0;


func get_list_proxy(filename string) []string{
	file, err := os.Open(filename)
 
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
 
    // for _, eachline := range txtlines {
    //     fmt.Println(eachline)
    // }
    return txtlines
}

func select_untrust_proxy() string{
	var proxy_select string
	listproxies := get_list_proxy("listproxy.txt")

	//Get random value proxy
	// rand.Seed(time.Now().Unix())
	// return listproxies[rand.Intn(len(listproxies))]
	if(len(list_proxy_used) >= len(listproxies)){
		//Reset list used proxy
		// log.Printf("Reset list proxy used:!!!!!!!", list_proxy_used)
		list_proxy_used = list_proxy_used[:0]
		current_item = 0
	}
	proxy_select = listproxies[current_item]
	current_item = current_item +1
	list_proxy_used = append(list_proxy_used, proxy_select)
	return proxy_select
}

func select_trust_proxy() string{
	listproxies := get_list_proxy("trustproxy.txt")

	//Get random value proxy
	rand.Seed(time.Now().Unix())
	return listproxies[rand.Intn(len(listproxies))]
}

func main() {
	// if len(os.Args) <= 2 {
	// 	log.Fatal("usage: portfw local:port remote:port")
	// }
	localAddrString := "0.0.0.0:57812"

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
		remoteAddrString := select_untrust_proxy()
		go forward(conn, remoteAddrString)
	}
}

/// forward requests to other host.
func forward(local net.Conn, remoteAddr string) {

	remote, err := net.DialTimeout("tcp", remoteAddr, time.Duration(5*time.Second))
	if remote == nil {
		log.Printf("remote dial failed1: %v\n", err)
		remoteAddr = select_trust_proxy()
		remote, err = net.DialTimeout("tcp", remoteAddr, time.Duration(5*time.Second))
		log.Printf("Using trust proxy: %s\n", remoteAddr)
		if remote == nil{
			log.Printf("remote dial failed2: %v\n", err)
			local.Close()
			return
		}
		// local.Close()
		// return
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
