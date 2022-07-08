# Building and running the Docker images
1. Set environment vars in the config `.env` e.g.:
```
#common settings
APP_PORT=3000

#required for db service inside app
DATABASE_HOST=postgres
DATABASE_PORT=5432
DATABASE_USER=indefinite_studies_api_user
DATABASE_PASSWORD=password
DATABASE_NAME=indefinite_studies_api_db
DATABASE_SSL_MODE=disable

#required for liquibase
DATABASE_URL=jdbc:postgresql://postgres:5432/indefinite_studies_api_db
```
2. Check `docker-compose.yml` is appropriate to config that you are going to use (e.g.`docker-compose config`)
3. Build images: `docker-compose  build`
4. Run it: `docker-compose up`
5.* Stop it: `docker-compose down`

P.S. It uses the services from https://github.com/ArtemVoronov/indefinite-studies-environment