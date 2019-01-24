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

func formatRaw(content string) string {
	var ret string
	ret += "<pre>\n<code>\n"
	ret += content
	ret += "\n</code>\n<pre>\n"
	return ret
}

func preparePastePage(title, date, content, templDir string, raw bool) (string, error) {

	s := ""
	if !raw {
		templFile := templDir + "/paste.html"

		fTempl, err := os.Open(templFile)
		if err != nil {
			return "", errors.New("Error opening template file")
		}
		defer fTempl.Close()

		if cont, err := ioutil.ReadFile(templFile); err == nil {
			tmpl := string(cont)
			re, _ := regexp.Compile("{{TITLE}}")
			tmpl = string(re.ReplaceAllLiteralString(tmpl, title))

			re, _ = regexp.Compile("{{DATE}}")
			tmpl = string(re.ReplaceAllLiteralString(tmpl, date))
			re, _ = regexp.Compile("{{CONTENT}}")
			tmpl = string(re.ReplaceAllLiteralString(tmpl, formatRows(content)))
			s += tmpl
		} else {
			return "", errors.New("Error opening template file")
		}
	} else {
		s += "<html>\n<head>\n</head>\n<body>\n"
		s += formatRaw(content)
		s += "</body>\n</html>"
	}
	return s, nil
}
