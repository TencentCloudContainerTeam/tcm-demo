FROM maven:3.5-jdk-8 as builder
COPY ./ /src/app
WORKDIR /src/app
RUN mvn package -Dmaven.test.skip=true

FROM scratch
COPY --from=builder /root/.m2/repository /root/.m2/repository

# docker build -t tcmdemo/spring-boot-cache -f ./cache.dockerfile .
