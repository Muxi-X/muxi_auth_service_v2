FROM golang:latest 
RUN apt install git
CMD ["mkdir", "$GOPATH/src"] \
&&  ["mkdir", "$GOPATH/src/github.com"] \
&&  ["mkdir", "$GOPATH/src/github.com/Muxi-X"] \
&&  ["mkdir", "$GOPATH/src/github.com/Muxi-X/muxi_auth_service_v2"]
ADD . $GOPATH/src/github.com/Muxi-X/muxi_auth_service_v2
WORKDIR $GOPATH/src/github.com/Muxi-X/muxi_auth_service_v2
RUN go build -o main . 
EXPOSE 8033 25 465 587
CMD ["./main"]
