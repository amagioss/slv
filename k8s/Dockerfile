FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
WORKDIR /
COPY ./dist/k8slv_linux_${TARGETARCH}*/ .
ENTRYPOINT ["/k8slv"]