.PHONY: app
app:
	sudo docker-compose -f ./docker-compose-app.yml up -d --build


.PHONY: nt
nt:
	ghz --config ./infa/ghz_config.yaml localhost:50051


