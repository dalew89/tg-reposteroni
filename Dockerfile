FROM golang:onbuild
RUN mkdir /app
ADD . /app/
WORKDIR /app
VOLUME . /app/data
RUN go build -o reposteroni .

CMD ["/app/reposteroni"]