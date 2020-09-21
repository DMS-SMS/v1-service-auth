FROM alpine
MAINTAINER Park, Jinhong <jinhong0719@naver.com>

COPY ./auth-service ./auth-service
ENTRYPOINT [ "/auth-service" ]
