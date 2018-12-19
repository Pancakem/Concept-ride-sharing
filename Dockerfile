FROM golang:alpine

RUN apk update && apk add --no-cache git

COPY . $GOPATH/src/github.com/pancakem/rides


WORKDIR $GOPATH/src/github.com/pancakem/rides

RUN go get -d -v ./...

RUN go build ./v1/src/cmd/app 

EXPOSE 4000
ENTRYPOINT [ "./app" ]
