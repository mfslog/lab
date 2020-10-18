FROM alpine:3.12
COPY bin/echo /usr/local/bin
COPY bootstrap.json  /usr/local/bin

ENV PORT 50051

ENV GRPC_XDS_BOOTSTRAP=/usr/local/bin/bootstrap.json
ENV GRPC_XDS_EXPERIMENTAL_V3_SUPPORT=false
ENV GRPC_TRACE=all
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=2
ENV GRPC_GO_LOG_SEVERITY_LEVEL=info

EXPOSE $PORT

ARG git_commit=0
ENV git_commit=$git_commit

CMD ["/usr/local/bin/echo"]