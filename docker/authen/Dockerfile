FROM debian:bookworm as runtime

RUN apt-get update && apt-get install

FROM golang:latest as builder

WORKDIR /app

COPY cmd/authen ./
COPY api api
COPY go.mod ./
COPY go.sum ./

# RUN go install github.com/cosmtrek/air@latest
# COPY .air.toml ./

RUN go build -v -o ./dist/authen

FROM runtime

ARG timezone=Asia/Bangkok
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app

# copy app release
COPY --from=builder /app/dist/authen ./authen

# expose default port
EXPOSE 5000

# default run entrypoint
CMD ["./authen"]
