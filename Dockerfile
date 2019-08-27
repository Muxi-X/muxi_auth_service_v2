FROM golang:latest 
RUN mkdir $GOPATH/src/github.com/Muxi-X/muxi_auth_service_v2 
ADD . $GOPATH/src/github.com/Muxi-X/muxi_auth_service_v2
WORKDIR $GOPATH/src/github.com/Muxi-X/muxi_auth_service_v2
RUN go build -o main . 
CMD ["./main"]
