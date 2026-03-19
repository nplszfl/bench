package main

import (
	"flag"
	"fmt"
	"log"
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
	http.HandleFunc("/api/now", handleNow)
	http.HandleFunc("/api/convert", handleConvert)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleNow(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"unix":%d,"iso":"%s","date":"%s","time":"%s"}`,
		now.Unix(), now.Format(time.RFC3339),
		now.Format("2006-01-02"), now.Format("15:04:05"))
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	timestamp := r.URL.Query().Get("timestamp")
	format := r.URL.Query().Get("format")
	
	w.Header().Set("Content-Type", "application/json")
	
	if timestamp != "" {
		var t int64
		fmt.Sscanf(timestamp, "%d", &t)
		
		// 毫秒转秒
		if t > 10000000000 {
			t = t / 1000
		}
		
		tm := time.Unix(t, 0)
		if format != "" {
			fmt.Fprintf(w, `{"result":"%s"}`, tm.Format(format))
		} else {
			fmt.Fprintf(w, `{"unix":%d,"iso":"%s","local":"%s"}`,
				tm.Unix(), tm.Format(time.RFC3339), tm.Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Fprintf(w, `{"error":"missing timestamp"}`)
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

func init() { flag.Parse() }