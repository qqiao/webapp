module github.com/qqiao/webapp

go 1.16

retract (
	v1.4.0
	v1.8.0
	v1.8.1
)

require (
	cloud.google.com/go/firestore v1.6.1
	github.com/golang-jwt/jwt/v4 v4.4.1
	github.com/google/uuid v1.3.0
	google.golang.org/api v0.74.0
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	golang.org/x/net v0.0.0-20220403103023-749bd193bc2b // indirect
	golang.org/x/sys v0.0.0-20220406163625-3f8b81556e12 // indirect
	google.golang.org/genproto v0.0.0-20220405205423-9d709892a2bf // indirect
)
