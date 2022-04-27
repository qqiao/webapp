module github.com/qqiao/webapp

go 1.16

retract (
	v1.8.1
	v1.8.0
	v1.4.0
)

require (
	cloud.google.com/go/firestore v1.6.1
	github.com/golang-jwt/jwt/v4 v4.4.1
	github.com/google/uuid v1.3.0
	google.golang.org/api v0.75.0
	google.golang.org/grpc v1.46.0
)

require github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
