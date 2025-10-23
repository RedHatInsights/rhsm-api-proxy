FROM registry.access.redhat.com/ubi9/ubi as build
LABEL maintainer=teamnado-ops@redhat.com

ENV GODIR=/usr/local/go APPDIR=/opt/app CGO_ENABLED=0

RUN curl -sfL --retry 10 -o /tmp/go.tar.gz https://go.dev/dl/go1.25.3.linux-amd64.tar.gz && \
    echo "647ddaa978db85623f51f6698bc0c8a5e5fce350397a7fc362f081561954f6df /tmp/go.tar.gz" | sha256sum -c && \
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
