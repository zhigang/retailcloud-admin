FROM  alpine:3.10

WORKDIR /go/bin/

RUN mkdir config

COPY ./config/config.yml ./config/config.yml
COPY build/app ./

CMD ["./app"]