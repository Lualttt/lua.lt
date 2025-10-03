# lua.lt
---

this is meant to be ran as a docker container behind a proxy

## environment variable
 * `LUALT_DRAW_ADDRESS` address of the draw backend (exampe: "http://192.168.50.2")

## volumes
 * `/app/visits.json` for visits persistance

## docker compose.yml example
```yaml
services:
  lualt:
    image: lualt # make sure to "docker build -t lualt ."
    container_name: lualt
    restart: unless-stopped
    environment:
      - LUALT_DRAW_ADDRESS=http://192.168.50.2
    volumes:
      - ./visits.json:/app/visits.json # make sure to run the container and copy the file before mounting
    ports:
      - 8080:8080
```
