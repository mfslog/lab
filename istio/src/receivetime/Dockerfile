FROM alpine:3.12
COPY bin/receivetime /usr/local/bin

ENV PORT 50051

EXPOSE $PORT


CMD ["sh","-c","cd /usr/local/bin && ./receivetime"]