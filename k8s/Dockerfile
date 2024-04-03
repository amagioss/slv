FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
WORKDIR /
COPY ./dist/slv-k8s-operator_linux_${TARGETARCH}*/slv-k8s-operator .
USER 65532:65532
ENTRYPOINT ["/slv-k8s-operator"]