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

/*
*
* minimal Templating support for binit
*
*/

package main


import (
	"os"
	"io/ioutil"
	"regexp"
	"errors"
)


func prepare_paste_page(c *Config, paste_ID string) (string, error) {

	s:= ""

	// insert header

	head_file := c.templ_dir + "/header.html"

	f_head, err := os.Open(head_file)
	defer f_head.Close()
	
	if  err == nil {
		cont, err := ioutil.ReadFile(head_file)
		if err == nil{
			s += string(cont)
		}
	}
	
	// insert content

	cont_file := c.paste_dir + "/" + paste_ID
	f_cont, err := os.Open(cont_file)
	defer f_cont.Close()

	if err == nil {
		// Let's read the content of the paste

		cont, err := ioutil.ReadFile(cont_file)
		if err == nil {
			paste_buf := string(cont)

			// ...Let's read the template
			templ_file := c.templ_dir + "/templ.html" 
			f_templ, err := os.Open(templ_file)
			defer f_templ.Close()

			cont, err := ioutil.ReadFile(templ_file)
			if err == nil {
				tmpl := string(cont)
				// ...and replace  {{CONTENT}} with the paste itself!
				re,_ := regexp.Compile("{{CONTENT}}")
				tmpl = string(re.ReplaceAll([]byte(tmpl), []byte(paste_buf)))
				
				s += tmpl
				
			} else {
				return "", errors.New("Error opening template file")
			}
			
		} else {
			return "", errors.New("Error opening paste")
		}
	}
	// insert footer
	foot_file := c.templ_dir + "/footer.html"
	f_foot, err := os.Open(foot_file)
	defer f_foot.Close()

	if  err == nil {
		cont, err := ioutil.ReadFile(foot_file)
		if err == nil{
			s += string(cont)
		}
	}
	
	return s, nil
}
