FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
WORKDIR /
COPY ./dist/slv-operator_linux_${TARGETARCH}*/ .
USER 65532:65532
ENTRYPOINT ["/slv-operator"]