run: build
	@./bin/app

build: 
	@go build -o bin/app .

css: 
	npx tailwindcss -i views/css/app.css -o api/public/css/styles.css --watch

templ:
	templ generate --watch --proxy=http://localhost:3000

dev: 
	air &
	@while ! nc -z localhost 3000; do sleep 1; done
	$(MAKE) css & $(MAKE) templ
	@wait
