deploy:
	helm --kubeconfig ~/.kube/config --kube-context docker-desktop dependency build \
		./deploy/helm/postgresql
	helm --kubeconfig ~/.kube/config --kube-context docker-desktop upgrade \
		--install postgresql \
		--values=./deploy/helm/postgresql/values.yaml \
		--values=./deploy/helm/postgresql/values_$(env).yaml \
		./deploy/helm/postgresql

.PHONY: deploy