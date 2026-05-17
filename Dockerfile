FROM cgr.dev/chainguard/glibc-dynamic:latest
ARG TARGETARCH
COPY ./dist/slv-app-linux-${TARGETARCH}*/slv /bin/
WORKDIR /workspace
USER 65532:65532
ENV GODEBUG=madvdontneed=1
ENTRYPOINT ["slv"]
