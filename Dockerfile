FROM scratch
ARG TARGETARCH
COPY ./dist/slv_linux_${TARGETARCH}*/ /
ENTRYPOINT ["/slv"]