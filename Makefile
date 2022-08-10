SHELL := /bin/bash

.PHONY: shared
shared:
	@echo "Copying shared files"
	cp ./shared/* ./packages/user/auth/
	cp ./shared/* ./packages/user/migrate/
	cp ./shared/* ./packages/user/register/
	cp ./shared/* ./packages/gallery/single/
	cp ./shared/* ./packages/gallery/list/
