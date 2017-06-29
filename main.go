package main

import (
	"fmt"
	"net/http"
	"crypto/sha256"
	"log"
	"path/filepath"
	"os"
	"time"
	"io/ioutil"
)


func handle_get_paste(w http.ResponseWriter, r *http.Request){
	
	var paste_name, orig_name string
	var err error

	orig_name = filepath.Clean(r.URL.Path)
	paste_name = "./pastes/" + orig_name

	fmt.Printf("orig_name: '%s'\npaste_name: '%s'\n", orig_name, paste_name)

	// The default is to serve index.html
	if (orig_name == "/" ) || ( orig_name == "/index.html" ) {
		http.ServeFile(w, r, "index.html")
	} else {
		// otherwise, if the requested paste exists, we serve it...
		if _, err = os.Stat(paste_name); err == nil && orig_name != "./" {
			http.ServeFile(w, r, paste_name)
		}	else {
			// otherwise, we give say we didn't find it 
			fmt.Fprintf(w, "Paste '%s' not found\n", orig_name)
			return
		}
	}
}



func handle_put_paste(w http.ResponseWriter, r *http.Request){


	fmt.Printf("We are inside handle_put_paste\n");
	
	if err := r.ParseForm(); err != nil{
		// Invalid POST -- let's serve the default file
		http.ServeFile(w, r, "index.html")
	} else {
		h := sha256.New()
		req_body := r.PostForm
		// get title, body, and time
		title := req_body.Get("title")
		paste := req_body.Get("paste")
		now := time.Now().String()
		// format content 
		content := fmt.Sprintf("# Title: %s\n# Pasted: %s\n------------\n%s", title, now, paste)
		
		// ccompute the sha256 hash using title, body, and time
		h.Write([]byte (content))

		paste_hash := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Printf("hash: %s fname: %s\n", paste_hash, paste_hash[:16])
		paste_dir := "./pastes/"
		
		// Now we save the file
		for i := 0; i < len(paste_hash) - 16; i ++ {

			if _, err := os.Stat(paste_dir + paste_hash[i:i+16]); os.IsNotExist(err) {
				// The file does not exist, so we can create it
				if err := ioutil.WriteFile(paste_dir + paste_hash[i:i+16], []byte (content), 0644); err == nil{
					// and then we return the URL:
					fmt.Fprintf(w, "<html><body>Link: <a href='%s'>%s</a></body></html>",
						paste_hash[i:i+16], paste_hash[i:i+16])
					return
				} else {
					fmt.Fprintf(w, "Cannot create the paste!!!\n")
				}
			}
		}
	}
}

func req_handler(w http.ResponseWriter, r *http.Request){

	switch r.Method {
	case "GET":
		handle_get_paste(w, r)
	case "POST":
		handle_put_paste(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main(){
	http.HandleFunc("/", req_handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
