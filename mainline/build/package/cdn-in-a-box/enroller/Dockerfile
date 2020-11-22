# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

FROM debian:buster AS enroller-builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        wget && \
    apt-get clean

COPY build/package/GO_VERSION /
RUN go_version=$(cat /GO_VERSION) && \
    wget -O go.tar.gz https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/go && \
    rm go.tar.gz
ENV GOPATH=/go

# enroller source and dependencies
# todo tentatively these are no longer necessary
#COPY pkg/ /go/src/github.com/apache/trafficcontrol/lib/
#COPY ./vendor/ /go/src/github.com/apache/trafficcontrol/vendor/
#COPY ./pkg/client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/client/
#COPY ./traffic_ops/vendor/ /go/src/github.com/apache/trafficcontrol/traffic_ops/vendor/
#COPY ./infrastructure/cdn-in-a-box/ /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/
# todo copying the entire mainline dir is a little much
COPY ./ /trafficcontrol
COPY build/package/cdn-in-a-box/traffic_ops/to-access.sh /traffic_ops/to-access.sh
WORKDIR /trafficcontrol
RUN go clean && go get ./... && go build -o enroller ./build/package/cdn-in-a-box/enroller/enroller.go

FROM debian:buster

RUN apt-get update && apt-get install -y \
        netcat curl dnsutils net-tools \
        #to-access dependencies
        jq gettext && \
    apt-get clean
COPY --from=enroller-builder \
   /trafficcontrol/enroller \
   /trafficcontrol/build/package/cdn-in-a-box/enroller/run.sh \
   /trafficcontrol/build/package/cdn-in-a-box/enroller/server_template.json \
   /trafficcontrol/build/package/cdn-in-a-box/traffic_ops/to-access.sh \
   /

COPY build/package/cdn-in-a-box/dns/set-dns.sh \
     build/package/cdn-in-a-box/dns/insert-self-into-dns.sh \
     /usr/local/sbin/
RUN chmod +x /run.sh; chmod -R +x /usr/local/sbin;
WORKDIR /shared/enroller
CMD /run.sh
