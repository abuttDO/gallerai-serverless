SHELL := /bin/bash

.PHONY: shared
shared:
	@echo "Copying shared files"
	cp ./shared/* ./packages/user/auth/
	cp ./shared/* ./packages/user/migrate/
	cp ./shared/* ./packages/user/register/
	cp ./shared/* ./packages/user/forgot-password/
	cp ./shared/* ./packages/gallery/single/
	cp ./shared/* ./packages/gallery/list/
	cp ./shared/* ./packages/gallery/reaction/
