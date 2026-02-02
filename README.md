# seller-accounts-backend

## HTTPS / TLS

TLS can be enabled via `Server.tls` config. When `tls.enabled` is `true` the server starts with HTTPS and, if the cert/key files are missing, it will auto-generate a self-signed pair (when `tls.selfSigned.enabled` is `true`).

Example (see `config/config.dev.yaml` for a ready-to-use block):

```yaml
Server:
  host: "0.0.0.0"
  port: "8080"
  tls:
    enabled: true
    certFile: "config/certs/dev.crt"
    keyFile: "config/certs/dev.key"
    selfSigned:
      enabled: true
      commonName: "localhost"
      hosts:
        - "localhost"
        - "127.0.0.1"
      validForDays: 365
```

First start will create the files; subsequent starts reuse them.

Set `APP_ENV=prod` (or `test`) to pick the matching `config/config.<env>.yaml` file.

## Docker

Собрать образ:
```bash
docker build -t sellers-accounts-backend .
```

Запустить (подставь свой конфиг):
```bash
docker run --rm -p 3000:3000 \
  -e APP_ENV=prod \
  -v $(pwd)/config/config.prod.yaml:/app/config/config.prod.yaml:ro \
  sellers-accounts-backend
```

В `.dockerignore` конфиги исключены, поэтому YAML нужно пробрасывать томом или через переменные окружения/секреты.

Если нужны миграции init.sql — смонтируй/скопируй файл и примени psql перед стартом контейнера базы данных.
