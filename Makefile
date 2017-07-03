GO=go

all: binnit

binnit: binnit.go templ.go config.go 
	$(GO) build -o binnit binnit.go templ.go config.go

