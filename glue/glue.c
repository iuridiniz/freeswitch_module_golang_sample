#include <switch.h>
#include "glue.h"
#include "_cgo_export.h"

/* these are declared in _cgo_export.h and defined in glue.go */
// switch_status_t _ModuleRuntime();
// switch_status_t _ModuleShutdown();
// switch_status_t _ModuleLoad();
// switch_status_t _ModuleApiHandler();

SWITCH_MODULE_LOAD_FUNCTION(_wrap_load);
SWITCH_MODULE_SHUTDOWN_FUNCTION(_wrap_shutdown);
SWITCH_MODULE_RUNTIME_FUNCTION(_wrap_runtime);

/// This macro needs to be called once in the module. It will generate the definitions
/// required to be loaded by FreeSWITCH. FS requires the exported table to have a name
/// of <filename>_module_interface. If your mod is called mod_foo, then the first param
/// to this macro must be mod_foo_module_interface.
SWITCH_MODULE_DEFINITION(mod_hello_world, _wrap_load, _wrap_shutdown, _wrap_runtime);

SWITCH_MODULE_LOAD_FUNCTION(_wrap_load)
{
    switch_api_interface_t *api_interface;

    *module_interface = switch_loadable_module_create_module_interface(pool, modname);

    /* TODO: implement a way to call this from go */
    SWITCH_ADD_API(api_interface, "hello", "Hello API", (_ModuleApiHandler), "hello syntax");

    return (_ModuleLoad)(*module_interface);
}

SWITCH_MODULE_SHUTDOWN_FUNCTION(_wrap_shutdown)
{
    return (_ModuleShutdown)();
}

SWITCH_MODULE_RUNTIME_FUNCTION(_wrap_runtime)
{
    return (_ModuleRuntime)();
}

/* Calling variadic C functions is not supported by cgo */
void _log_on_channel(switch_log_level_t level, const char *file, const char *func, int line, char *msg)
{
    switch_log_printf(SWITCH_CHANNEL_ID_LOG, file, func, line, NULL, level, "%s", msg);
}

/* Calling C function pointers is currently not supported by cgo */
void _stream_write_function(switch_stream_handle_t *stream, char *msg)
{
    stream->write_function(stream, "%s", msg);
}
