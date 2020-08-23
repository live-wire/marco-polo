# MarcoPolo gRPC Client (go)
---
`Golang`

#### Usage
- Import the client like:
```
import (
...
mpc "github.com/live-wire/marco-polo/client"
...
)
```

- Initialize the client once when registering your routes in your server with the address of MarcoPolo server and Name of your service. You can use the `MarcoPoloDecorator` around your handlerFunctions for ease of fetching IP addresses from a client request.
```
{
mpClient, _ := mpc.InitMarcoPoloClient("localhost:1254", "api-service")
defer mpClient.Cleanup() // cleans up gRPC connections

...

r := mux.NewRouter()

r.HandleFunc("/list", mpClient.MarcoPoloDecorator(yourListHandler))
r.HandleFunc("/blah", mpClient.MarcoPoloDecorator(yourBlahHandler))
http.Handle("/", r)
log.Fatal(http.ListenAndServe("<yourport>", nil))
}
``` 

- Your handlers are typical HTTP handler functions like:
```
func yourListHandler(w http.ResponseWriter, req *http.Request) {
    ...
}
```

- You can also use `mpClient.Consume` for directly sending an IP address to the MarcoPolo server yourself.
