FROM ksb-dev.keysystems.local:4567/docker/images/golang:ksb-golang-1.17.1-buster
COPY . $PROJECT_DIR
# качаем зависимости
RUN go mod download
# даем права на выполнение скриптам
RUN chmod +x $PROJECT_DIR/build/scripts/*
RUN mkdir $PROJECT_DIR/builded $PROJECT_DIR/builded/microksbscanner $PROJECT_DIR/builded/ksbagent
# запускаем сборку microKSBScanner и KsbAgent
RUN true \
    && $PROJECT_DIR/build/scripts/build-micro-ksb-scanner.sh \
    && $PROJECT_DIR/build/scripts/build-ksb-agent.sh
VOLUME ["/opt/project/builded"]
