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
	"strconv"
	"strings"
)

func formatRows(content string) string {

	var ret string

	lines := strings.Split(content, "\n")

	ret += "<table class='content'>"

	for lNum, l := range lines {
		ret += "<tr>\n"
		ret += "<td class='lineno'><pre>" + strconv.Itoa(lNum+1) + "</pre></td>"
		ret += "<td class='line'><pre>" + l + "</pre></td>"
		ret += "</tr>"
	}
	ret += "</table>"
	return ret
}

func preparePastePage(title, date, content, templDir string) (string, error) {

	s := ""

	// insert header

	headFile := templDir + "/header.html"

	fHead, err := os.Open(headFile)
	defer fHead.Close()

	if err == nil {
		cont, err := ioutil.ReadFile(headFile)
		if err == nil {
			s += string(cont)
		}
	}

	// insert content

	// ...Let's read the template
	templFile := templDir + "/templ.html"
	fTempl, err := os.Open(templFile)
	defer fTempl.Close()

	if cont, err := ioutil.ReadFile(templFile); err == nil {
		tmpl := string(cont)
		re, _ := regexp.Compile("{{TITLE}}")
		tmpl = string(re.ReplaceAllLiteralString(tmpl, title))

		re, _ = regexp.Compile("{{DATE}}")
		tmpl = string(re.ReplaceAllLiteralString(tmpl, date))

		re, _ = regexp.Compile("{{CONTENT}}")
		tmpl = string(re.ReplaceAllLiteralString(tmpl, formatRows(content)))

		re, _ = regexp.Compile("{{RAW_CONTENT}}")
		tmpl = string(re.ReplaceAllLiteralString(tmpl, content))

		s += tmpl

	} else {
		return "", errors.New("Error opening template file")
	}

	// insert footer
	footFile := templDir + "/footer.html"
	fFoot, err := os.Open(footFile)
	defer fFoot.Close()

	if err == nil {
		cont, err := ioutil.ReadFile(footFile)
		if err == nil {
			s += string(cont)
		}
	}
	return s, nil
}
