.PHONY: run test ext

build:
	scripts/build.sh

run:
	@scripts/run_app.sh $(ARGS)

test:
	scripts/test.sh

ext:
	scripts/make_extension_list.sh

# if target 'run', treat remaining make targets as run args (swallow them)
ifeq (run,$(firstword $(MAKECMDGOALS)))
ARGS := $(filter-out run --,$(MAKECMDGOALS))
# Create no-op targets so make does not error on them
$(foreach a,$(ARGS),$(eval $(a):;@true))
endif
