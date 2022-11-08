ARG BUILD_CONTAINER_IMAGE
FROM $BUILD_CONTAINER_IMAGE as build

FROM ksb-dev.keysystems.local:4567/docker/images/golang:ksb-lightweight-golang-1.12.9
# заимствуем собранные файлы проекта
COPY --from=build /opt/project/builded $PROJECT_DIR

WORKDIR $PROJECT_DIR/microksbscanner
# создаем папку без которого сервис не запуститься
RUN mkdir ./agentzip
EXPOSE 3003 3004
VOLUME ["/opt/project/microksbscanner", "/opt/project/ksbagent"]
CMD ["/opt/project/microksbscanner/microKSBScanner"]
