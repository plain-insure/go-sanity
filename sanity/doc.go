/*
Package sanity provides a client for the Sanity HTTP API.

This package does not directly handle authentication. You should provide an
`http.Client` instance that can handle authentication for you when creating a
Sanity client, such as with the https://golang.org/x/oauth2 package.

	ctx := context.Background()

	tokenSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "YOUR_SANITY_API_TOKEN"},
	)
	httpClient := oauth2.NewClient(ctx, tokenSrc)

	client := sanity.NewClient(httpClient)

	projects, err := client.Projects.List(ctx)
	// ...
*/
package sanity
