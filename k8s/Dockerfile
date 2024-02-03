FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
WORKDIR /
COPY ./dist/slv-operator_linux_${TARGETARCH}*/ .
ENTRYPOINT ["/slv-operator"]