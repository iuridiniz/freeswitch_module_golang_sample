// Package freeswitch contains all the FreeSWITCH↔Go plumbing (cgo, //export entry points,
// C shims) required to build this sample as a FreeSWITCH module. It is the only
// package in this project that uses cgo.
//
// Interaction with FreeSWITCH:
//   - The package uses CGo to interface with FreeSWITCH's C API.
//   - Exported C functions (_ModuleLoad, _ModuleRuntime, _ModuleShutdown, _ModuleApiHandler)
//     serve as entry points for FreeSWITCH to call into Go, via freeswitch.c.
//   - The `Stream` struct wraps the FreeSWITCH stream handle for writing messages to the console.
//   - The `Log` variable provides logging functionality via FreeSWITCH's logging API.
//
// Consumers (package main) implement the Module interface and call Register() from an
// init() function; Go guarantees all package init() functions run at dlopen time, before
// FreeSWITCH calls any exported entry point.
package freeswitch

/*
#include "freeswitch.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

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

// callerInfo returns the file, function name, and line number of the caller,
// skip frames up the stack (skip=1 targets the direct caller of callerInfo's caller).
func callerInfo(skip int) (file, fn string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", "", 0
	}
	if f := runtime.FuncForPC(pc); f != nil {
		fn = f.Name()
	}
	return file, fn, line
}

// Notice logs messages with the NOTICE severity level.
func (_ _Log) Notice(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	file, fn, line := callerInfo(2) // skip: callerInfo, Notice -> lands on the caller of Notice
	c_msg := C.CString(msg)
	c_file := C.CString(file)
	c_func := C.CString(fn)
	defer C.free(unsafe.Pointer(c_msg))
	defer C.free(unsafe.Pointer(c_file))
	defer C.free(unsafe.Pointer(c_func))

	C._log_on_channel(C.SWITCH_LOG_NOTICE, c_file, c_func, C.int(line), c_msg)
}

// Module is the interface a consumer (package main) must implement and register with
// Register() so the glue code can dispatch FreeSWITCH lifecycle/API events to it.
type Module interface {
	// Load is called by FreeSWITCH when the module is loaded.
	Load()
	// Runtime is called after the module has been loaded and is ready to run.
	Runtime()
	// Shutdown is called by FreeSWITCH when the module is being unloaded.
	Shutdown()
	// ApiHandler is called by FreeSWITCH when an API command registered by this
	// module is executed.
	ApiHandler(cmd string, session Session, stream Stream)
}

// module holds the Module implementation registered by package main.
var module Module

// Register associates a Module implementation with the glue code. It must be called
// from an init() function in package main so that it runs before FreeSWITCH invokes
// any of the exported entry points below.
func Register(m Module) {
	module = m
}

// _ModuleLoad is called by FreeSWITCH when the module is loaded.
// It receives a pointer to the module interface, which can be used to register APIs, applications, etc.
// It should return SWITCH_STATUS_SUCCESS on success, or an error status on failure.
// Currently, due to a Go bug (https://github.com/golang/go/issues/11100), Go shared libraries cannot be unloaded,
// so we return SWITCH_STATUS_NOUNLOAD to prevent FreeSWITCH from attempting to unload the module.
//
//export _ModuleLoad
func _ModuleLoad(module_interface *C.switch_loadable_module_interface_t) C.switch_status_t {
	if module != nil {
		module.Load()
	}
	/*
		BUG: Go shared library cannot be unloaded.
		see: https://github.com/golang/go/issues/11100
	*/
	return C.SWITCH_STATUS_NOUNLOAD
}

// _ModuleRuntime is called by FreeSWITCH after the module has been loaded and is ready to run.
// This function is typically used for long-running tasks or event loops.
// It should return SWITCH_STATUS_TERM when the module is finished running.
//
//export _ModuleRuntime
func _ModuleRuntime() C.switch_status_t {
	if module != nil {
		module.Runtime()
	}
	return C.SWITCH_STATUS_TERM
}

// _ModuleShutdown is called by FreeSWITCH when the module is being unloaded.
// This function should be used to clean up any resources used by the module.
// It should return SWITCH_STATUS_SUCCESS.
//
//export _ModuleShutdown
func _ModuleShutdown() C.switch_status_t {
	if module != nil {
		module.Shutdown()
	}
	return C.SWITCH_STATUS_SUCCESS
}

// _ModuleApiHandler is called by FreeSWITCH when an API command registered by this module is executed.
// c_cmd: The API command string.
// c_session: A pointer to the FreeSWITCH session, if the command was executed from within a call. Otherwise, it's NULL.
// c_stream: A pointer to the FreeSWITCH stream handle, used for sending output back to the caller.
// It should return SWITCH_STATUS_SUCCESS on success, or an error status on failure.
//
//export _ModuleApiHandler
func _ModuleApiHandler(c_cmd *C.cchar_t, c_session *C.switch_core_session_t, c_stream *C.switch_stream_handle_t) C.switch_status_t {
	if module == nil {
		return C.SWITCH_STATUS_SUCCESS
	}

	cmd := C.GoString(c_cmd)
	stream := Stream{c_stream}
	session := Session{}
	module.ApiHandler(cmd, session, stream)

	return C.SWITCH_STATUS_SUCCESS
}
