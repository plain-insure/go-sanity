# go-sanity

Go library for accessing the [Sanity HTTP API](https://www.sanity.io/docs/http-api).

## Installation

go-sanity is compatible with modern Go releases in module mode, with Go installed:

```
go get github.com/tessellator/go-sanity
```

## Usage

go-sanity does not directly handle authentication. You should provide an
`http.Client` instance that can handle authentication for you when creating a
Sanity client. The following example uses the
[oauth2](https://pkg.go.dev/golang.org/x/oauth2) library to create an
`http.Client`, creates a new Sanity client, and gets all projects from the
Sanity account.

```go
import (
	"context"

	"github.com/tessellator/go-sanity/sanity"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	tokenSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "YOUR_SANITY_API_TOKEN"},
	)
	httpClient := oauth2.NewClient(ctx, tokenSrc)

	client := sanity.NewClient(httpClient)

	projects, err := client.Projects.List(ctx)
	// ...

	// List webhooks for a project
	webhooks, err := client.Webhooks.List(ctx, "projectId")
	// ...
}
```

To discover your personal auth token, you can run the command
`sanity debug --secrets` at a terminal. You may then create new tokens via the
API.

## Supported APIs

- **Projects API**: Manage Sanity projects, datasets, CORS entries, users, roles, and tokens
- **Webhooks API**: Manage webhook configurations for real-time notifications

## Code structure

The code structure was inspired by [jianyuan/go-sentry](https://github.com/jianyuan/go-sentry).

## License

This library is distributed under the [MIT License](LICENSE).
