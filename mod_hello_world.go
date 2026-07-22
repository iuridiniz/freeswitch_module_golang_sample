// This file implements a simple FreeSWITCH module in Go.
// The module logs messages at different stages of its lifecycle (Load, Runtime, Shutdown)
// and provides a basic API command handler.
//
// All FreeSWITCH↔Go plumbing (cgo, //export entry points, C shims) lives in the
// freeswitch package (see freeswitch/freeswitch.go); this file contains only module logic and has no
// dependency on cgo. To build your own module from this sample, edit this file.
package main

import "example.com/freeswitch_mod_hello_world/freeswitch"

// helloModule implements freeswitch.Module and holds this sample's module logic.
type helloModule struct{}

func init() {
	freeswitch.Register(helloModule{})
}

func main() {} //No Op

// Load is called when the module is loaded by FreeSWITCH.
// It's a good place to perform initialization tasks,
// such as registering API commands or event handlers.
func (helloModule) Load() {
	freeswitch.Log.Notice("Hello World!\n")
}

// Runtime is called after the module has been loaded and is ready to run.
// This function can be used for long-running tasks or event loops.
// In this example, it simply logs a message.
func (helloModule) Runtime() {
	freeswitch.Log.Notice("Ruling the World!\n")
}

// Shutdown is called when the module is being unloaded by FreeSWITCH.
// It's a good place to perform cleanup tasks,
// such as releasing resources or unregistering handlers.
func (helloModule) Shutdown() {
	freeswitch.Log.Notice("Good bye World!\n")
}

// ApiHandler is called when an API command registered by this module is executed.
// cmd: The API command string.
// session: The FreeSWITCH session, if the command was executed from within a call.
// stream: The FreeSWITCH stream handle, used for sending output back to the caller.
func (helloModule) ApiHandler(cmd string, session freeswitch.Session, stream freeswitch.Stream) {
	stream.Write("Hello from api handler: %s\n", cmd)
}
