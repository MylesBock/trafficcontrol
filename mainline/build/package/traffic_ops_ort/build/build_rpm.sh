#!/usr/bin/env sh
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# shellcheck shell=ash
trap 'exit_code=$?; [ $exit_code -ne 0 ] && echo "Error on line ${LINENO} of ${0}" >/dev/stderr; exit $exit_code' EXIT;
set -o errexit -o nounset -o pipefail -o xtrace;

#----------------------------------------
importFunctions() {
	TC_DIR=$(\
	 pwd | \
	 xargs -n1 -I {} echo '"{}"' | \
	 jq -r '. | split("/") |  to_entries | .[:(.[] | select(.value == "trafficcontrol").key + 1)] | [.[].value] | join("/")' \
	)
  ORT_DIR=$(find ${TC_DIR} -wholename "*/cmd/traffic_ops_ort" -type d)
	export ORT_DIR TC_DIR;
  functions_sh=$(find ${TC_DIR} -wholename "*functions.sh")
  if [ -z "$functions_sh" ]; then
		echo "error: can't find $functions_sh" >&2;
		return 1;
	fi
	. "$functions_sh";
}

#----------------------------------------
initBuildArea() {
	echo "Initializing the build area for Traffic Ops ORT";
	(mkdir -p "$RPMBUILD"
	 cd "$RPMBUILD"
	 mkdir -p SPECS SOURCES RPMS SRPMS BUILD BUILDROOT) || { echo "Could not create $RPMBUILD: $?"; return 1; }

	local dest;
	dest=$(createSourceDir traffic_ops_ort);
	cd "$ORT_DIR";

	echo "PATH: $PATH";
	echo "GOPATH: $GOPATH";
	go version;
	go env;

	go get -v golang.org/x/crypto/ed25519 golang.org/x/crypto/scrypt golang.org/x/net/ipv4 golang.org/x/net/ipv6 golang.org/x/sys/unix;

	gcflags=''
	ldflags=''
	tags='osusergo netgo'
	{ set +o nounset;
	if [ "$DEBUG_BUILD" = true ]; then
		echo 'DEBUG_BUILD is enabled, building without optimization or inlining...';
		gcflags="${gcflags} all=-N -l";
	else
		ldflags="${ldflags} -s -w"; # strip binary
	fi;
	set -o nounset; }

	(cd atstccfg;
	go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}" -tags "$tags")

	(cd t3c;
	go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}" -tags "$tags")

	cp -p traffic_ops_ort.pl "$dest";
	cp -p supermicro_udev_mapper.pl "$dest";
	mkdir -p "${dest}/build";

	echo "build_rpm.sh lsing for logrotate";
	ls -lah .;
	ls -lah ./build;

	cp -p build/atstccfg.logrotate "$dest"/build;
	mkdir -p "${dest}/atstccfg";
	cp -a atstccfg/* "${dest}/atstccfg";
	tar -czvf "$dest".tgz -C "$RPMBUILD"/SOURCES "$(basename "$dest")";
	cp build/traffic_ops_ort.spec "$RPMBUILD"/SPECS/.;
	cp build/atstccfg.logrotate "$RPMBUILD"/.;

	echo "The build area has been initialized.";
}

#----------------------------------------
importFunctions;
checkEnvironment go;
initBuildArea;
buildRpm traffic_ops_ort;
