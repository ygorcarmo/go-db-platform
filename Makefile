run: build
	@./bin/app

build: 
	@go build -o bin/app .

css: 
	npx @tailwindcss/cli -i views/css/app.css -o api/public/css/styles.css --watch

templ:
	templ generate --watch

dev:
	(trap 'kill 0' SIGINT; $(MAKE) css & air &  $(MAKE) templ )  
