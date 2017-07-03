GO=go

all: binnit

binnit: main.go templ.go config.go
	$(GO) build -o binnit main.go templ.go config.go

