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
set -o errexit -o nounset -o pipefail;
set -o xtrace

#----------------------------------------
importFunctions() {
	# take the cwd, wrap it in double quotes, pipe it into jq
	# then split it on forward slashes, convert it to entries
	# { "key": ARRAY_INDEX, "value": SPLIT_VALUE }
	# select the entry that contains "trafficcontrol" and grab the slice containing all values leading up to it
	# then wrap the values up in an array and join them with forward slashes
	TC_DIR=$(\
	 pwd | \
	 xargs -n1 -I {} echo '"{}"' | \
	 jq -r '. | split("/") |  to_entries | .[:(.[] | select(.value == "trafficcontrol").key + 1)] | [.[].value] | join("/")' \
	)
  TO_DIR=$(dirname $(find ${TC_DIR} -wholename "*traffic_ops/traffic_ops.go"))
  # todo (restructure): gotta figure out how to handle plugins without keeping them in internal/pkg/traffic_ops/plugin
  PLUGIN_DIR=$(find ${TC_DIR} -wholename "*traffic_ops/plugin*"-not -wholename "*traffic_ops/plugins*" -type d)
  echo $TC_DIR
  echo $TO_DIR
	export TO_DIR TC_DIR
	functions_sh=$(find ${TC_DIR} -wholename "*functions.sh")
	if [ ! -r "$functions_sh" ]; then
		echo "error: can't find $functions_sh"
		return 1
	fi
	. "$functions_sh"
}

# ---------------------------------------
initBuildArea() {
	echo "Initializing the build area."
	(mkdir -p "$RPMBUILD"
	 cd "$RPMBUILD"
	 mkdir -p SPECS SOURCES RPMS SRPMS BUILD BUILDROOT) || { echo "Could not create $RPMBUILD: $?"; return 1; }

	local dest
	dest="$(createSourceDir traffic_ops)"
	cd "$TO_DIR" || \
		 { echo "Could not cd to $TO_DIR: $?"; return 1; }

	echo "PATH: $PATH"
	echo "GOPATH: $GOPATH"
	go version
	go env

	# get x/* packages (everything else should be properly vendored)
	go get -v \
		golang.org/x/crypto/ed25519 \
		golang.org/x/crypto/scrypt \
		golang.org/x/net/idna \
		golang.org/x/net/ipv4 \
		golang.org/x/net/ipv6 \
		golang.org/x/sys/unix \
		golang.org/x/text/secure/bidirule ||
		{ echo "Could not get go package dependencies"; return 1; }

	# compile traffic_ops
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
	go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.version=traffic_ops-${TC_VERSION}-${BUILD_NUMBER}.${RHEL_VERSION} -B 0x$(git rev-parse HEAD)" -tags "$tags" || \
								{ echo "Could not build traffic_ops_golang binary"; return 1; }
	cd -

  # todo (restructure) move these to tools
  # also properly change to that directory
	# compile db/admin
	(cd traffic_ops/app/db
	go build -v -o admin -gcflags "$gcflags" -ldflags "$ldflags" -tags "$tags" || \
								{ echo "Could not build db/admin binary"; return 1;})

	# compile TO profile converter
	(cd traffic_ops/install/bin/convert_profile
	go build -v -gcflags "$gcflags" -ldflags "$ldflags" -tags="$tags" || \
								{ echo "Could not build convert_profile binary"; return 1; })

	rsync -av traffic_ops/etc traffic_ops/install "$dest"/ || \
		 { echo "Could not copy to $dest: $?"; return 1; }
	if ! (cd traffic_ops/app; rsync -av bin conf cpanfile db lib public script templates "${dest}/app"); then
		echo "Could not copy to $dest/app"
		return 1
	fi
	tar -czvf "$dest".tgz -C "$RPMBUILD"/SOURCES "$(basename "$dest")" || \
		 { echo "Could not create tar archive $dest.tgz: $?"; return 1; }
	cp traffic_ops/build/traffic_ops.spec "$RPMBUILD"/SPECS/. || \
		 { echo "Could not copy spec files: $?"; return 1; }
  # todo (restructure) update this to reflect new structure
	PLUGINS="$(grep -l 'AddPlugin(' "${TC_DIR}/mainline/internal/pkg/traffic_ops/plugin/"*.go | grep -v 'func AddPlugin(' | xargs -I '{}' basename {} '.go')"
	export PLUGINS

	echo "The build area has been initialized."
}

# ---------------------------------------
importFunctions
checkEnvironment -i go,rsync
initBuildArea
buildRpm traffic_ops
