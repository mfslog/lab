FROM alpine:3.12
COPY bin/echo /usr/local/bin

ENV PORT 50051

EXPOSE $PORT

ARG git_commit=0
ENV git_commit=$git_commit

CMD ["sh","-c","cd usr/local/bin && echo"]