package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var port = flag.Int("port", 8080, "端口")

func main() {
	flag.Parse()
	localIP := getLocalIP()
	
	fmt.Printf("\n🛠️  Bench - 开发工具箱\n")
	fmt.Printf("访问: http://%s:%d\n", localIP, *port)
	
	webDir := getWebDir()
	http.HandleFunc("/api/diff", handleDiff)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handleDiff(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
	
	json.NewDecoder(r.Body).Decode(&req)
	
	oldLines := strings.Split(req.Old, "\n")
	newLines := strings.Split(req.New, "\n")
	
	var diff []string
	max := len(oldLines)
	if len(newLines) > max {
		max = len(newLines)
	}
	
	for i := 0; i < max; i++ {
		var oldLine, newLine string
		if i < len(oldLines) {
			oldLine = oldLines[i]
		}
		if i < len(newLines) {
			newLine = newLines[i]
		}
		
		if oldLine != newLine {
			if oldLine != "" {
				diff = append(diff, fmt.Sprintf("- %s", oldLine))
			}
			if newLine != "" {
				diff = append(diff, fmt.Sprintf("+ %s", newLine))
			}
		} else {
			diff = append(diff, fmt.Sprintf("  %s", oldLine))
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"diff":%s}`, toJSONArray(diff))
}

func toJSONArray(arr []string) string {
	result := "["
	for i, s := range arr {
		if i > 0 { result += "," }
		result += fmt.Sprintf(`"%s"`, escapeJSON(s))
	}
	return result + "]"
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
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