FROM debian:bookworm as runtime

RUN apt-get update && apt-get install

FROM golang:latest as builder

WORKDIR /app

COPY cmd/ticket ./
COPY go.mod ./
COPY go.sum ./

RUN go build -v -o ./dist/ticket

FROM runtime

ARG timezone=Asia/Bangkok
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app

# copy app release
COPY --from=builder /app/dist/ticket ./ticket

# expose default port
EXPOSE 5001

# default run entrypoint
CMD ["./ticket"]
