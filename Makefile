.SILENT:

PACKAGE = github.com/sudachen/smwlt/
TESTS = test0 testwlt testcli testnet
TESTDIR = .data/tests

mk-data-dir:
	mkdir -p $(TESTDIR)
	cp tests/accounts.json $(TESTDIR)/..

build:
	go build ./...

build-windows-tests: mk-data-dir
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc make build-cross-tests EXT=exe

build-osx-tests: mk-data-dir
	env GOOS=darwin GOARCH=amd64 make build-cross-tests EXT=osx

build-linux-tests: mk-data-dir
	make build-cross-tests EXT=test

build-cross-tests: mk-data-dir
	for i in $$(find $(TESTDIR) -name '*.$(EXT)'); do rm $$i; done
	for i in $(TESTS); do \
		cd tests/$$i; \
		go test -o ../../$(TESTDIR)/$$i.$(EXT) -c -covermode=atomic -coverpkg=../../...; \
		cd ../..; \
	done

run-linux-tests: build-linux-tests
	cd $(TESTDIR) && \
		for i in ./*.test; do \
			$$i -test.v=true -test.coverprofile=$$i.out > $$i.log; \
		done

run-windows-tests: build-windows-tests
	cd $(TESTDIR) && \
		for i in *.exe; \
			do wine $$i -test.v=true -test.coverprofile=$$i.out > $$i.log; \
		done

collect-tests:
	if [ -f $(TESTDIR)/c.out ]; then rm $(TESTDIR)/c.out; fi
	for i in $$(find $(TESTDIR) -name '*.test.out'); do tail -n +2 $$i >> $(TESTDIR)/c.out; done
	for i in $$(find $(TESTDIR) -name '*.exe.out'); do tail -n +2 $$i >> $(TESTDIR)/c.out; done
	for i in $$(find $(TESTDIR) -name '*.osx.out'); do tail -n +2 $$i >> $(TESTDIR)/c.out; done
	echo "mode: atomic" > c.out
	cat $(TESTDIR)/c.out | sort >> c.out
	sed -i -e '\:^$(PACKAGE)/tests:d' c.out
	sed -i -e 's:$(PACKAGE)::g' c.out
	awk '/\.go/{print "$(PACKAGE)"$$0}/^mode/{print $$0}' < c.out > gocov.txt

check-fail:
	cd $(TESTDIR) && if [ $$(cat *.log | grep FAIL | wc -l) -gt 0 ]; then \
		grep '.-- FAIL:' *.log; \
		exit 1; \
		fi

run-cover:
	go tool cover -html=gocov.txt

clean-tests:
	for i in $$(find $(TESTDIR) -name '*.out'); do rm $$i; done

run-cover-all: clean-tests run-linux-tests run-windows-tests collect-tests check-fail run-cover
run-cover-linux: clean-tests run-linux-tests collect-tests check-fail run-cover
run-all-tests: clean-tests run-linux-tests run-windows-tests check-fail
run-tests: clean-tests run-linux-tests check-fail

