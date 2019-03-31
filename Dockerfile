FROM golang:latest
WORKDIR /app
RUN go get github.com/go-sql-driver/mysql && \
    go get github.com/julienschmidt/httprouter
COPY start.sh .
COPY main.go .
COPY template.html .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/start.sh .
COPY --from=0 /app/app .
COPY --from=0 /app/template.html .
RUN chmod 755 /root/start.sh
CMD ["/root/start.sh"]
