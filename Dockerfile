FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
COPY ./dist/slv_linux_${TARGETARCH}*/slv /slv
WORKDIR /workspace
USER 65532:65532
ARG SLV_MODE
ENV SLV_MODE=${SLV_MODE}
ENTRYPOINT ["/slv"]