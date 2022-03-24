package grpc

import `github.com/pangum/pangu`

func init() {
	pangu.New().Dependencies(
		newServer,
		newClient,
	)
}
