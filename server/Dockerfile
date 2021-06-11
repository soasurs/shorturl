FROM golang:1.16.5

WORKDIR /app

COPY . /app/

ENV GOPROXY "https://goproxy.cn,direct"

RUN go build -o /app/exec shorturl/server

CMD ["/app/exec"]

EXPOSE 3000