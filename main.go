package main

// #cgo pkg-config: xcb-atom
// #include <stdlib.h>
// #include <string.h>
// #include <xcb/xcb.h>
import "C"
import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

func die(err error, conn *C.xcb_connection_t) {
	fmt.Fprintln(os.Stderr, err)
	if conn != nil {
		C.xcb_disconnect(conn)
	}
	os.Exit(1)
}

func connect() (*C.xcb_connection_t, error) {
	conn := C.xcb_connect(nil, nil)
	if C.xcb_connection_has_error(conn) != 0 {
		return conn, errors.New("cannot connect to X server")
	}
	return conn, nil
}

func getScreen(conn *C.xcb_connection_t) (*C.xcb_screen_t, error) {
	screen := C.xcb_setup_roots_iterator(C.xcb_get_setup(conn)).data
	if screen == nil {
		return nil, errors.New("cannot get screen")
	}
	return screen, nil
}

func getAtom(conn *C.xcb_connection_t, name string) (C.xcb_atom_t, error) {
	nameCS := C.CString(name)
	cookie := C.xcb_intern_atom(conn, 0, C.ushort(C.strlen(nameCS)), nameCS)
	C.free(unsafe.Pointer(nameCS))
	reply := C.xcb_intern_atom_reply(conn, cookie, nil)
	if reply == nil {
		return 0, fmt.Errorf("cannot get %s atom", name)
	}
	atom := reply.atom
	C.free(unsafe.Pointer(reply))
	return atom, nil
}

func getWindowPropertyReply(
	conn *C.xcb_connection_t,
	window C.xcb_window_t,
	property string,
	propertyType string,
) (*C.xcb_get_property_reply_t, error) {

	propertyAtom, err := getAtom(conn, property)
	if err != nil {
		return nil, fmt.Errorf("cannot get %s property atom", property)
	}
	typeAtom, err := getAtom(conn, propertyType)
	if err != nil {
		return nil, fmt.Errorf("cannot get %s type atom", propertyType)
	}
	cookie := C.xcb_get_property(conn, 0, window, propertyAtom, typeAtom, 0, 64)
	reply := C.xcb_get_property_reply(conn, cookie, nil)
	if reply == nil {
		return nil, fmt.Errorf("cannot get %s property", property)
	}
	return reply, nil
}

func getFocusedWindow(conn *C.xcb_connection_t, root C.xcb_window_t) (C.xcb_window_t, error) {
	reply, err := getWindowPropertyReply(conn, root, "_NET_ACTIVE_WINDOW", "WINDOW")
	if err != nil {
		return 0, err
	}
	valuePointer := C.xcb_get_property_value(reply)
	window := *(*C.xcb_window_t)(valuePointer)
	C.free(unsafe.Pointer(reply))
	return window, nil
}

func getWindowStringProperty(
	conn *C.xcb_connection_t,
	window C.xcb_window_t,
	property string,
	propertyType string,
) (string, error) {
	reply, err := getWindowPropertyReply(conn, window, property, propertyType)
	if err != nil {
		return "", err
	}
	length := C.xcb_get_property_value_length(reply)
	value := C.GoStringN((*C.char)(C.xcb_get_property_value(reply)), length)
	C.free(unsafe.Pointer(reply))
	return value, nil
}

var nameAtoms = [...][2]string{
	{"_NET_WM_NAME", "UTF8_STRING"},
	{"WM_NAME", "UTF8_STRING"},
	{"WM_NAME", "STRING"},
}

func getWindowName(conn *C.xcb_connection_t, window C.xcb_window_t) (string, error) {
	for _, atoms := range nameAtoms {
		name, err := getWindowStringProperty(conn, window, atoms[0], atoms[1])
		if err != nil {
			return "", err
		}
		if name != "" {
			return name, nil
		}
	}
	return "", nil
}

func getWindowClass(conn *C.xcb_connection_t, window C.xcb_window_t) (string, string, error) {
	rawClass, err := getWindowStringProperty(conn, window, "WM_CLASS", "STRING")
	if err != nil {
		return "", "", errors.New("cannot find class")
	}
	parts := strings.Split(rawClass, "\000")
	// "instance\0class\0" -> ["instance", "class", ""]
	if len(parts) != 3 {
		return "", "", fmt.Errorf("cannot parse class: %s", rawClass)
	}
	return parts[0], parts[1], nil
}

func main() {
	conn, err := connect()
	if err != nil {
		die(err, conn)
	}
	screen, err := getScreen(conn)
	if err != nil {
		die(err, conn)
	}
	window, err := getFocusedWindow(conn, screen.root)
	if err != nil {
		die(err, conn)
	}
	name, err := getWindowName(conn, window)
	if err != nil {
		die(err, conn)
	}
	instance, class, err := getWindowClass(conn, window)
	if err != nil {
		die(err, conn)
	}
	fmt.Printf("id: %d\n", window)
	fmt.Printf("name: %s\n", name)
	fmt.Printf("instance: %s\n", instance)
	fmt.Printf("class: %s\n", class)
}
