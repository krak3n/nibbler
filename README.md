# Nibbler

Nibbler is a little URL shortener.

## Usage

Build the Docker image:

```
make image
```

Start the server:

```
make up
```

Run migrations:

```
make migrate_up
```

Using `curl` you can now shorten a URL and get a unique ID:

```
short=$(curl http://localhost:3000/shorten\?url\=http://google.com)
```

And reverse it back:

```
curl -L http://localhost:3000/$short
```
