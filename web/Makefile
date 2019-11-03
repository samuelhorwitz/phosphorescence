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
	docker build -t phosphor-static --target static --build-arg EOS_ORIGIN=https://eoserigeneia.run --build-arg API_ORIGIN=https://api.phosphor.me --build-arg SCRIPTS_ORIGIN=https://scripts.phosphor.me --build-arg STATIC_ORIGIN=https://static.phosphor.me .
	docker build -t ${PHOSPHOR_LATEST} --build-arg EOS_ORIGIN=https://eoserigeneia.run --build-arg API_ORIGIN=https://api.phosphor.me --build-arg SCRIPTS_ORIGIN=https://scripts.phosphor.me --build-arg STATIC_ORIGIN=https://static.phosphor.me .

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

apply:
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml apply -f k8s.yaml

deploy: apply
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods -l type=server

deploy-phosphor: deploy-phosphor-static deploy-phosphor-k8s

deploy-phosphor-k8s: apply
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods -l app=phosphor

deploy-phosphor-static:
	rm -rf ./.production
	-docker rm phosphor-static
	docker create --name phosphor-static phosphor-static
	docker cp phosphor-static:/app/public ./.production
	# js, css, html, svg, ttf, xml, json, txt
	find ./.production -type f -name '*.js' -exec gzip -9 -- {} \;
	find ./.production -type f -name '*.css' -exec gzip -9 -- {} \;
	find ./.production -type f -name '*.html' -exec gzip -9 -- {} \;
	find ./.production -type f -name '*.svg' -exec gzip -9 -- {} \;
	find ./.production -type f -name '*.ttf' -exec gzip -9 -- {} \;
	find ./.production -type f -name '*.xml' -exec gzip -9 -- {} \;
	find ./.production -type f -name '*.json' -exec gzip -9 -- {} \;
	find ./.production -type f -name '*.txt' -exec gzip -9 -- {} \;
	find ./.production -type f -name '*.gz' -exec sh -c 'mv "$$0" "$${0%.gz}"' {} \;
	aws s3 sync ./.production s3://phosphorescence-static/ --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --endpoint=https://nyc3.digitaloceanspaces.com --exclude ".DS_Store" --exclude ".nojekyll" --exclude "*.js" --exclude "*.css" --exclude "*.html" --exclude "*.svg" --exclude "*.ttf" --exclude "*.xml" --exclude "*.json" --exclude "*.txt"
	aws s3 sync ./.production s3://phosphorescence-static/ --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --endpoint=https://nyc3.digitaloceanspaces.com --exclude "*" --include "*.js" --include "*.css" --include "*.html" --include "*.svg" --include "*.ttf" --include "*.xml" --include "*.json" --include "*.txt" --content-encoding gzip

deploy-eos: apply
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods -l app=eos

deploy-api: apply
	kubectl --kubeconfig=/Users/${USER}/.kube/phosphor-kubeconfig.yaml --namespace phosphor delete pods -l app=api

postgres-local:
	docker run --rm --name phosphor-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=phosphor -d -p 5432:5432 -v ${PWD}/pgdata:/var/lib/postgresql/data samuelhorwitz/postgres-perl:11.5

redis-local:
	docker run -d -p 6379:6379 --name phosphor-redis redis

redis-cache-local:
	docker run -d -p 6389:6379 --name phosphor-cache redis
