FROM registry.access.redhat.com/ubi9/ubi as build
LABEL maintainer=teamnado-ops@redhat.com

ENV GODIR=/usr/local/go APPDIR=/opt/app CGO_ENABLED=0

RUN curl -sfL --retry 10 -o /tmp/go.tar.gz https://go.dev/dl/go1.25.3.linux-amd64.tar.gz && \
    echo "0335f314b6e7bfe08c3d0cfaa7c19db961b7b99fb20be62b0a826c992ad14e0f /tmp/go.tar.gz" | sha256sum -c && \
    mkdir -p $GODIR && \
    tar --strip-components=1 -zxf /tmp/go.tar.gz --directory $GODIR && \
    rm /tmp/go.tar.gz && \
    mkdir -p $APPDIR

ENV PATH="$GODIR/bin:$PATH"

WORKDIR $APPDIR

COPY . $APPDIR

RUN go build -o caddy cmd/caddy/caddy.go

# Multi-stage
FROM registry.access.redhat.com/ubi9/ubi-minimal as run

WORKDIR /opt/app/

COPY --from=build /opt/app/caddy .

CMD /opt/app/caddy
