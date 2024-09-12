.PHONY:buildup
buildup:
	make down
	docker-compose up --build

.PHONY:down
down:
	docker-compose down

.PHONY:up
up:
	make down
	docker-compose up	

.PHONY:up
upd:
	make down
	docker-compose up -d
