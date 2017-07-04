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
 * minimal Templating support for binnit
 *
 */

package main

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"strconv"
)

func format_rows(content string) (string) {

	var ret string

	lines := strings.Split(content, "\n")

	ret += "<table class='content'>"
	
	for l_num, l := range lines {
		ret += "<tr>\n"
		ret += "<td class='lineno'><pre>"+ strconv.Itoa(l_num+1) + "</pre></td>"
		ret += "<td class='line'><pre>"+ l +"</pre></td>"
		ret += "</tr>"
	}
	ret += "</table>"
	return ret
}


func prepare_paste_page(title, date, content, templ_dir string) (string, error) {

	s := ""

	// insert header

	head_file := templ_dir + "/header.html"

	f_head, err := os.Open(head_file)
	defer f_head.Close()

	if err == nil {
		cont, err := ioutil.ReadFile(head_file)
		if err == nil {
			s += string(cont)
		}
	}

	// insert content

	// ...Let's read the template
	templ_file := templ_dir + "/templ.html"
	f_templ, err := os.Open(templ_file)
	defer f_templ.Close()

	
	if cont, err := ioutil.ReadFile(templ_file); err == nil {
		tmpl := string(cont)
		// ...and replace  {{CONTENT}} with the paste itself!
		re, _ := regexp.Compile("{{TITLE}}")
		tmpl = string(re.ReplaceAll([]byte(tmpl), []byte(title)))

		re, _ = regexp.Compile("{{DATE}}")
		tmpl = string(re.ReplaceAll([]byte(tmpl), []byte(date)))

		re, _ = regexp.Compile("{{CONTENT}}")
		tmpl = string(re.ReplaceAll([]byte(tmpl), []byte(format_rows(content))))

		re, _ = regexp.Compile("{{RAW_CONTENT}}")
		tmpl = string(re.ReplaceAll([]byte(tmpl), []byte(content)))

		s += tmpl
		
	} else {
		return "", errors.New("Error opening template file")
	}
	
	// insert footer
	foot_file := templ_dir + "/footer.html"
	f_foot, err := os.Open(foot_file)
	defer f_foot.Close()

	if err == nil {
		cont, err := ioutil.ReadFile(foot_file)
		if err == nil {
			s += string(cont)
		}
	}
	return s, nil
}
