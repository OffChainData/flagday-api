FROM golang:1.13

WORKDIR /go/src
RUN git clone https://github.com/pinzolo/flagday

WORKDIR /go/src/flagday-api
COPY src/* /go/src/flagday-api/

RUN go get
RUN go install
CMD ["sh", "-c", "flagday-api"]