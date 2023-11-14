# Phishhook Backend Application

Backend server to support the phishhook android application.

## Docker

To deploy locally to Docker, you can run `make rebuild` and check your local docker app for the container. It will expose the application on port `8080`.

## Tests

### Locally

```
curl -H "X-API-KEY: {API_KEY}"  http://localhost:8080/users
```

### AWS

```
curl -H'X-API-KEY: {API_KEY}' http://ec2-18-224-251-242.us-east-2.compute.amazonaws.com:8080/users
```

## Deployment to ECS from ECR

This repository is served to the outside world using the Elastic Container Service from Amazon AWS. To make changes and deploy them to ECS, first we must update our Elastic Container Registry image by running `make push`.

Then, login to the AWS console, and update the Task Definition with the URI of the image. Save. Stop the currently running task and deploy the new Task Definition Revision. Once running, the container is updated!
