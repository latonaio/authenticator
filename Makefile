.PHONY: docker-build
docker-build:
	bash scripts/docker-build.sh

.PHONY: statik
statik:
	statik -src=configs -include=*.yaml -dest=configs