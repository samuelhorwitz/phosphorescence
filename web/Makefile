REGISTRY        := docker.phosphor.me
PHOSPHOR_NAME   := ${REGISTRY}/samuelhorwitz/phosphor
EOS_NAME        := ${REGISTRY}/samuelhorwitz/eos
JOBS_NAME       := ${REGISTRY}/samuelhorwitz/phosphor-jobs
API_NAME        := ${REGISTRY}/samuelhorwitz/phosphor-api
PHOSPHOR_LATEST := ${PHOSPHOR_NAME}:latest
EOS_LATEST      := ${EOS_NAME}:latest
JOBS_LATEST     := ${JOBS_NAME}:latest
API_LATEST      := ${API_NAME}:latest

build: build-phosphor build-eos build-jobs build-api

build-phosphor:
	docker build -t ${PHOSPHOR_LATEST} --build-arg EOS_ORIGIN=https://eoserigeneia.run --build-arg API_ORIGIN=https://api.phosphor.me --build-arg SCRIPTS_ORIGIN=https://scripts.phosphor.me .

build-eos:
	docker build -t ${EOS_LATEST} -f ./eos/Dockerfile --build-arg PHOSPHOR_ORIGIN=https://phosphor.me .

build-jobs:
	docker build -t ${JOBS_LATEST} -f ./jobs/Dockerfile .

build-api:
	docker build -t ${API_LATEST} -f ./api/Dockerfile .

build-twitter:
	$(MAKE) -C twitter build

push: push-phosphor push-eos push-jobs push-api

push-phosphor:
	docker push ${PHOSPHOR_LATEST}

push-eos:
	docker push ${EOS_LATEST}

push-jobs:
	docker push ${JOBS_LATEST}

push-api:
	docker push ${API_LATEST}

push-twitter:
	${MAKE} -C twitter push

apply:
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml apply -f k8s.yaml

deploy: apply
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods -l type=server

deploy-phosphor: apply
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods -l app=phosphor

deploy-eos: apply
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods -l app=eos

deploy-api: apply
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods -l app=api

deploy-jobs: apply

postgres-local:
	docker run --rm --name phosphor-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=phosphor -d -p 5432:5432 -v ${PWD}/pgdata:/var/lib/postgresql/data samuelhorwitz/postgres-perl:11.5

redis-local:
	-docker rm phosphor-redis
	docker run -d -p 6379:6379 --name phosphor-redis redis

redis-cache-local:
	-docker rm phosphor-cache
	docker run -d -p 6389:6379 --name phosphor-cache redis
