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

### TODO

* Check maximum paste length
* Add a config file (hostname, port, pastedir)
* Add a simple template system
