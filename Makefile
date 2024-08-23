zrun: build
	@./bin/app

build: 
	@go build -o bin/app .

css: 
	npx tailwindcss -i views/css/app.css -o api/public/css/styles.css --watch

templ:
	templ generate --watch --proxy=http://localhost:3000

dev:
	$(MAKE) css &  # Start CSS processing in the background
	air &  # Start air in the background
	$(MAKE) templ &  # Start templ in the background
	@wait  # Wait for all background processes to finish
