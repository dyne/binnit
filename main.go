package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"io"
)


var p_conf = Config{
	host: "localhost",
	port: "8000",
	paste_dir: "./pastes",
	templ_dir: "./tmpl",
	log_fname: "./binit.log",
	max_size: 4096,
}



func min (a, b int) int {

	if a > b {
		return b
	} else {
		return a
	}

}

func handle_get_paste(w http.ResponseWriter, r *http.Request) {

	var paste_name, orig_name string
	var err error

	orig_name = filepath.Clean(r.URL.Path)
	paste_name = p_conf.paste_dir + "/" + orig_name

	orig_IP := r.RemoteAddr
	
	log.Printf("Received GET from %s for  '%s'\n", orig_IP, orig_name)

	// The default is to serve index.html
	if (orig_name == "/") || (orig_name == "/index.html") {
		http.ServeFile(w, r, p_conf.templ_dir + "/index.html")
	} else {
		// otherwise, if the requested paste exists, we serve it...
		if _, err = os.Stat(paste_name); err == nil && orig_name != "./" {
			http.ServeFile(w, r, paste_name)
		} else {
			// otherwise, we give say we didn't find it
			fmt.Fprintf(w, "Paste '%s' not found\n", orig_name)
			return
		}
	}
}

func handle_put_paste(w http.ResponseWriter, r *http.Request) {

	
	if err := r.ParseForm(); err != nil {
		// Invalid POST -- let's serve the default file
		http.ServeFile(w, r, p_conf.templ_dir + "/index.html")
	} else {
		h := sha256.New()
		req_body := r.PostForm

		orig_IP := r.RemoteAddr

		log.Printf("Received new POST from %s\n", orig_IP)
		
		// get title, body, and time
		title := req_body.Get("title")
		paste := req_body.Get("paste")
		now := time.Now().String()
		// format content

		paste = paste[0:min(len(paste), int(p_conf.max_size))]
		
		content := fmt.Sprintf("# Title: %s\n# Pasted: %s\n------------\n%s", title, now, paste)

		// ccompute the sha256 hash using title, body, and time
		h.Write([]byte(content))

		paste_hash := fmt.Sprintf("%x", h.Sum(nil))
		log.Printf("  `-- hash: %s\n", paste_hash)
		paste_dir := p_conf.paste_dir + "/"

		// Now we save the file
		for i := 0; i < len(paste_hash)-16; i++ {
			paste_name := paste_hash[i:i+16]
			if _, err := os.Stat(paste_dir + paste_name); os.IsNotExist(err) {
				// The file does not exist, so we can create it
				if err := ioutil.WriteFile(paste_dir+ paste_name, []byte(content), 0644); err == nil {
					// and then we return the URL:
					log.Printf("  `-- saving paste to : %s", paste_dir + paste_name)
					hostname := r.Host
					if show := req_body.Get("show"); show != "1" {
						fmt.Fprintf(w, "%s/%s", hostname, paste_name)
						return
					} else{
						fmt.Fprintf(w, "<html><body>Link: <a href='%s'>%s</a></body></html>",
							paste_hash[i:i+16], paste_hash[i:i+16])
						return
					}
				} else {
					fmt.Fprintf(w, "Cannot create the paste.. Sorry!\n")
					return
				}
			}
		}
	}
}

func req_handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		handle_get_paste(w, r)
	case "POST":
		handle_put_paste(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {

	
	
	parse_config("binit.cfg", &p_conf)
	

	f, err := os.OpenFile(p_conf.log_fname, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening logfile: %s. Exiting\n", p_conf.log_fname)
		os.Exit(1)
	}
	
	log.SetOutput(io.Writer(f))
	log.SetPrefix("[binit]: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("Binit version 0.1 -- Starting ")
	log.Printf("  + listening on: %s:%s\n", p_conf.host, p_conf.port )
	log.Printf("  + paste_dir: %s\n", p_conf.paste_dir)
	log.Printf("  + templ_dir: %s\n", p_conf.templ_dir)
	log.Printf("  + max_size: %d\n", p_conf.max_size)

	
	http.HandleFunc("/", req_handler)
	log.Fatal(http.ListenAndServe(p_conf.host + ":" +  p_conf.port, nil))
}
