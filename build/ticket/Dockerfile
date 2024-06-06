FROM cosmtrek/air AS development

WORKDIR /app

COPY cmd/ticket ./
COPY api/ticket api/ticket
COPY pkg pkg
COPY go.mod go.sum ./
COPY config config
COPY certs certs
COPY build/ticket ./

RUN go mod tidy
RUN go build -v -o ./dist/ticket

CMD ["air", "-c", ".air.toml"]

FROM debian:bookworm

ARG timezone=Asia/Bangkok
ARG hotreload=false
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app

RUN apt-get update && apt-get install

COPY --from=development /app/dist/ticket .
COPY --from=development  /app/config config
COPY --from=development  /app/certs certs

RUN chmod +x ticket

CMD ["./ticket"]