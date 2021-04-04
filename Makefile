# (C) Copyright IBM Corporation 2021.
# LICENSE: GPL-3.0 https://opensource.org/licenses/GPL-3.0

.PHONY: docker-push			##pushes the operator to a docker registry. in order to run it need to add <IMAGE> var as an argument.
docker-push: build
	@operator-sdk build ${IMAGE}
	@docker push ${IMAGE}

.PHONY: install			##install the oeprator in the specified <NAMESPACE>.
install:
	@build/install.sh ${NAMESPACE}


.PHONY: clean			##cleans all objects from kubernetes in the specified <NAMESPACE>.
clean:
	@build/clean.sh ${NAMESPACE}

.PHONY: help				##show this help message
help:
	@echo "usage: make [target]\n"; echo "options:"; \fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//' | sed 's/.PHONY:*//' | sed -e 's/^/  /'; echo "";
