all: 

build:
	@echo "Building locally..."
	go run .

docker:
	@echo "Building Docker container..."
	docker build -t phishhook .

# having the - before docker stop/rm will ignore any errors, like if the container is not running
# the -p flag of docker build maps the host port to the container port
# the -d flag of docker run runs the container in detached mode. i.e., in the background
rebuild:
	@echo "Rebuilding the Docker container..."
	docker build -t phishhook:alpha . 
	-docker stop phishhook-backend
	-docker rm phishhook-backend
	docker run -d -p 8080:8081 --name phishhook-backend phishhook:alpha

push:
	@echo "updating Docker Hub with new image, then pushing to AWS ECR..."
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/n1f8c0d4
	docker build -t phishhook-backend .
	docker tag phishhook-backend:latest public.ecr.aws/n1f8c0d4/phishhook-backend:latest
	docker push public.ecr.aws/n1f8c0d4/phishhook-backend:latest

# run when you are done with your working session to save computer resources
pause:
	docker pause 1f99d198ceca