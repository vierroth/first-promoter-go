This library provides a thin go wrapper for the [FirstPromoter API](https://docs.firstpromoter.com/api-reference-v2/api-admin/introduction), providing a typesafe and easy to use interface.

## Usage

### Installation

The package is available on [GitHub](https://github.com/vierroth/first-promoter-go) and can be installed as follows:

```bash
go get github.com/vierroth/first-promoter-go@latest
```

### Setup

To begin with, create a client instance, providing the `accountId`, `API key` and transport:

```go
Client := firstpromoter.New( "AccountId", "ApiKey", http.DefaultClient)
```

### Usage

```go
resp, err := Client.TrackSignUp(ctx, firstpromoter.TrackSignUpInput{
	Email: firstpromoter.String("test@test.com"),
	Tid:   firstpromoter.String("tid"),
})
```
