ARG GO_VERSION=1.18

FROM golang:${GO_VERSION} AS build

ARG app_name="nhl-recap"
ARG app_path="/app"

ENV GO111MODULE=on \
    APP_BUILD_PATH="${app_path}" \
    APP_BUILD_NAME=${app_name} \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR ${APP_BUILD_PATH}

COPY . .

RUN go mod download

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH}  go build   -o ${APP_BUILD_NAME}

FROM ubuntu:22.04

ARG app_name="nhl-recap"
ARG app_path="/app"

ENV APP_BUILD_PATH=${app_path} \
    APP_BUILD_NAME=${app_name}

RUN apt update -y && \
    apt install -y ca-certificates && \
    apt clean -y && \
    apt autoremove -y && \
    rm -rf /var/cache/apt/* && \
    rm -rf /var/lib/apt/lists/*

WORKDIR ${APP_BUILD_PATH}

COPY --from=build --chmod=+x ${APP_BUILD_PATH}/${APP_BUILD_NAME} ${APP_BUILD_PATH}/${APP_BUILD_NAME}
COPY --from=build ${APP_BUILD_PATH}/nhl/logos/*.png ${APP_BUILD_PATH}/nhl/logos/

ENTRYPOINT ["/app/nhl-recap"]

