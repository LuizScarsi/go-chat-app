# Instructions:
## Docker:
`cd docker`
### Docker compose:
- `docker compose --env-file ./.env up -d`
- `docker compose down`
  - Assuming that your .env file is in project's root folder
### Open psql via terminal:
`docker exec -it postgres-go-chat psql -U gouser -W gochatapp`
