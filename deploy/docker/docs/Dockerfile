FROM polinux/mkdocs:1.0.4 as builder

ARG VERSION="undefined"

COPY docs/ /workdir

RUN pip install pymdown-extensions mkdocs-material markdown-include
RUN echo "_Documentation Version: ${VERSION}_" > docs/versioninfo.md
RUN mkdocs build


FROM alpine:edge
RUN apk add --no-cache caddy
COPY --from=builder /workdir/site /docs

EXPOSE 80

WORKDIR /docs
ENTRYPOINT ["/usr/sbin/caddy"]
CMD ["--host", "0.0.0.0", "--port", "80", "--log", "stdout"]
