package x11

// #cgo pkg-config: xcb
// #include <stdlib.h>
// #include <string.h>
// #include <xcb/xcb.h>
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// Connection ...
type Connection struct {
	connection *C.xcb_connection_t
	screen     *C.xcb_screen_t
}

// WindowID ...
type WindowID uint32

// Window ...
type Window struct {
	window     C.xcb_window_t
	connection *C.xcb_connection_t
}

// Connect ...
func Connect() (*Connection, error) {
	connection := C.xcb_connect(nil, nil)
	if C.xcb_connection_has_error(connection) != 0 {
		return nil, errors.New("cannot connect to X server")
	}
	screen := C.xcb_setup_roots_iterator(C.xcb_get_setup(connection)).data
	if screen == nil {
		return nil, errors.New("cannot get screen")
	}
	return &Connection{connection, screen}, nil
}

// Disconnect ...
func (c *Connection) Disconnect() {
	C.xcb_disconnect(c.connection)
}

// RootWindow ...
func (c *Connection) RootWindow() *Window {
	root := c.screen.root
	return &Window{root, c.connection}
}

// FocusedWindow ...
func (c *Connection) FocusedWindow() (*Window, error) {
	root := c.RootWindow()
	reply, err := getWindowPropertyReply(c.connection, root.window, "_NET_ACTIVE_WINDOW", "WINDOW")
	if err != nil {
		return nil, err
	}
	valuePointer := C.xcb_get_property_value(reply)
	window := *(*C.xcb_window_t)(valuePointer)
	C.free(unsafe.Pointer(reply))
	return &Window{window, c.connection}, nil
}

var nameAtoms = [...][2]string{
	{"_NET_WM_NAME", "UTF8_STRING"},
	{"WM_NAME", "UTF8_STRING"},
	{"WM_NAME", "STRING"},
}

// ID ...
func (w *Window) ID() WindowID {
	return WindowID(w.window)
}

// Name ...
func (w *Window) Name() (string, error) {
	for _, atoms := range nameAtoms {
		name, err := getWindowStringProperty(w.connection, w.window, atoms[0], atoms[1])
		if err != nil {
			return "", err
		}
		if name != "" {
			return name, nil
		}
	}
	return "", nil
}

// Class ...
func (w *Window) Class() (string, string, error) {
	rawClass, err := getWindowStringProperty(w.connection, w.window, "WM_CLASS", "STRING")
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
