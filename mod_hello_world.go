// This file implements a simple FreeSWITCH module in Go.
// The module logs messages at different stages of its lifecycle (Load, Runtime, Shutdown)
// and provides a basic API command handler.
//
// Interaction with FreeSWITCH:
// - The module uses CGo to interface with FreeSWITCH's C API.
// - Exported C functions (_ModuleLoad, _ModuleRuntime, _ModuleShutdown, _ModuleApiHandler)
//   serve as entry points for FreeSWITCH to call into the Go module.
// - The `Stream` struct wraps the FreeSWITCH stream handle for writing messages to the console.
// - The `_Log` struct provides logging functionality via FreeSWITCH's logging API.
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

// Stream is a wrapper around the FreeSWITCH stream handle.
// It is used to send messages to the FreeSWITCH console.
type Stream struct {
	c_stream *C.switch_stream_handle_t
}

// Write sends a message to the FreeSWITCH console.
func (s Stream) Write(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	c_msg := C.CString(msg)
	defer C.free(unsafe.Pointer(c_msg))
	C._stream_write_function(s.c_stream, c_msg)
}

// Session is a placeholder for future session-related functionality.
type Session struct{}

// _Log is used for logging messages to the FreeSWITCH log.
type _Log struct{}

var Log = _Log{}

// Notice logs messages with the NOTICE severity level.
func (_ _Log) Notice(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	c_msg := C.CString(msg)
	defer C.free(unsafe.Pointer(c_msg))

	C._log_on_channel(C.SWITCH_LOG_NOTICE, c_msg)
}

// _ModuleLoad is called by FreeSWITCH when the module is loaded.
// It receives a pointer to the module interface, which can be used to register APIs, applications, etc.
// It should return SWITCH_STATUS_SUCCESS on success, or an error status on failure.
// Currently, due to a Go bug (https://github.com/golang/go/issues/11100), Go shared libraries cannot be unloaded,
// so we return SWITCH_STATUS_NOUNLOAD to prevent FreeSWITCH from attempting to unload the module.
//export _ModuleLoad
func _ModuleLoad(module_interface *C.switch_loadable_module_interface_t) C.switch_status_t {
	Load()
	/*
		BUG: Go shared library cannot be unloaded.
		see: https://github.com/golang/go/issues/11100
	*/
	return C.SWITCH_STATUS_NOUNLOAD
}

// _ModuleRuntime is called by FreeSWITCH after the module has been loaded and is ready to run.
// This function is typically used for long-running tasks or event loops.
// It should return SWITCH_STATUS_TERM when the module is finished running.
//export _ModuleRuntime
func _ModuleRuntime() C.switch_status_t {
	Runtime()
	return C.SWITCH_STATUS_TERM
}

// _ModuleShutdown is called by FreeSWITCH when the module is being unloaded.
// This function should be used to clean up any resources used by the module.
// It should return SWITCH_STATUS_SUCCESS.
//export _ModuleShutdown
func _ModuleShutdown() C.switch_status_t {
	Shutdown()
	return C.SWITCH_STATUS_SUCCESS
}

// _ModuleApiHandler is called by FreeSWITCH when an API command registered by this module is executed.
// c_cmd: The API command string.
// c_session: A pointer to the FreeSWITCH session, if the command was executed from within a call. Otherwise, it's NULL.
// c_stream: A pointer to the FreeSWITCH stream handle, used for sending output back to the caller.
// It should return SWITCH_STATUS_SUCCESS on success, or an error status on failure.
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

// Load is called when the module is loaded by FreeSWITCH.
// It's a good place to perform initialization tasks,
// such as registering API commands or event handlers.
func Load() {
	Log.Notice("Hello World!\n")
}

// Runtime is called after the module has been loaded and is ready to run.
// This function can be used for long-running tasks or event loops.
// In this example, it simply logs a message.
func Runtime() {
	Log.Notice("Ruling the World!\n")
}

// Shutdown is called when the module is being unloaded by FreeSWITCH.
// It's a good place to perform cleanup tasks,
// such as releasing resources or unregistering handlers.
func Shutdown() {
	Log.Notice("Good bye World!\n")
}

// ApiHandler is called when an API command registered by this module is executed.
// cmd: The API command string.
// session: The FreeSWITCH session, if the command was executed from within a call.
// stream: The FreeSWITCH stream handle, used for sending output back to the caller.
func ApiHandler(cmd string, session Session, stream Stream) {
	stream.Write("Hello from api handler: %s\n", cmd)
}
