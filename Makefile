all: 

build:
	@echo "Building locally..."
	go run .

# having the - before docker stop/rm will ignore any errors, like if the container is not running
# the -p flag of docker build maps the host port to the container port
# the -d flag of docker run runs the container in detached mode. i.e., in the background
rebuild:
	@echo "Rebuilding the Docker container..."
	docker build -t phishhook:alpha . 
	-docker stop phishhook-backend
	-docker rm phishhook-backend
	docker run -d -p 8080:8081 --name phishhook-backend phishhook:alpha