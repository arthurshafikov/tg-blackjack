FROM golang:1.18 as build

ENV CODE_DIR /go/src/
ENV BIN_FILE /go/src/app

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS

RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/app/main.go

FROM alpine:3.9

LABEL MAINTAINERS="arthurshafikov"

ENV BIN_FILE /go/src/app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV ENV_FILE /go/src/app.env
COPY --from=build ${ENV_FILE} ${ENV_FILE}

ENV CONFIGS_FOLDER /go/src/configs
COPY --from=build ${CONFIGS_FOLDER} ${CONFIGS_FOLDER}

# Add docker-compose-wait tool -------------------
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait

CMD /wait && \
    ${BIN_FILE} -cfgFolder ${CONFIGS_FOLDER} -env /go/src/
