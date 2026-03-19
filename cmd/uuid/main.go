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
	http.HandleFunc("/api/uuid", handleUUID)
	http.HandleFunc("/api/uuid-batch", handleUUIDBatch)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleUUID(w http.ResponseWriter, r *http.Request) {
	count := 1
	if c := r.URL.Query().Get("count"); c != "" {
		fmt.Sscanf(c, "%d", &count)
		if count > 100 {
			count = 100
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"uuids":[`)
	for i := 0; i < count; i++ {
		if i > 0 { fmt.Fprintf(w, ",") }
		fmt.Fprintf(w, `"%s"`, generateUUID())
	}
	fmt.Fprintf(w, `]}`)
}

func handleUUIDBatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Count int `json:"count"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	if req.Count <= 0 || req.Count > 1000 {
		req.Count = 10
	}
	
	uuids := make([]string, req.Count)
	for i := 0; i < req.Count; i++ {
		uuids[i] = generateUUID()
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"uuids": uuids, "count": req.Count})
}

func generateUUID() string {
	// 简化的UUID v4生成
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
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