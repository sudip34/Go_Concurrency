Subcription-service
===

## connect to database [postgres]

> Download postgres driver

```
go get github.com/jackc/pgconn

go get github.com/jackc/pgx/v4

go get github.com/jackc/pgx/v4/stdlib
```
## Session

> Download for SESSION

```
go get github.com/alexedwards/scs/v2
```
- this allows differnt stores like store as `cookies`/ store in `database' / `Radis`
- Redis will be used to store user SESSION
- Redis - very fast `in-memeory cach`

> Install/Get Redis package

```
go get github.com/alexedwards/scs/redisstore
```

## Web library for routing : ***go-chi***

> go-chi

```
go get github.com/go-chi/chi/v5

go get github.com/go-chi/chi/middleware 
```


## We will use docker for our post-database, radis and mail server








