.PHONY: models generate

# ==============================================================================
# Swagger Models
# 	find ./internal/api/ -type f -not -name '*_test.go' -delete
models:
	$(call print-target)
	swagger generate model -f docs/basket_service.yml -m internal/api

generate: models