# basic info
app := maxwell
module := github.com/dstgo/maxwell/cmd/$(app)
output := $(shell pwd)/bin
# meta info
build_time := $(shell date +"%Y/%m/%dT%H:%M:%SZ%z")
git_version := $(shell git describe --tags --always)
# build info
host_os := $(shell go env GOHOSTOS)
host_arch := $(shell go env GOHOSTARCH)
os := $(host_os)
arch := $(host_arch)

ifeq ($(os), windows)
	exe := .exe
endif


.PHONY: build
build:
	# go lint
	go vet ./...

	# prepare target environment $(os)/$(arch)
	go env -w GOOS=$(os)
	go env -w GOARCH=$(arch)

	# build go module
	go build -a -trimpath \
		-ldflags="-X main.AppName=$(app) -X main.Version=$(git_version) -X main.BuildTime=$(build_time)" \
		-o $(output)/$(app)$(exe) \
		$(module)

	# resume host environment $(host_os)/$(host_arch)
	go env -w GOOS=$(host_os)
	go env -w GOARCH=$(host_arch)

# ent schema path
schema =
ent_out := ./app/data/ent
ent_target := schema
ent_generated := $(shell find $(ent_out)/* ! -path "*$(ent_target)*")

.PHONY: ent_gen
ent_gen:
ifneq ($(schema),)
	# generate new $(schema) schema
	ent new --target $(ent_out)/$(ent_target) $(schema)
endif
	# generate ent code
	ent generate $(ent_out)/$(ent_target)

.PHONY: ent_clean
ent_clean:
	@rm -rf $(ent_generated)

api_path := ./app/api

.PHONY: swag
swag:
	go generate $(api_path)