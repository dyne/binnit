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

/*
 *
 * Store/Retrieve functions for FS-based paste storage
 *
 */

package main

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func store(title, date, content, destDir, lang string) (string, error) {

	h := sha256.New()

	h.Write([]byte(title))
	h.Write([]byte(date))
	h.Write([]byte(lang))
	h.Write([]byte(content))

	paste := fmt.Sprintf("# Title: %s\n# Date: %s\n# Language: %s\n%s", title, date, lang, content)

	pasteHash := fmt.Sprintf("%x", h.Sum(nil))
	log.Printf("  `-- hash: %s\n", pasteHash)
	pasteDir := destDir + "/"

	// Now we save the file
	for i := 0; i < len(pasteHash)-16; i++ {
		pasteName := pasteHash[i : i+16]
		if _, err := os.Stat(pasteDir + pasteName); os.IsNotExist(err) {
			// The file does not exist, so we can create it
			if err := ioutil.WriteFile(pasteDir+pasteName, []byte(paste), 0644); err == nil {
				// and then we return the URL:
				log.Printf("  `-- saving new paste to : %s", pasteDir+pasteName)
				return pasteName, nil
			}
			log.Printf("Cannot create the paste: %s!\n", pasteDir+pasteName)

		}
	}
	return "", errors.New("cannot store the paste...sorry")
}

func retrieve(URI string) (title string, date string, lang string, content string, err error) {

	fCont, err := os.Open(URI)
	defer fCont.Close()

	stuff := bufio.NewScanner(fCont)
	// The first line contains the title
	stuff.Scan()
	title = strings.Trim(strings.Split(stuff.Text(), ":")[1], " ")
	stuff.Scan()
	date = strings.Trim(strings.Join(strings.Split(stuff.Text(), ":")[1:], ":"), " ")
	stuff.Scan()
	lang = strings.Trim(strings.Split(stuff.Text(), ":")[1], " ")
	for stuff.Scan() {
		content += stuff.Text() + "\n"
	}

	return
}
