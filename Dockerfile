FROM golang:1.14.4-alpine3.12

ENV APP_NAME feeder

WORKDIR /go/src/${APP_NAME}
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN adduser -D appto
# RUN chown -R appto:root /go/webserver

USER appto

CMD ["webserver"]