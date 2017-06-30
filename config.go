package main


import (
	"fmt"
	"os"
	"bufio"
	"regexp"
	"strings"
	"log"
)


type Config struct {
	host string
	port string
	paste_dir string
	templ_dir string
	log_fname string
	logger *log.Logger
}


func (c Config) String() string {

	var s string

	s+= "Host: " + c.host + "\n"
	s+= "Port: " + c.port + "\n"
	s+= "paste_dir: " + c.paste_dir + "\n"
	s+= "templ_dir: " + c.templ_dir + "\n"

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
				if matched, _ := regexp.MatchString("^([a-z_]+)=.*", s);  matched == true {
					// and contains an assignment
					fields := strings.Split(s, "=")
					switch fields[0]{
					case "host":
						c.host = fields[1]
					case "port":
						c.port = fields[1]
					case "paste_dir":
						c.paste_dir = fields[1]					
					case "templ_dir":
						c.templ_dir = fields[1]
					case "log_fname":
						c.log_fname = fields[1]
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


