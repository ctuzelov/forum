build:
	docker build -t forum .
run:
	docker run -dp 8080:8080 --rm --name container1 forum
stop:
	docker stop container1