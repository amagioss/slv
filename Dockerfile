FROM scratch
COPY ./dist/slv_linux_${TARGETARCH}*/ /
ENTRYPOINT ["/slv"]