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
	http.HandleFunc("/api/generate", handleGenerate)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	length := 16
	hasUpper := true
	hasLower := true
	hasNumber := true
	hasSpecial := false
	
	if l := r.URL.Query().Get("length"); l != "" {
		fmt.Sscanf(l, "%d", &length)
	}
	
	charset := "abcdefghijklmnopqrstuvwxyz"
	if hasUpper {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if hasNumber {
		charset += "0123456789"
	}
	if hasSpecial {
		charset += "!@#$%^&*"
	}
	
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"password":"%s","length":%d}`, result, length)
}

func getWebDir() string {
	for _, p := range []string{"web", filepath.Dir(os.Args[0]) + "/web"} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "web"
}

func getLocalIP() string {
	for _, iface := range getIfaces() {
		for _, addr := range iface.Addrs {
			if ip, ok := addr.(*net.IPNet); ok {
				if ip.IP.To4() != nil && !ip.IP.IsLoopback() {
					return ip.IP.String()
				}
			}
		}
	}
	return "127.0.0.1"
}

func getIfaces() []net.Interface {
	ifaces, _ := net.Interfaces()
	var result []net.Interface
	for _, i := range ifaces {
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagLoopback == 0 {
			result = append(result, i)
		}
	}
	return result
}

func init() { 
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}