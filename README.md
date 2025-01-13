# Basics
- Go will read the .env file, the property SERVER.PORT will control the server port

# Run in development mode
We need to install air globally with `go install github.com/air-verse/air@latest`, for more info: https://github.com/air-verse/air

With air installed, just run `air`

.env will be committed to be used as an example, in production, we're going to overwrite it 

# Folder structure

- `cmd/`: Application's entry points
- `internal/`: Code for internal use only
- `internal/routes`: Responsible for linking the handler's methods into some router
- `internal/handler`: Handles the http request and response, by gathering the necessary information for service layer and controlling the response
- `internal/service`: Keeps the business rules and validations
- `internal/repository`: Persistence layer abstractions
- `internal/utils`: All-purpose functions and variables
- `internal/pkg`: Domain structures and contracts(interfaces)


# Code standards
- Every domain, when creating the necessary logic, should use its name inside their respective package, e.g: `internal/services/products.go` or `internal/respository/mysql-products.go`

# Git flows & standards 
1. Create the branch: `git checkout -b feature/<NOME_DA_FEATURE> main`
1. Start the development: `git commit -m "feat: start <NOME_DO_REQUISITO>"`
1. In-progress commits:`git commit -m "products POST end-point"`
1. In-progress commits: `git commit -m "fix products save service "`
1. In-progress commits: `git commit -m "..."`
1. Finish the development: `git commit -m "feat: end <NOME_DO_REQUISITO>"`


# Testing
- Unit tests must have **Unit** after Test, e.g 
```go
func TestUnitSeller_Create(t *testing.T){}
```
With that, we can run the unit tests apart from integration tests


# Features
**We're going
**Nice to Have**