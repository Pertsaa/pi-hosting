build:
	go build -o dist/pi-hosting

start:
	./dist/pi-hosting -t=<DISCORD_BOT_TOKEN> -u=<IP_API_URL>
