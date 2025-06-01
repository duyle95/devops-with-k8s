docker-build-push:
	docker build ./$(name) -t duysmartum/$(name):latest
	docker push duysmartum/$(name):latest
.PHONY: docker-build-push