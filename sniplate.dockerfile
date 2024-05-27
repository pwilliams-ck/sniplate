FROM alpine:latest

RUN mkdir /app

COPY sniplateApp /app

CMD [ "/app/sniplateApp" ]

