FROM golang:1.12 as builder
WORKDIR /go/src
RUN mkdir naive-bayes
WORKDIR /go/src/naive-bayes
ADD . .
RUN go get github.com/gorilla/mux
RUN go get k8s.io/klog
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /go/src/naive-bayes/main /app/
WORKDIR /app
EXPOSE 8080
CMD ["./main"]
