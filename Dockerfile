FROM golang:1.19.3 as builder
LABEL maintainer="Dipak Parmar <hi@dipak.tech>"
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o server .

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/server .
RUN chmod +x server

ENV PORT=8080
ENV DOMAIN=example.com

EXPOSE 8080
CMD ["./server"]

