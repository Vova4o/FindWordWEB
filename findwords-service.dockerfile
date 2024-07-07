FROM alpine:latest

RUN mkdir /app

COPY findwords /app
COPY templates /templates
COPY static /static
COPY russian_nouns.txt .

CMD ["/app/findwords"]