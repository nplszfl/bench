package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var port = flag.Int("port", 8080, "端口")

func main() {
	flag.Parse()
	localIP := getLocalIP()
	
	fmt.Printf("\n🛠️  Bench - 开发工具箱\n")
	fmt.Printf("访问: http://%s:%d\n", localIP, *port)
	
	webDir := getWebDir()
	http.HandleFunc("/api/convert", handleConvert)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Query().Get("value")
	base := r.URL.Query().Get("base")
	
	w.Header().Set("Content-Type", "application/json")
	
	fromBase := 10
	fmt.Sscanf(base, "%d", &fromBase)
	
	// 统一转换为十进制
	n, err := strconv.ParseInt(value, fromBase, 64)
	if err != nil {
		fmt.Fprintf(w, `{"error":"%s"}`, err)
		return
	}
	
	fmt.Fprintf(w, `{"decimal":%d,"binary":"%s","octal":"%s","hex":"%s"}`,
		n, strconv.FormatInt(n, 2), strconv.FormatInt(n, 8), strconv.FormatInt(n, 16))
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

func init() { flag.Parse() }