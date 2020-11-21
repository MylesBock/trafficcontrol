module github.com/apache/trafficcontrol/experimental

go 1.15

replace github.com/apache/trafficcontrol => ../go

require (
	github.com/apache/trafficcontrol v1.1.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/lib/pq v1.8.0
)
