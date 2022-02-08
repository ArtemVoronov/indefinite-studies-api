# Building and running the Docker images
1. Set environment vars in the config `./config/.env` e.g.:
```
DATABASE_NAME=indefinite_studies_api_db
DATABASE_USER=indefinite_studies_api_user
DATABASE_PASSWORD=password
DATABASE_ROOT_PASSWORD=password
DATABASE_URL=jdbc:mysql://localhost:3306/indefinite_studies_api_db
DATABASE_DRIVER_NAME=com.mysql.cj.jdbc.Driver
```
2. Check `docker-compose.yml` is appropriate to config that you are going to use (e.g.`docker-compose --env-file ./config/.env config`)
3. Build project: `./gradlew clean installDist`
4. Build images: `docker-compose --env-file ./config/.env build`
5. Run it: `docker-compose --env-file ./config/.env up`

# Database migration

1. Go to config dir:
```
cd ./config/migrations
```
2. Run liqubase with appropriate params , e.g.
```
liquibase --username=indefinite_studies_api_user --password=password --changeLogFile=db.changelog-root.xml --url=jdbc:mysql://localhost:3306/indefinite_studies_api_db update
```
3. For case of rollback for last migration
```
liquibase --username=indefinite_studies_api_user --password=password --changeLogFile=db.changelog-root.xml --url=jdbc:mysql://localhost:3306/indefinite_studies_api_db rollback-count 1
```

More details see on https://liquibase.org/documentation/