.PHONY: test release install

test:
	go test -i
	go test -v

release: test
	godocdown -signature . > README.markdown

install: test
	go install .
