FROM golang:1.16.5

RUN mkdir /app

WORKDIR /app

COPY . /app/

ENV GOPROXY "https://goproxy.cn,direct"

RUN go build -o /app/exec shorturl/transform

CMD ["/app/exec"]

EXPOSE 3000