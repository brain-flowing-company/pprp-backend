## Getting Started

**Prerequisite**

- Golang
- Makefile (optional)
- docker
- git
- swagger

1. Clone repository

```bash
git clone https://github.com/brain-flowing-company/pprp-backend.git
```

2. Copy `.env.example` to `.env`

3. Get Swagger cli

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. Run development server

```bash
docker-compose up -d
# or
make up
```

> Note: make sure you have `Makefile` before using `make` command

- This will automatically start postgres database and development server with auto-reload.
- **Swagger (API docs)** is at [localhost:3000/docs](http://localhost:3000/docs) \*_change port if your app is running on different port_
- If you've made any changes to API docs (comments above handler function), make sure you run this command to update API docs page.

```bash
swag init -g ./cmd/main.go -o ./docs/
# or
make docs
```

## Project structures

- `cmd/` contains `main.go`
- `config/` contains env var loader
- `database/` contains database (postgres) connector
- `internal/`
  - name folder by service name (snake_case)
  - `.handler.go` handles http request
  - `.service.go` holds core logic
  - `.repository.go` holds data fetcher (database queries or external api calls)

## Contribution

1. Make sure to pull the latest commit from `dev` branch

```bash
git pull origin dev
```

2. Create new branch with your git GUI tools or use this command

```bash
git checkout -b <branch-name>
```

3. Make sure you on the correct branch
4. Craft your wonderful code and don't forget to commit frequently

```bash
git add <file-path> # add specific file
# or
git add . # add all files
```

```bash
git commit -m "<prefix>: <commit message>"
```

> Note: _check out [commit message convention](#commit-message-convention)_

5. Push code to remote repository
6. Create pull request on [github](https://github.com/brain-flowing-company/pprp-backend/pulls)

- compare changes with **base**: `dev` &#8592; **compare**: `<branch-name>`
- title also apply [commit message convention](#commit-message-convention)
- put fancy description

### Commit message convention

```bash
git commit -m "<prefix>: <commit message>"
```

- use lowercase
- **meaningful** commit message

**Prefix**

- **`feat`**: introduce new feature
- **`fix`**: fix bug
- **`refactor`**: changes which neither fix bug nor add a feature
- **`chore`**: changes to the build process or extra tools and libraries
