# Marco Polo
---
`Marco Polo born: 1254, died: 1324`
This is a lightweight service that will plot all your incomming request traffic on a map view. The traffic markers keep disappearing after a few seconds (configurable) making it usable for very high traffic scenarios too.

### Development Setup
#### Docker (recommended)
- TODO
#### Local
- Set your `$GOPATH` to `~/go/` if not already set.
- Clone this repository in the following path: `$GOPATH/src/github.com/live-wire/marco-polo`
- Run the server file: `go run server.go`

### API 
- `localhost:1324/` Map UI
- `localhost:1324/list` List of services sending traffic 
- `localhost:1324/flush` All GeoJSON points for all services
- `localhost:1324/flush/{service}` GeoJSON points for a particular service
- If a service doesn't send any name for itself, it is mapped to a service called `default`.
