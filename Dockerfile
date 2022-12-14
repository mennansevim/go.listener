FROM [REPLACEME]/golang:1.17 as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o main

FROM [REPLACEME]/alpine:3.11.0 as runner
COPY --from=builder /app/main /app/main
COPY --from=builder /app/config/configs.yml /app/config/configs.yml

WORKDIR /app
RUN chmod +x main
ENTRYPOINT ["./main"]
