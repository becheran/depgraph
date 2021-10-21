package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
)

//go:embed www/*
var staticWebSite embed.FS

func openbrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}

func Run(host string, handler func(http.ResponseWriter, *http.Request)) {
	files, err := fs.Sub(staticWebSite, "www")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(files)))

	http.HandleFunc("/api", handler)

	url := "http://" + host
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Unable to server files. Err: %s", err.Error())
	}
	defer listener.Close()

	fmt.Printf("Listen on %s", url)

	err = openbrowser(url)
	if err != nil {
		log.Printf("Failed to open browser. Err: %s", err.Error())
	}

	if err := http.Serve(listener, nil); err != nil {
		log.Fatal(err)
	}
}
