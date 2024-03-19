FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
COPY ./dist/slv_linux_${TARGETARCH}*/slv /slv
WORKDIR /workspace
ENTRYPOINT ["/slv"]