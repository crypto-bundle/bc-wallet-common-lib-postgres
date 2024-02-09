
default: deploy_postgresql

deploy_postgresql:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval context=$(or $(context),k0s-dev-cluster))
	$(eval platform=$(or $(platform),linux/amd64))

	helm --kube-context $(context) upgrade \
		--install postgresql \
		--values=./deploy/helm/postgresql/values.yaml \
		--values=./deploy/helm/postgresql/values_$(env).yaml \
		./deploy/helm/postgresql

destroy_postgresql:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval context=$(or $(context),k0s-dev-cluster))
	$(eval platform=$(or $(platform),linux/amd64))

	helm --kube-context $(context) uninstall postgresql

.PHONY: deploy_postgresql