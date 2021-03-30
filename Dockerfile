FROM golang as builder
WORKDIR /logbot

COPY src/ /logbot
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o logbot github.com/omrigan/logbot/pkg/main


FROM alpine:3.6
RUN apk --no-cache add ca-certificates
COPY --from=builder /logbot/logbot /logbot
WORKDIR /
ENTRYPOINT ["/logbot"]