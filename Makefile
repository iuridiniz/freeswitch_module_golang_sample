.PHONY: clean install

GO_BINARY=go
FREESWITCH_DIR=/usr/local/freeswitch/
FREESWITCH_LIB_DIR=$(FREESWITCH_DIR)/lib
FREESWITCH_INCLUDE_DIR=$(FREESWITCH_DIR)/include/freeswitch
FREESWITCH_MOD_DIR=$(FREESWITCH_DIR)/mod

CFLAGS=-I$(FREESWITCH_INCLUDE_DIR)
LDFLAGS=-L$(FREESWITCH_LIB_DIR) -lfreeswitch -Wl,-rpath=$(FREESWITCH_LIB_DIR)

mod_hello_world.so: $(wildcard *.go *.c *.h) go.mod
	CGO_CFLAGS="$(CGO_CFLAGS) $(CFLAGS)" \
	CGO_LDFLAGS="$(CGO_LDFLAGS) $(LDFLAGS)" \
		$(GO_BINARY) build -buildmode=c-shared -o $@

clean:
	rm -f mod_hello_world.so mod_hello_world.h

install: mod_hello_world.so
	install -m 0755 -t $(FREESWITCH_MOD_DIR) $<