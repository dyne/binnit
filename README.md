## binnit -- minimal pastebin-like in golang

That's just it. Preliminary version of a minimal, no-fuss
pastebin-like service in golang. 

It serves pastes in the format:

    mypasteserver.org/abcdef1234567890

and stores them in a folder, one file per paste, whose filename is
equal to the paste ID. The unique ID of a paste is obtained from the
SHA256 of the concatenation of title, time, and content. Rendering is
minimal, but can be enhanced.

`binnit` is currently configured through a simple key=value
configuration file, whose name can be specified on the command line
through the option `-c <config\_file>`. The available options are:

* server\_name  (the FQDN where the service is reachable from outside)
* bind\_addr (the address to listen on)
* bind\_port (the port to bind)
* paste\_dir (the folder where pastes are kept)
* templ\_dir (the folder where HTML files and templates are kept)
* max_size (the maximum allowed length of a paste, in bytes. Larger
    pastes will be trimmed to that length)
* log_fname (path to the logfile)



### TODO

* provide a better standard template
