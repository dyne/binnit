/*
 *  This program is free software: you can redistribute it and/or
 *  modify it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 *  General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public
 *  License along with this program.  If not, see
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
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dyne/binnit/paste"
)

var confFile = flag.String("c", "./binnit.cfg", "Configuration file for binnit")

var pConf = config{
	serverName: "localhost",
	bindAddr:   "0.0.0.0",
	bindPort:   "8000",
	pasteDir:   "./pastes",
	templDir:   "./tmpl",
	maxSize:    4096,
	logFile:    "./binnit.log",
}

func min(a, b int) int {

	if a > b {
		return b
	}
	return a

}

func handleGetPaste(w http.ResponseWriter, r *http.Request) {

	var pasteName, origName string

	origName = filepath.Clean(r.URL.Path)
	pasteName = pConf.pasteDir + "/" + origName

	origIP := r.RemoteAddr

	log.Printf("Received GET from %s for  '%s'\n", origIP, origName)

	// The default is to serve index.html
	if (origName == "/") || (origName == "/index.html") {
		http.ServeFile(w, r, pConf.templDir+"/index.html")
	} else {
		// otherwise, if the requested paste exists, we serve it...

		title, date, content, err := paste.Retrieve(pasteName)

		title = html.EscapeString(title)
		date = html.EscapeString(date)
		content = html.EscapeString(content)

		if err == nil {
			s, err := preparePastePage(title, date, content, pConf.templDir)
			if err == nil {
				fmt.Fprintf(w, "%s", s)
				return
			}
			fmt.Fprintf(w, "Error recovering paste '%s'\n", origName)
			return

		}
		// otherwise, we give say we didn't find it
		fmt.Fprintf(w, "%s\n", err)
		return
	}
}

func handlePutPaste(w http.ResponseWriter, r *http.Request) {

	err1 := r.ParseForm()
	err2 := r.ParseMultipartForm(int64(2 * pConf.maxSize))

	if err1 != nil && err2 != nil {
		// Invalid POST -- let's serve the default file
		http.ServeFile(w, r, pConf.templDir+"/index.html")
	} else {
		reqBody := r.PostForm

		origIP := r.RemoteAddr

		log.Printf("Received new POST from %s\n", origIP)

		// get title, body, and time
		title := reqBody.Get("title")
		date := time.Now().String()
		content := reqBody.Get("paste")

		content = content[0:min(len(content), int(pConf.maxSize))]

		ID, err := paste.Store(title, date, content, pConf.pasteDir)

		log.Printf("   title: %s\npaste: %s\n", title, content)
		log.Printf("   ID: %s; err: %v\n", ID, err)

		if err == nil {
			hostname := pConf.serverName
			if show := reqBody.Get("show"); show != "1" {
				fmt.Fprintf(w, "http://%s/%s\n", hostname, ID)
				return
			}
			fmt.Fprintf(w, "<html><body>Link: <a href='http://%s/%s'>http://%s/%s</a></body></html>",
				hostname, ID, hostname, ID)
			return

		}
		fmt.Fprintf(w, "%s\n", err)

	}
}

func reqHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		handleGetPaste(w, r)
	case "POST":
		handlePutPaste(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {

	flag.Parse()

	parseConfig(*confFile, &pConf)

	f, err := os.OpenFile(pConf.logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file: %s. Exiting\n", pConf.logFile)
		os.Exit(1)
	}
	defer f.Close()

	log.SetOutput(io.Writer(f))
	log.SetPrefix("[binnit]: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("Binnit version 0.1 -- Starting ")
	log.Printf("  + Config file: %s\n", *confFile)
	log.Printf("  + Serving pastes on: %s\n", pConf.serverName)
	log.Printf("  + listening on: %s:%s\n", pConf.bindAddr, pConf.bindPort)
	log.Printf("  + paste_dir: %s\n", pConf.pasteDir)
	log.Printf("  + templ_dir: %s\n", pConf.templDir)
	log.Printf("  + max_size: %d\n", pConf.maxSize)

	// FIXME: create paste_dir if it does not exist

	http.HandleFunc("/", reqHandler)
	log.Fatal(http.ListenAndServe(pConf.bindAddr+":"+pConf.bindPort, nil))
}
