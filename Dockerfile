FROM golang:1.8.1

WORKDIR /go/src/docker-dashboard
COPY . .
RUN go get -u gopkg.in/godo.v2/cmd/godo
RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

EXPOSE 8000

CMD ["godo", "server", "-w"]
