# lumana-alerts

Go + [chi](https://github.com/go-chi/chi) service that receives Lumana HTTP alerts.

| Method | Path             | Event        |
|--------|------------------|--------------|
| POST   | `/alerts/plug`   | plug event   |
| POST   | `/alerts/unplug` | unplug event |
| GET    | `/`              | hello world  |
| GET    | `/healthz`       | health check |

Payloads are parsed leniently (any JSON accepted, logged as-is). Tighten once the real Lumana payload is known.

```sh
go run .                # listens on :$PORT (default 8080); ADDR overrides
go test ./...
curl -X POST localhost:8080/alerts/plug -d '{"device":"x"}'
```
