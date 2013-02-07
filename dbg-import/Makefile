.PHONY: build release install clean

build:
	go build

release:
	$(MAKE) -C .. $@

install: build
	go install

clean:
	rm -f dbg-import
