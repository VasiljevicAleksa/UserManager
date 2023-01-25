FROM golang:1

WORKDIR /go/src/app

# copy everything to the image
COPY . .

RUN go mod download

# build service
RUN go build -o usermanager app/main.go

# run tests to make sure everything is ok
RUN go test -v ./...

# expose port where container listen
EXPOSE 9000

CMD ["/go/src/app/usermanager"]