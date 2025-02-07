package graphql

import "github.com/ArtemSarafannikov/OzonTestTask/internal/service"

type Resolver struct {
	PostService *service.PostService
}
