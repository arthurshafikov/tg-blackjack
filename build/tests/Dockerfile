FROM golang:1.18-alpine as build

RUN set -ex; \
    apk update; \
    apk add --no-cache git; \
    apk add build-base

ENV CODE_DIR /go/src/
ENV TEST_DIR ./tests

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Add docker-compose-wait tool -------------------
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait

CMD /wait && \
    go test ${TEST_DIR}
