tailwind-watch:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

tailwind-build:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --minify

templ-generate:
	templ generate

templ-watch:
	templ generate --watch

up:
	@docker-compose up -d

dev: up
	@go build -o ./tmp/$(APP_NAME) ./cmd/$(APP_NAME)/main.go && air