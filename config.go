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
	"os"
	"bufio"
	"regexp"
	"strings"
	"strconv"
)

type Options struct {
	conf_file string
}


type Config struct {
	server_name string
	bind_addr string 
	bind_port string
	paste_dir string
	templ_dir string
	max_size uint16
	log_file string
}


func (c Config) String() string {

	var s string

	s+= "Server name: " + c.server_name + "\n"
	s+= "Listening on: " + c.bind_addr + ":" + c.bind_port +"\n"
	s+= "paste_dir: " + c.paste_dir + "\n"
	s+= "templ_dir: " + c.templ_dir + "\n"
	s+= "max_size: " + string(c.max_size) + "\n"
	s+= "log_file: " + c.log_file + "\n"

	return s
	
}

func parse_config (fname string, c *Config) error {

	
	f, err := os.Open(fname);
	if  err != nil {
		return err
	}

	r := bufio.NewScanner(f)

	line := 0
	for r.Scan (){
		s := r.Text()
		line += 1
		if matched, _ := regexp.MatchString("^([ \t]*)$", s); matched != true {
			// it's not a blank line
			if matched, _ := regexp.MatchString("^#", s); matched != true  {
				// This is not a comment...
				if matched, _ := regexp.MatchString("^([a-z_ ]+)=.*", s);  matched == true {
					// and contains an assignment
					fields := strings.Split(s, "=")
					switch strings.Trim(fields[0], " \t\""){
					case "server_name":
						c.server_name = strings.Trim(fields[1], " \t\"")
					case "bind_addr":
						c.bind_addr = strings.Trim(fields[1], " \t\"")
					case "bind_port":
						c.bind_port = strings.Trim(fields[1], " \t\"")
					case "paste_dir":
						c.paste_dir = strings.Trim(fields[1], " \t\"")
					case "templ_dir":
						c.templ_dir = strings.Trim(fields[1], " \t\"")
					case "log_file":
						c.log_file = strings.Trim(fields[1], " \t\"")
					case "max_size":
						if m_size, err := strconv.ParseUint(fields[1], 10, 16); err == nil {
							c.max_size = uint16(m_size)
						} else {
							fmt.Fprintf(os.Stderr, "Invalid max_size value %s at line %d (max: 65535)\n",
								fields[1], line)
						}
					default:
						fmt.Fprintf(os.Stderr, "Error reading config file %s at line %d: unknown variable '%s'\n",
							fname, line, fields[0])
					}
				} else {
					fmt.Fprintf(os.Stderr, "Error reading config file %s at line %d: unknown statement '%s'\n",
						fname, line, s)
				}
			}
		}
	}
	return nil
}


