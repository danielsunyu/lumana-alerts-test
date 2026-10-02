# lumana-alerts

Go + [chi](https://github.com/go-chi/chi) service that receives Lumana HTTP alerts.

| Method | Path             | Event        |
|--------|------------------|--------------|
| POST   | `/alerts/plug`   | plug event   |
| POST   | `/alerts/unplug` | unplug event |
| GET    | `/healthz`       | health check |

Payloads are parsed leniently (any JSON accepted, logged as-is). Tighten once the real Lumana payload is known.

```sh
go run .                # listens on :8080 (override with ADDR)
go test ./...
curl -X POST localhost:8080/alerts/plug -d '{"device":"x"}'
```
