# gitcache

```
┌─────────────┐
│             │
│ GitHub repo │
└──────▲──────┘
       │
 fetch │
       │
┌──────┴──────┐           ┌─────────────┐
│             │   fetch   │             │
│  gitcache   ◄───────────┤    CI/CD    │
└─────────────┘           └─────────────┘
```

Build binary:

    go build -o ./bin/gitcache -trimpath ./cmd/gitcache

Binary usage:

    ./bin/gitcache -h

Deploy to docker:

    CGO_ENABLED=0 GOOS=linux go build -o ./build/gitcache -trimpath ./cmd/gitcache
    docker compose --file ./deployment/docker-compose.local.yaml up --detach --build --force-recreate
    git clone http://localhost:8080/learn-git.git
