FROM golang:1.12.13
WORKDIR $GOPATH/src/github.com/Muxi-X/muxi_auth_service_v2
COPY . $GOPATH/src/github.com/Muxi-X/muxi_auth_service_v2
RUN go build -o main . 
EXPOSE 8083 25 465 587
CMD ["./main", "-c", "conf/config.yaml"]
