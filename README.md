## binnit -- minimal pastebin clone in golang

That's just it. A minimalist, no-fuss pastebin clone server in
golang. It supports only two operations:

* store a new paste, through a POST request
* retrieve a paste using its unique ID, through a GET request

what else do you need? 

## WTF?

`binnit` is a single executable with **no dependencies**. You **don't
need** a web server. You **don't need** a SQL server. You **don't
need** any external library. 

`binnit` serves pastes in the format:

    http://<server_name>/abcdef1234567890

and stores them in a folder on the server, one file per paste, whose
filename is identical to the paste ID. The unique ID of a paste is
obtained from the SHA256 of the concatenation of title, time, and
content. Rendering is minimal, on purpose, but based on a customisable
template.

`binnit` is currently configured through a simple key=value
configuration file, whose name can be specified on the command line
through the option `-c <config_file>`. If no config file is specified,
`binnit` looks for `./binnit.cfg`. The configurable options are:

* server\_name  (the FQDN where the service is reachable from outside)
* bind\_addr (the address to listen on)
* bind\_port (the port to bind)
* paste\_dir (the folder where pastes are kept)
* templ\_dir (the folder where HTML files and templates are kept)
* max\_size (the maximum allowed length of a paste, in bytes. Larger
    pastes will be trimmed to that length.)
* log_file (path to the logfile)

As with other pastebin-like services, you can send a paste to `binnit`
using `curl`. For instance, if your `binnit` server is running on
`http://servername.net`, you can paste a file there using:


    curl -F 'paste=<myfile' http://servername.net


and obtain on output the ID associated to the newly created
paste. Similarly

    mylongcommand | curl -F 'paste=<-' http://servername.net

will paste the output of `mylongcommand` to `http://servername.net`,
and show on output the ID of the new paste.


## Why another pastebin?

There are hundreds of pastebin-like servers in the wild. But the
overwhelming majority of them is _overbloated_ software, depending on
lots of libraries/frameworks/tools, providing a whole lot of useless
features, and implying a useless amount of complexity. 

A paste server must be able to do two things, 1) create a new paste
and return its ID, and 2) retrieve an existing paste using its
ID. `binnit` does just and only these two things, in the simplest
possible way, without any external dependency. If you need more than
that, then `binnit` is not for you. But do you really need anything
more?

## About minimalism

> It seems that perfection is attained not when there is nothing more
> to add, but when there is nothing more to remove (Antoine de Saint
> Exupéry)

`binnit` is intended to be truly minimal. It consists of about 500
lines of golang source code in total, including:

* ~110 lines for License statements (comments)
* ~110 lines of core logic
* ~90 blank lines
* ~75 lines for template management
* ~75 lines for config management
* ~30 lines of pure comments

If you want to strip `binnit` down even further, you could consider
removing:

* blank lines
* the external configuration file
* the template system 
* sanity checks and error management
* logging 
* code comments

You **CANNOT** remove the licence statements on each source file.


## LICENSE

`binnit` is Copyright (2017) by Vincenzo "KatolaZ" Nicosia.

`binnit` is free software. You can use, modify, and redistribute it
under the terms of the GNU Affero General Public Licence, version 3 of
the Licence or, at your option, any later version. Please see
LICENSE.md for details.

