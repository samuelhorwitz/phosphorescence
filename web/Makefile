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

push: push-phosphor push-eos push-jobs push-api

push-phosphor:
	docker push ${PHOSPHOR_LATEST}

push-eos:
	docker push ${EOS_LATEST}

push-jobs:
	docker push ${JOBS_LATEST}

push-api:
	docker push ${API_LATEST}

deploy:
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml apply -f k8s.yaml
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods --all

postgres-local:
	docker run --rm --name phosphor-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=phosphor -d -p 5432:5432 -v ${PWD}/pgdata:/var/lib/postgresql/data postgres
