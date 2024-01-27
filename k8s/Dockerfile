FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY ./dist/k8slv_linux_${TARGETARCH}*/ .
USER 65532:65532
ENTRYPOINT ["/k8slv"]