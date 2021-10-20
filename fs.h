#ifndef _FS_H
#define _FS_H

#include <switch.h>

/* https://stackoverflow.com/a/32940436/1522342 */

typedef const char cchar_t;

void _log_on_channel(switch_log_level_t level, char *msg);
void _stream_write_function(switch_stream_handle_t *stream, char *msg);

#endif