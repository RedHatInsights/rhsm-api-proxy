FROM registry.access.redhat.com/ubi9/ubi as build
LABEL maintainer=teamnado-ops@redhat.com

ENV GODIR=/usr/local/go APPDIR=/opt/app CGO_ENABLED=0

RUN curl -sfL --retry 10 -o /tmp/go.tar.gz https://go.dev/dl/go1.24.1.linux-amd64.tar.gz && \
    echo "cb2396bae64183cdccf81a9a6df0aea3bce9511fc21469fb89a0c00470088073 /tmp/go.tar.gz" | sha256sum -c && \
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
