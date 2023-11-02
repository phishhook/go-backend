# specify the node base image with your desired version node:<version>
FROM golang:1.21.3-bullseye

# Set the working directory to /app
WORKDIR /app

# Add the env file
COPY .env ./

# Effectively tracks changes within your go.mod file
COPY go.mod go.sum ./

# Install Go dependencies
RUN go mod download

# Copies your source code into the app directory
COPY main.go .

RUN go build -o /phishhook-backend

# Tells Docker which network port your container listens on
EXPOSE 8080

# specify the command to run on container start
CMD [ "/phishhook-backend" ]