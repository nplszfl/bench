package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
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
	http.HandleFunc("/api/hash", handleHash)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleHash(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	algo := r.URL.Query().Get("algo")
	
	md5Hash := md5.Sum(data)
	sha1Hash := sha1.Sum(data)
	sha256Hash := sha256.Sum256(data)
	
	w.Header().Set("Content-Type", "application/json")
	
	switch algo {
	case "md5":
		fmt.Fprintf(w, `{"algo":"md5","hash":"%s"}`, hex.EncodeToString(md5Hash[:]))
	case "sha1":
		fmt.Fprintf(w, `{"algo":"sha1","hash":"%s"}`, hex.EncodeToString(sha1Hash[:]))
	case "sha256":
		fmt.Fprintf(w, `{"algo":"sha256","hash":"%s"}`, hex.EncodeToString(sha256Hash[:]))
	default:
		fmt.Fprintf(w, `{"md5":"%s","sha1":"%s","sha256":"%s"}`,
			hex.EncodeToString(md5Hash[:]),
			hex.EncodeToString(sha1Hash[:]),
			hex.EncodeToString(sha256Hash[:]))
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