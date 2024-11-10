REPO ?= neovim/neovim
run:
	scripts/run_app.sh -g "$(REPO)"

test:
	scripts/test.sh
