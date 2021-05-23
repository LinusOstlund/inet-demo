# docker build -t my-go-client .
# docker run --name go-client-1 -it --rm --network="host" my-go-client
# docker run --name go-client-2 -it --rm --network="host" my-go-client

# TODO tunna ut denna image, behöver egentligen bara gcc sen kan jag ha en alpine image?

FROM golang:1.15

RUN apt-get update && apt-get install -y --no-install-recommends \
        gcc \
        pkg-config \
        libncurses5-dev \
        libncursesw5-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN go get -u github.com/rthornton128/goncurses

COPY . .

RUN go build -o main .
# vilken kommentar för att cacha?
# docker creates host which conncets to the name 'server' sepcified in the docker-compose
CMD ["/app/main", "0.0.0.0:3000"]