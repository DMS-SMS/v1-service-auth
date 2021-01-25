
.PHONY: build
build: proto
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service *.go

.PHONY: image
image:
	docker build . -t dms-sms-service-auth:${VERSION}

.PHONY: upload
upload:
	docker tag dms-sms-service-auth:${VERSION} jinhong0719/dms-sms-service-auth:${VERSION}.RELEASE
	docker push jinhong0719/dms-sms-service-auth:${VERSION}.RELEASE

.PHONY: pull
pull:
	docker pull jinhong0719/dms-sms-service-auth:${VERSION}.RELEASE

.PHONY: run
run:
	docker-compose -f ./docker-compose.yml up -d

.PHONY: deploy
deploy:
	envsubst < ./service-auth-deployment.yaml | kubectl apply -f -

.PHONY: stack
stack:
	env VERSION=${VERSION} docker stack deploy -c docker-compose.yml DSM_SMS
