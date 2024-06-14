FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
COPY ./dist/slv-cli_linux_${TARGETARCH}*/slv /slv
WORKDIR /workspace
ENTRYPOINT ["/slv"]