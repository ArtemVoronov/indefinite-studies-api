#Building and running the Docker image
1. `./gradlew installDist`
2. `docker build -t indefinite-studies-api .`
3. `docker run -p 8080:8080 indefinite-studies-api`