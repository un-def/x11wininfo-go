project := 'x11wininfo'

export GOPATH := invocation_directory()


build +build_args='':
  go build -o {{project}} -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH {{build_args}}

run +run_args='':
  go run main.go {{run_args}}
