project := 'x11wininfo'

cwd := invocation_directory()

common_build_args := "-a -gcflags=-trimpath=" + cwd + " -asmflags=-trimpath=" + cwd
static_build_args := "-ldflags '-w -linkmode external -extldflags \"-static -lXau -lXdmcp\"'"
cmd_path := cwd + "/cmd/" + project + "/main.go"

build +extra_build_args='':
  go build -o {{project}} {{common_build_args}} {{extra_build_args}} {{cmd_path}}

build-static +extra_build_args='':
  go build -o {{project}}-static {{common_build_args}} {{static_build_args}} {{extra_build_args}} {{cmd_path}}

run +run_args='':
  go run {{cmd_path}} {{run_args}}
