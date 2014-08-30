.PHONY: all 

all: build write read

build:
	go build benchmark.go

write:
	./benchmark -type=write

read:
	./benchmark -type=read

