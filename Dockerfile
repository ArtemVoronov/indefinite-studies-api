FROM openjdk:8-jdk
EXPOSE 8080:8080
RUN mkdir /app
COPY ./build/install/indefinite-studies-api/ /app/
COPY ./.env /app/bin
WORKDIR /app/bin
CMD ["./indefinite-studies-api"]