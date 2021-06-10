VERSION := $(shell cat ./VERSION)

release:
	git tag -a $(VERSION) -m "release $(VERSION)" || true
	git push origin master --tags
.PHONY: release

.PHONY: golinux
