FROM alpine:3.12
COPY bin/receivetime /usr/local/bin
COPY bootstrap.json  /usr/local/bin

ENV GRPC_XDS_BOOTSTRAP=/usr/local/bin/bootstrap.json
ENV GRPC_XDS_EXPERIMENTAL_V3_SUPPORT=true

ENV PORT 50051

EXPOSE $PORT

ARG git_commit=0
ENV git_commit=$git_commit

CMD ["/usr/local/bin/receivetime"]