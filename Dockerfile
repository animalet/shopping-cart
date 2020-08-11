FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./main.go
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
EXPOSE 8000
CMD ["./main", "-server"]
