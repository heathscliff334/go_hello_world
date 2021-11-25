# Initialization build stage
# FROM ubuntu
FROM golang:alpine

# To execute bash script from golang:alpine
RUN apk update && apk add git

# Arguments
# ARG DEFAULT_PORT
# ARG MYSQL_CONN_STRING

# Initialize environment variable manually
ENV PORT=8083
# ENV PORT=${DEFAULT_PORT}
ENV INSTANCE_ID = "Docker"

# To use /app become working directory for dockerfile statements
WORKDIR /app

# show current working directory
# RUN ls

# To copying from host to docker current working directory
#  This folder to /app
# COPY go.mod /app/go.mod
# COPY main.go /app/main.go
# OR
# COPY go.mod /go.mod
# COPY main.go /main.go
# OR
COPY . .

# To validate our source code, is all third parties installed or not, if not then automatically get the third parties
RUN go mod tidy
# To build golang's code to binary executable
RUN go build -o binary

# RUN .binary
ENTRYPOINT [ "./binary" ]

# to build in terminal
# docker build <current directory, or just use period "."> -t image-name:tag
# docker build <current directory, or just use period "."> -t image-name:tag --build-arg DEFAULT_PORT = .... > to build with argument
# run image:
# docker run -i -t image-name:tag
# -i = enable interactive mode, -t = enable tty
# -e = environment variable
# -p = expose port
# if using port env
# docker run -i -t -e PORT=.... -p <host>....:<container>.... image-name:tag

# To create container and start manually:
# docker container create -e PORT:.... -p ....:.... --name <container names (cannot duplicate)> image-name
# docker container start <ID>

# view list docker
# docker container ls (just for currently running container)
# docker container ls -a / --all (for all container)

# To remote container via SSH
# docker exec -it <Container ID> sh/bash/(binary name)