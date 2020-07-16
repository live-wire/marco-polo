# Marco Polo
---
`Marco Polo born: 1254, died: 1324`

### API 
- `localhost:1324/` Map UI
- `localhost:1324/list` List of services sending traffic 
- `localhost:1324/flush` All GeoJSON points for all services
- `localhost:1324/flush/{service}` GeoJSON points for a particular service
- If a service doesn't send any name for itself, it is mapped to a service called `default`.
