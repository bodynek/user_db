# User_db REST service

This is basic microservice with REST API supporting storage of User information and retrieving it.

## Deployment

Docker and docker-compose-plugin are needed for deployment. 

```
docker compose build
docker compose up -d
```

## Testing

Integration testing requires running docker containers on machine. Test will use test user with `id=eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee`.

You can run test using 
```
go test -v --count=1
```
