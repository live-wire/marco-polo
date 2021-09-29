# Marco Polo
---
`Marco Polo born: 1254, died: 1324`
This is a lightweight service that will plot all your incomming request traffic on a map view. The traffic markers keep disappearing after a few seconds (configurable) making it usable for very high traffic scenarios too.

![DockerHub](https://github.com/live-wire/marco-polo/actions/workflows/push-to-dockerhub-on-release.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/live-wire/marco-polo/client.svg)](https://pkg.go.dev/github.com/live-wire/marco-polo/client)
[![Go Report Card](https://goreportcard.com/badge/github.com/live-wire/marco-polo)](https://goreportcard.com/report/github.com/live-wire/marco-polo)

### Deployment Instructions
- `docker run -p 1254:1254 -p 1324:1324 dbatheja/marco-polo:v0.1-alpha`
- :point_up: This will expect traffic to be sent to port 1254 as a `proto/` message. See `client/` for a dummy client.
- Add a ` -dummy` at the end of the command above to run the server with dummy traffic.
- UI wil be available at `localhost:1324`

### Development Setup
#### Docker (recommended)
- `docker build -t marcopolocal`
- `docker run -p 1254:1254 -p 1324:1324 marcopolocal -dummy` (`-dummy` will seed dummy data to the server)

#### Local
- Set your `$GOPATH` to `~/go/` if not already set.
- Clone this repository in the following path: `$GOPATH/src/github.com/live-wire/marco-polo`
- Run the server file: `go run server.go`
- If you don't have any client feeding in any data yet, use `go run server.go -dummy` for feeding in dummy data. 
    - See your live dummy traffic @ `localhost:1324`

### API 
- `localhost:1324` Map UI
- `localhost:1324/list` List of services sending traffic 
- `localhost:1324/flush` All GeoJSON points for all services
- `localhost:1324/flush/{service}` GeoJSON points for a particular service
- If a service doesn't send any name for itself, it is mapped to a service called `default`.
