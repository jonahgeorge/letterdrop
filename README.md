LetterDrop
===

LetterDrop is an open source alternative to FormKeep that is easy to self-host.

## Developing

```sh
$ cp .env.example .env
```

```sh
$ docker run -d \
  -p 5432:5432 \
  -e POSTGRES_USER=letterdrop \
  -e POSTGRES_PASSWORD=letterdrop \
  postgres:10
```

```sh
$ go run .
```

### Running migrations

https://github.com/golang-migrate/migrate/tree/master/cli#installation

```sh
$ migrate \
  -path=migrations \
  -database=postgres://letterdrop:letterdrop@localhost:5432/letterdrop?sslmode=disable \
  up
```

```sh
$ PGPASSWORD=letterdrop psql -h 127.0.0.1 --user letterdrop <<EOF
  insert into users (name, email, password_digest, is_email_confirmed) 
  values ('admin', lower('admin@letterdrop.com'), crypt('keyboardcat', gen_salt('bf', 8)), true) 
  returning *
EOF
```
