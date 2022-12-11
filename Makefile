SHELL:=/bin/bash

ifdef test_run
	TEST_ARGS := -run $(test_run)
endif

check-modd-exists:
	@modd --version > /dev/null

run: check-modd-exists
	@modd -f ./.modd/server.modd.conf



migrate_up=go run main.go migrate --direction=up --step=0

migrate:
	@if [ "$(DIRECTION)" = "" ] || [ "$(STEP)" = "" ]; then\
    	$(migrate_up);\
	else\
		go run main.go migrate --direction=$(DIRECTION) --step=$(STEP);\
    fi