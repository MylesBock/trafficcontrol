module github.com/apache/trafficcontrol

go 1.15
replace github.com/basho/riak-go-client => github.com/basho/riak-go-client v1.7.1-0.20170327205844-5587c16e0b8b
require (
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	github.com/GehirnInc/crypt v0.0.0-20200316065508-bb7000b8a962
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/basho/backoff v0.0.0-20150307023525-2ff7c4694083 // indirect
	github.com/basho/riak-go-client v1.7.0
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/dchest/siphash v1.2.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/gofrs/flock v0.8.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.1.2
	github.com/hydrogen18/stoppableListener v0.0.0-20161101122645-827d760f0663
	github.com/influxdata/influxdb v1.8.3
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.10
	github.com/lestrrat-go/jwx v1.0.5
	github.com/lestrrat/go-jwx v0.0.0-20180221005942-b7d4802280ae
	github.com/lestrrat/go-pdebug v0.0.0-20180220043741-569c97477ae8 // indirect
	github.com/lib/pq v1.8.0
	github.com/miekg/dns v1.1.35
	github.com/ogier/pflag v0.0.1
	github.com/pborman/getopt v1.1.0
	github.com/pkg/errors v0.9.1
	go.etcd.io/bbolt v1.3.5
	golang.org/x/crypto v0.0.0-20201117144127-c1f2f97bffc9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)
