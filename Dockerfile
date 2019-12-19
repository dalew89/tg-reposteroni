FROM golang:onbuild
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o reposteroni .
VOLUME /app/data
CMD ["/app/reposteroni"]