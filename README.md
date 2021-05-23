# TCP in Go with Docker and Goncurses

Ett sätt att hantera Server och Client i Golang med Goncurses som gränssnitt. Docker hanterar dependencies för Goncurses.
Detta lilla projekt gjordes för `inet` i kursen DD1362 på KTH. Det är en spartansk demo som bryter mot många goda programmeringsprinciper. 

**Var på er vakt!**

## Dockerfile
För att köra, öppna två terminaler. En fallgrop för Dockernybörjaren är nätverkshantering. I Dockerfilerna specificeras ett nätverk med `--network="host"`, där `host` är ett Docker-keyword. På så sätt kommunicerar containerserna med varandra över localhost. Se övriga kommentarer i Dockerfilerna.

```bash
# Terminal 1
# navigera in i /Server/
docker build -t my-go-server .
docker run --name go-server -it --rm --network="host" my-go-server

# Terminal 2
# navigera in i /Client/
# -it innebär 'interactive mode', så du kopplar på containern direkt
# --rm innebär 'remove' så du tar bort containern när du är klar
docker build -t my-go-client .
docker run --name go-client-1 -it --rm --network="host" my-go-client

# för att lägga till en spelare, öppna ett nytt terminalfönster och kör samma kommando med ett annat namn:
docker run --name go-client-2 -it --rm --network="host" my-go-client
```
