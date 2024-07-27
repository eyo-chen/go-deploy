clean:
	docker-compose down

remove-image:
	docker rmi go-deploy-image || true

build:
	docker-compose build --force-rm

start:
	docker-compose up -d

run: clean remove-image build start
