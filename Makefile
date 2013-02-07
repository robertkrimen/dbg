.PHONY: test release install clean

test:
	go test -i
	go test

release: test
	(cd dbg-import && godocdown -signature . > README.markdown) || false
	godocdown -signature . > README.markdown

install: test
	go install
	$(MAKE) -C dbg-import $@

clean:
	$(MAKE) -C dbg-import $@
