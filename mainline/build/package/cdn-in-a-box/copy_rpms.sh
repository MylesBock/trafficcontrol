#!/usr/bin/env bash
TC_DIR=$(\
	 pwd | \
	 xargs -n1 -I {} echo '"{}"' | \
	 jq -r '. | split("/") |  to_entries | .[:(.[] | select(.value == "trafficcontrol").key + 1)] | [.[].value] | join("/")' \
)
DIST_DIR=$(find ${TC_DIR} -wholename "*trafficcontrol/dist" -type d)

if [  -z "$DIST_DIR" ]; then
  echo "Run docker-compose -f ${TC_DIR}/mainline/build/package/build/docker-compose.yml up before running this!"
  return 1
fi
#edge/traffic_ops_ort.rpm
#mid/traffic_ops_ort.rpm
#traffic_monitor/traffic_monitor.rpm
#traffic_ops/traffic_ops.rpm
#traffic_portal/traffic_portal.rpm
mkdir -p ${TC_DIR}/mainline/build/package/cdn-in-a-box/rpm
NEEDED_RPMS=("traffic_ops_ort" "traffic_ops" "traffic_portal" "traffic_monitor" "tomcat" "traffic_router" "traffic_stats")
for NEEDED_RPM in ${NEEDED_RPMS[@]}
do
  COPY_ME=$(find ${ACTUAL_DIST} -name "*${NEEDED_RPM}-*x86_64.rpm" -print0 | \
  xargs -r -0 ls -1 -t | \
  xargs -n1 -I {} echo '"{}"' | \
  jq --slurp -r '.[0]')
  cp $COPY_ME "${TC_DIR}/mainline/build/package/cdn-in-a-box/rpm/${NEEDED_RPM}.rpm"
done
#find -name "*.rpm")
