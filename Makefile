.PHONY: docker-build
docker-build:
	bash scripts/docker-build.sh

.PHONY: statik
statik:
	statik -src=configs -include=*.yaml -dest=configs

.PHONY: generate-key-pair
generate-key-pair:
	openssl genrsa 4096 > private.key
	openssl rsa -pubout < private.key > public.key

.PHONY: local-run
local-run: statik generate-key-pair
	export PRIVATE_KEY=`cat private.key`; go run ./cmd/server/.