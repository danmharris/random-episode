FROM golang:1.23-alpine3.20 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /random-episode ./cmd/main.go

FROM alpine:3.20

LABEL org.opencontainers.image.title="Random Episode"
LABEL org.opencontainers.image.description="Web app for selecting random episodes of TV shows"
LABEL org.opencontainers.image.url="https://github.com/danmharris/random-episode"
LABEL org.opencontainers.image.source="https://github.com/danmharris/random-episode"

ENV UID=1000
ENV GID=1000
ENV PORT=8000

RUN addgroup -S random-episode -g $GID && \
    adduser -S -u $UID -G random-episode random-episode

COPY --from=build /random-episode /usr/local/bin/random-episode

USER $UID
EXPOSE $PORT
CMD ["/usr/local/bin/random-episode"]
