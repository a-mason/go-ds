# Install protobuf onto the machine
# Figure out which machine we're running on.
# Run the right installation command for the operating system.
install-protobuf:
UNAME := $(shell uname)
PROTOC := $(shell which protoc)
$(PROTOC):
ifeq ($(UNAME), Darwin)
	brew install protobuf
	brew install protoc-gen-go
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif

# Get the go proto modules.
go-get-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/protobuf

build-proto:
	protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.

# This is a "phony" target - an alias for the above command, so "make compile"
# still works.
all: install-protobuf go-get-proto build-proto
