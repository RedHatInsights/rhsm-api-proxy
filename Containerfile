FROM registry.access.redhat.com/ubi9/ubi as build
LABEL maintainer=teamnado-ops@redhat.com

ENV GODIR=/usr/local/go APPDIR=/opt/app CGO_ENABLED=0

RUN curl -sfL --retry 10 -o /tmp/go.tar.gz https://go.dev/dl/go1.21.1.linux-amd64.tar.gz && \
    echo "b3075ae1ce5dab85f89bc7905d1632de23ca196bd8336afd93fa97434cfa55ae /tmp/go.tar.gz" | sha256sum -c && \
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

COPY misc/2015-IT-Root-CA.crt /etc/pki/ca-trust/source/anchors/2015-IT-Root-CA.crt
COPY misc/2022-IT-Root-CA.crt /etc/pki/ca-trust/source/anchors/2022-IT-Root-CA.crt
RUN update-ca-trust

WORKDIR /opt/app/

COPY --from=build /opt/app/caddy .

CMD /opt/app/caddy
