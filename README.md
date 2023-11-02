# go-backend

backend server to support the phishhook android application

# Docker

here: https://www.docker.com/blog/developing-go-apps-docker/
`docker build --rm -t phishhook:alpha .`

`docker run -d -p 8080:8081 --name phishhook-backend phishhook:alpha`

# API Keys

Hunter
`phishhookxnGvHDpBDrQ8eynzXypw181R4PptZXfnpUZkx0r96Hy5DgcIudsmrogfoc4wf4ugwB9kakwxVxBYDTPVJCQ5bGKfVpRMtNLtDjIxgh8QVGBT4nC6mDUbQ2U`

Kab
`phishhook7A4t8BrNThTSeJZBZYwZ9OAb0HpkWldiVQcHF9cSNJk50vdD4VgZ7NJvXdALVXkElKv9rcvK91pCkPUMLBjlswxQLbK2IA7GHcGHFBexXX1xvdrNfbuk779`

Lucas
`phishhookRyJHCenIz97Q5LIDPmHhDyg9eddxaBO29omDuzM1D5BsDRKH5mo3j8pmBehoO2Roj0Z4zWuDHlNW4AJVrSnLZF6lUravmyje13YB1LBriXHxYlxLUDYeXmV`

test

```
curl -H "X-API-KEY: phishhookRyJHCenIz97Q5LIDPmHhDyg9eddxaBO29omDuzM1D5BsDRKH5mo3j8pmBehoO2Roj0Z4zWuDHlNW4AJVrSnLZF6lUravmyje13YB1LBriXHxYlxLUDYeXmV"  http://localhost:8081/albums
```
