FROM openjdk:8-jdk
EXPOSE 8080:8080
RUN mkdir /app
COPY ./build/install/indefinite-studies-api/ /app/
WORKDIR /app/bin
CMD ["./indefinite-studies-api"]