.PHONY: test release 

export TERST_BASE=$(PWD)

test:
	go test -i && go test -v

release: test
	godocdown -signature . > README.markdown

