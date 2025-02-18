FROM hub.furycloud.io/mercadolibre/distroless-go-dev:1.21-mini

RUN apk add mysql-server

ADD .ci/ /commands/
RUN chmod a+x /commands/*