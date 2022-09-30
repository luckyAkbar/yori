SHELL:=/bin/bash

ifdef test_run
	TEST_ARGS := -run $(test_run)
endif

check-modd-exists:
	@modd --version > /dev/null

run: check-modd-exists
	@modd -f ./.modd/server.modd.conf