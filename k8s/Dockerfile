FROM cgr.dev/chainguard/static:latest
WORKDIR /
COPY ./dist/k8slv_linux_${TARGETARCH}*/ .
ENTRYPOINT ["/k8slv"]