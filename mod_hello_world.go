package main

/*
#include "fs.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

/****************************************/
/* START OF GLUE CODE */
/****************************************/

type Stream struct {
	c_stream *C.switch_stream_handle_t
}

func (s Stream) Write(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	c_msg := C.CString(msg)
	defer C.free(unsafe.Pointer(c_msg))
	C._stream_write_function(s.c_stream, c_msg)
}

type Session struct{}

type _Log struct{}

var Log = _Log{}

func (_ _Log) Notice(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	c_msg := C.CString(msg)
	defer C.free(unsafe.Pointer(c_msg))

	C._log_on_channel(C.SWITCH_LOG_NOTICE, c_msg)
}

//export _ModuleLoad
func _ModuleLoad(module_interface *C.switch_loadable_module_interface_t) C.switch_status_t {
	Load()
	/*
		BUG: Go shared library cannot be unloaded.
		see: https://github.com/golang/go/issues/11100
	*/
	return C.SWITCH_STATUS_NOUNLOAD
}

//export _ModuleRuntime
func _ModuleRuntime() C.switch_status_t {
	Runtime()
	return C.SWITCH_STATUS_TERM
}

//export _ModuleShutdown
func _ModuleShutdown() C.switch_status_t {
	Shutdown()
	return C.SWITCH_STATUS_SUCCESS
}

//export _ModuleApiHandler
func _ModuleApiHandler(c_cmd *C.cchar_t, c_session *C.switch_core_session_t, c_stream *C.switch_stream_handle_t) C.switch_status_t {

	cmd := C.GoString(c_cmd)
	stream := Stream{c_stream}
	session := Session{}
	ApiHandler(cmd, session, stream)

	return C.SWITCH_STATUS_SUCCESS
}

/****************************************/
/* END OF GLUE CODE */
/****************************************/

func main() {} //No Op

func Load() {
	Log.Notice("Hello World!\n")
}

func Runtime() {
	Log.Notice("Ruling the World!\n")
}

func Shutdown() {
	Log.Notice("Good bye World!\n")
}

func ApiHandler(cmd string, session Session, stream Stream) {
	stream.Write("Hello from api handler: %s\n", cmd)
}
