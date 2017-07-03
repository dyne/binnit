/*
 *  This program is free software: you can redistribute it and/or
 *  modify it under the terms of the GNU General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 *  General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see
 *  <http://www.gnu.org/licenses/>.
 *
 *  (c) Vincenzo "KatolaZ" Nicosia 2017 -- <katolaz@freaknet.org>
 * 
 * 
 *  This file is part of "binnit", a minimal no-fuss pastebin-like 
 *  server written in golang
 *
 */


package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"io"
	"binnit/paste"
)


var p_conf = Config{
	server_name: "localhost",
	bind_addr: "0.0.0.0",
	bind_port: "8000",
	paste_dir: "./pastes",
	templ_dir: "./tmpl",
	max_size: 4096,
	log_file: "./binnit.log",
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
			//http.ServeFile(w, r, paste_name)
			s, err := prepare_paste_page(&p_conf, orig_name)
			if err == nil {
				fmt.Fprintf(w, "%s", s)
				return
			} else {
				fmt.Fprintf(w, "Error recovering paste '%s'\n", orig_name)
				return
			}
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
		req_body := r.PostForm

		orig_IP := r.RemoteAddr

		log.Printf("Received new POST from %s\n", orig_IP)
		
		// get title, body, and time
		title := req_body.Get("title")
		date := time.Now().String()
		content := req_body.Get("paste")
		
		content = content[0:min(len(content), int(p_conf.max_size))]

		ID, err := paste.Store(title, date, content, p_conf.paste_dir)
		
		log.Printf("   ID: %s; err: %s\n", ID, err)

		if  err == nil {
			hostname := p_conf.server_name
			if show := req_body.Get("show"); show != "1" {
				fmt.Fprintf(w, "http://%s/%s", hostname, ID)
				return
			} else{
				fmt.Fprintf(w, "<html><body>Link: <a href='http://%s/%s'>http://%s/%s</a></body></html>",
					hostname, ID, hostname, ID)
				return
			}
		} else {
			fmt.Fprintf(w, "%s\n", err)
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

	
	
	parse_config("binnit.cfg", &p_conf)
	

	f, err := os.OpenFile(p_conf.log_file, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log_file: %s. Exiting\n", p_conf.log_file)
		os.Exit(1)
	}
	defer f.Close()

	
	log.SetOutput(io.Writer(f))
	log.SetPrefix("[binnit]: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("Binnit version 0.1 -- Starting ")
	log.Printf("  + Serving pastes on: %s\n", p_conf.server_name)
	log.Printf("  + listening on: %s:%s\n", p_conf.bind_addr, p_conf.bind_port )
	log.Printf("  + paste_dir: %s\n", p_conf.paste_dir)
	log.Printf("  + templ_dir: %s\n", p_conf.templ_dir)
	log.Printf("  + max_size: %d\n", p_conf.max_size)

	// FIXME: create paste_dir if it does not exist
	
	http.HandleFunc("/", req_handler)
	log.Fatal(http.ListenAndServe(p_conf.bind_addr + ":" +  p_conf.bind_port, nil))
}
