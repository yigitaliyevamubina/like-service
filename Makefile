CURRENT_DIR =$(shell pwd)
DB_URL := "postgres://postgres:mubina2007@localhost:5432/likedb?sslmode=disable"

proto-gen:
	chmod +x ./scripts/gen-proto.sh
	./scripts/gen-proto.sh

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-file:
	migrate create -ext sql -dir migrations/ create_commentlikes_table
