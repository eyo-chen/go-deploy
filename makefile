build:
	docker build . -t go-deploy-image

run:
	docker run --env-file .env -d --name go-deploy -p 8080:8080 go-deploy-image