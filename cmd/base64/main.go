package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var port = flag.Int("port", 8080, "端口")

func main() {
	flag.Parse()
	localIP := getLocalIP()
	
	fmt.Printf("\n🛠️  Bench - 开发工具箱\n")
	fmt.Printf("访问: http://%s:%d\n", localIP, *port)
	
	webDir := getWebDir()
	http.HandleFunc("/api/encode", handleEncode)
	http.HandleFunc("/api/decode", handleDecode)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleEncode(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	encoded := base64.StdEncoding.EncodeToString(data)
	fmt.Fprintf(w, `{"result":"%s","original_size":%d,"encoded_size":%d}`, encoded, len(data), len(encoded))
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(decoded)
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

func init() { flag.Parse() }