.PHONY: dev build start test clean migrate seed docker-up docker-down install

install:
	npm install

dev:
	npm run dev

build:
	npm run build

start: build
	npm start

migrate:
	npm run migrate

seed:
	npm run seed

test:
	npm test

clean:
	rm -rf dist node_modules

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

lint:
	npm run lint
