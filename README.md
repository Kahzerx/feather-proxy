# feather-proxy

API that redirects HTTP maven traffic and adds redis distributed locks on metadata.xml download and upload for safe multi publishing on repos like reposilite.

It is designed to run inside a workflow

## env variables

- `API_PORT`
  + default: `8000`
  + desc: port on which the service will be listening at
- `REDIS_HOST`
  + default: `127.0.0.1`
  + desc: redis address
- `REDIS_PORT`
  + default: `6379`
  + desc: redis port
- `MAVEN_SCHEME`
  + default: `http`
- `MAVEN_HOST`
  + default: `127.0.0.1`
- `MAVEN_PORT`
  + default: `80`
