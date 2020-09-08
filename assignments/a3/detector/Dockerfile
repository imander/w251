FROM golang:1.15
ENV go111module=on

RUN apt-get update && apt-get install sudo
RUN go get -u -d gocv.io/x/gocv
RUN cd $GOPATH/src/gocv.io/x/gocv && make install

WORKDIR /go/mids.w251.edu/a3/detector
COPY go.mod .
COPY go.sum .
COPY main.go . 
RUN go build -o /bin/face_detector main.go
WORKDIR /
COPY data/haarcascade_frontalface_default.xml /data/
ENTRYPOINT ["/bin/face_detector"]
