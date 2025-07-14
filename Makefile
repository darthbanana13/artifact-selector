REPO ?= neovim/neovim
run:
	scripts/run_app.sh -g "$(REPO)" -e "$(EXT)" -a "$(ARCH)" -o "$(OS)"

test:
	scripts/test.sh
