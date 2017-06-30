## binit -- minimal pastebin-like in 100 lines of golang

That's just it. Preliminary version of a minimal, no-fuss
pastebin-like service in golang. 

Needs a folder "pastes/" to exist on the same dir where the program is
run from. At the moment, it binds on `localhost:8000` and serves
pastes in the format:

    localhost:8000/abcdef1234567890

The unique ID of a paste is obtained from the SHA256 of the
concatenation of title, time, and content. Rendering is minimal, but
can be enhanced. 

`binit` is currently configured through a simple key=value
configuration file. The available options are:

* host (the hostname to listen on)
* port (the port to bind)
* paste\_dir (the folder where pastes are kept)
* templ\_dir (the folder where HTML files and templates are kept)
* max_size (the maximum allowed length of a paste, in bytes. Larger
    pastes will be trimmed to that length)


### TODO

* Add a simple template system
