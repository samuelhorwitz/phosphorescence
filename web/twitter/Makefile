build:
	npm run clean
	npm run build
	cp -R dist distgz
	gzip ./distgz/*
	find ./distgz/*.gz -type f | while read f; do mv "$$f" "$${f%.gz}"; done

push:
	aws s3 sync ./distgz s3://phosphorescence-twitter/ --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --endpoint=https://nyc3.digitaloceanspaces.com --content-encoding gzip