# Building and running the Docker images
1. Set environment vars in a config, e.g `./config/.env.local`:
```
MYSQL_DATABASE=indefinite_studies_api_db
MYSQL_USER=indefinite_studies_api_user
MYSQL_PASSWORD=password
MYSQL_ROOT_PASSWORD=password
```
2. Check `docker-compose.yml` is appropriate to config that you are going to use (e.g.`docker-compose --env-file ./config/.env.local config`)
3. Build project: `./gradlew clean installDist`
4. Build images: `docker-compose --env-file ./config/.env.local build`
5. Run it: `docker-compose --env-file ./config/.env.local up`