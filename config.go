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
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type options struct {
	confFile string
}

type config struct {
	serverName string
	bindAddr   string
	bindPort   string
	pasteDir   string
	templDir   string
	staticDir  string
	maxSize    uint16
	logFile    string
}

func (c config) String() string {

	var s string

	s += "Server name: " + c.serverName + "\n"
	s += "Listening on: " + c.bindAddr + ":" + c.bindPort + "\n"
	s += "paste_dir: " + c.pasteDir + "\n"
	s += "templ_dir: " + c.templDir + "\n"
	s += "static_dir: " + c.staticDir + "\n"
	s += "max_size: " + string(c.maxSize) + "\n"
	s += "log_file: " + c.logFile + "\n"

	return s

}

func parseConfig(fname string, c *config) error {

	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	r := bufio.NewScanner(f)

	line := 0
	for r.Scan() {
		s := r.Text()
		line++
		if matched, _ := regexp.MatchString("^([ \t]*)$", s); matched != true {
			// it's not a blank line
			if matched, _ := regexp.MatchString("^#", s); matched != true {
				// This is not a comment...
				if matched, _ := regexp.MatchString("^([a-z_ ]+)=.*", s); matched == true {
					// and contains an assignment
					fields := strings.Split(s, "=")
					switch strings.Trim(fields[0], " \t\"") {
					case "server_name":
						c.serverName = strings.Trim(fields[1], " \t\"")
					case "bind_addr":
						c.bindAddr = strings.Trim(fields[1], " \t\"")
					case "bind_port":
						c.bindPort = strings.Trim(fields[1], " \t\"")
					case "paste_dir":
						c.pasteDir = strings.Trim(fields[1], " \t\"")
					case "templ_dir":
						c.templDir = strings.Trim(fields[1], " \t\"")
					case "static_dir":
						c.staticDir = strings.Trim(fields[1], " \t\"")
					case "log_file":
						c.logFile = strings.Trim(fields[1], " \t\"")
					case "max_size":
						if mSize, err := strconv.ParseUint(fields[1], 10, 16); err == nil {
							c.maxSize = uint16(mSize)
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
