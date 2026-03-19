package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var port = flag.Int("port", 8080, "端口")

func main() {
	flag.Parse()
	localIP := getLocalIP()
	
	fmt.Printf("\n🛠️  Bench - 开发工具箱\n")
	fmt.Printf("访问: http://%s:%d\n", localIP, *port)
	
	webDir := getWebDir()
	http.HandleFunc("/api/random", handleRandom)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleRandom(w http.ResponseWriter, r *http.Request) {
	min := 0
	max := 100
	count := 1
	
	if m := r.URL.Query().Get("min"); m != "" {
		fmt.Sscanf(m, "%d", &min)
	}
	if m := r.URL.Query().Get("max"); m != "" {
		fmt.Sscanf(m, "%d", &max)
	}
	if c := r.URL.Query().Get("count"); c != "" {
		fmt.Sscanf(c, "%d", &count)
		if count > 100 {
			count = 100
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	if count == 1 {
		n := rand.Intn(max-min+1) + min
		fmt.Fprintf(w, `{"number":%d}`, n)
	} else {
		nums := make([]int, count)
		for i := range nums {
			nums[i] = rand.Intn(max-min+1) + min
		}
		fmt.Fprintf(w, `{"numbers":[`)
		for i, n := range nums {
			if i > 0 { fmt.Fprintf(w, ",") }
			fmt.Fprintf(w, "%d", n)
		}
		fmt.Fprintf(w, `]}`)
	}
}

func getWebDir() string {
	for _, p := range []string{"web", filepath.Dir(os.Args[0]) + "/web"} {
		if _, err := os.Stat(p); err == nil { return p }
	}
	return "web"
}

func getLocalIP() string {
	for _, i := range getIfaces() {
		for _, a := range i.Addrs {
			if ip, ok := a.(*net.IPNet); ok && ip.IP.To4() != nil && !ip.IP.IsLoopback() {
				return ip.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func getIfaces() []net.Interface {
	var r []net.Interface
	for _, i := range getAllIfaces() {
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagLoopback == 0 {
			r = append(r, i)
		}
	}
	return r
}

func getAllIfaces() []net.Interface {
	ifaces, _ := net.Interfaces()
	return ifaces
}

func init() { 
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}