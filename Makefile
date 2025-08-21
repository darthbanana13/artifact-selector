REPO ?= neovim/neovim
run:
	scripts/run_app.sh -g "$(REPO)" -e "$(EXT)" -a "$(ARCH)" -o "$(OS)" -v "$(VERSION)"

test:
	scripts/test.sh

ext:
	scripts/make_extension_list.sh
