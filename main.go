package main

import (
	"fmt"
	"os"
	"x11wininfo/pkg/x11"
)

func die(err error, conn *x11.Connection) {
	fmt.Fprintln(os.Stderr, err)
	if conn != nil {
		conn.Disconnect()
	}
	os.Exit(1)
}

func main() {
	conn, err := x11.Connect()
	if err != nil {
		die(err, conn)
	}
	window, err := conn.FocusedWindow()
	if err != nil {
		die(err, conn)
	}
	name, err := window.Name()
	if err != nil {
		die(err, conn)
	}
	instance, class, err := window.Class()
	if err != nil {
		die(err, conn)
	}
	fmt.Printf("id: %d\n", window.ID)
	fmt.Printf("name: %s\n", name)
	fmt.Printf("instance: %s\n", instance)
	fmt.Printf("class: %s\n", class)
}
