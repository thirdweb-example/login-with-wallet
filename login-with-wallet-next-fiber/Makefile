.PHONY: server web

buildserver:
	cd server && go build -o build/server main.go
runserver:
	cd server && ./build/server
server:
	make buildserver && make runserver

web:
	cd web && yarn dev
