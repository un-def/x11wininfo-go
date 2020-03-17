package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"x11wininfo/pkg/x11"
)

var modes = [...]string{
	"text",
	"mintext",
	"json",
}

type jsonResponse struct {
	ID       x11.WindowID `json:"id,string"`
	Name     string       `json:"name"`
	Instance string       `json:"instance"`
	Class    string       `json:"class"`
}

func die(err error, conn *x11.Connection) {
	fmt.Fprintln(os.Stderr, err)
	if conn != nil {
		conn.Disconnect()
	}
	os.Exit(1)
}

func isModeSupported(mode string) bool {
	for _, supportedMode := range modes {
		if mode == supportedMode {
			return true
		}
	}
	return false
}

func main() {
	modePtr := flag.String("m", modes[0], fmt.Sprintf("output `mode`: %v", modes))
	flag.Parse()
	mode := *modePtr
	if !isModeSupported(mode) {
		die(fmt.Errorf("unsupported mode: %s", mode), nil)
	}
	conn, err := x11.Connect()
	if err != nil {
		die(err, conn)
	}
	window, err := conn.FocusedWindow()
	if err != nil {
		die(err, conn)
	}
	id := window.ID()
	name, err := window.Name()
	if err != nil {
		die(err, conn)
	}
	instance, class, err := window.Class()
	if err != nil {
		die(err, conn)
	}
	switch mode {
	case "text":
		fmt.Printf(
			"id: %d\nname: %s\ninstance: %s\nclass: %s\n",
			id, name, instance, class,
		)
	case "mintext":
		fmt.Printf(
			"%d\n%s\n%s\n%s\n",
			id, name, instance, class,
		)
	case "json":
		response, err := json.Marshal(jsonResponse{
			ID:       id,
			Name:     name,
			Instance: instance,
			Class:    class,
		})
		if err != nil {
			die(err, conn)
		}
		fmt.Println(string(response))
	}

}
