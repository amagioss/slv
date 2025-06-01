FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
COPY ./dist/slv_linux_${TARGETARCH}*/slv /bin/
WORKDIR /workspace
USER 65532:65532
ENV GODEBUG=madvdontneed=1
ENTRYPOINT ["slv"]