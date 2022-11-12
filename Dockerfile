FROM golang:1.19 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o httpproxy ./main.go

FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/httpproxy /httpproxy
USER nonroot:nonroot
ENTRYPOINT ["/httpproxy"]

EXPOSE 8080