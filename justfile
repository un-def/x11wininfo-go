project := 'x11wininfo'

export GOPATH := invocation_directory()

common_build_args := "-a -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH"
static_build_args := "-ldflags '-w -linkmode external -extldflags \"-static -lXau -lXdmcp\"'"

build +extra_build_args='':
  go build -o {{project}} {{common_build_args}} {{extra_build_args}}

build-static +extra_build_args='':
  go build -o {{project}}-static {{common_build_args}} {{static_build_args}} {{extra_build_args}}

run +run_args='':
  go run main.go {{run_args}}
