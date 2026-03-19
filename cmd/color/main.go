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
	"strings"
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
	color := r.URL.Query().Get("color")
	
	w.Header().Set("Content-Type", "application/json")
	
	var r_, g, b int
	if strings.HasPrefix(color, "#") {
		color = color[1:]
		if len(color) == 3 {
			r_, _ = strconv.ParseInt(string(color[0])+string(color[0]), 16, 16)
			g, _ = strconv.ParseInt(string(color[1])+string(color[1]), 16, 16)
			b, _ = strconv.ParseInt(string(color[2])+string(color[2]), 16, 16)
		} else if len(color) == 6 {
			r_, _ = strconv.ParseInt(color[0:2], 16)
			g, _ = strconv.ParseInt(color[2:4], 16)
			b, _ = strconv.ParseInt(color[4:6], 16)
		}
	}
	
	fmt.Fprintf(w, `{"hex":"#%02X%02X%02X","rgb":"rgb(%d,%d,%d)","r":%d,"g":%d,"b":%d}`,
		r_, g, b, r_, g, b, r_, g, b)
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