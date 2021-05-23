# docker build -t my-go-server .
# docker run -it my-go-server
#
# docker build -t my-go-client .
# docker run my-go-client
#
# https://dzone.com/articles/docker-to-docker-networking-between-tcp-client-and
# https://tutorialedge.net/golang/go-docker-tutorial/
# https://docs.docker.com/config/containers/container-networking/
# 
# docker build -t my-go-server .
# docker run --name go-server -it --rm --network="host" my-go-server
# Ta bort alla containers:
# docker container rm $(docker container ls -aq)

## We specify the base image we need for our
## go application
FROM golang:1.15
## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
ADD . /app
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app

## we run go build to compile the binary
## executable of our Go program
RUN go build -o server .
## Our start command which kicks off
## our newly created binary executable
EXPOSE 3000
CMD ["/app/server", "3000"]
