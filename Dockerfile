FROM debian:sid
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apt-get update && apt-get install -y \
	ca-certificates \
	fluidsynth \
	fluid-soundfont-gm \
	--no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

COPY . /go/src/github.com/jfrazelle/cliaoke

RUN buildDeps=' \
		gcc \
		golang \
		git \
		libc6-dev \
	' \
	set -x \
	&& apt-get update \
	&& apt-get install -y  $buildDeps --no-install-recommends \
	&& cd /go/src/github.com/jfrazelle/cliaoke \
	&& go build -o /usr/bin/cliaoke . \
	&& apt-get purge -y --auto-remove $buildDeps \
	&& rm -rf /var/lib/apt/lists/* \
	&& rm -rf /go \
	&& echo "Build complete."


ENTRYPOINT [ "cliaoke" ]
