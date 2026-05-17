FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH
COPY ./dist/slv-app-linux-static-${TARGETARCH}*/slv /bin/
WORKDIR /workspace
USER 65532:65532
ENV GODEBUG=madvdontneed=1
ENTRYPOINT ["slv"]
