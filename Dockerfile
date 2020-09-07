FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build
RUN apk add git
RUN go get github.com/prometheus/client_golang/prometheus
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main . 

FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /build/www/ /app/www
WORKDIR /app
CMD ["./main"]