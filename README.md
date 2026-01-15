# Craft CMS Broken Link Checker
A fast broken link checker based on SEOmatic sitemap* for Craft CMS projects.

## Available Commands
- `check-links`: Check links to see if they are broken.  
- `fetch-links`: Returns all URLs based on the provided sitemap.  

---

## Usage with Go

### Fetch links

```bash
go run main.go fetch-links [URL]
```
For example: `go run main.go fetch-links https://fredmansky.com/sitemap.xml`

### Check links

```bash
go run main.go check-links [URL]
```

For example: `go run main.go check-links https://fredmansky.com/sitemap.xml -l 100`

📌 **What does `-l` (`--rate-limit`) do?**
- Limits the number of **requests per second** (LPS).
- **Example:** `-l 100` means **100 links are checked per second**.
- Default: `200`

### Basic Authentication

For password-protected sitemaps, use the `-u` and `-p` flags:

```bash
go run main.go check-links [URL] -u [USERNAME] -p [PASSWORD]
```

For example: `go run main.go check-links https://stage.fredmansky.fredmansky.com/sitemap.xml -u username -p password`

| Flag | Long | Description |
|------|------|-------------|
| `-u` | `--username` | Username for Basic Auth |
| `-p` | `--password` | Password for Basic Auth |

---

## Usage with Docker
### Fetch links

```bash
docker run fredmansky/go-link-checker fetch-links [URL]
```

For example: `docker run fredmansky/go-link-checker fetch-links https://fredmansky.com/sitemap.xml`

### Check links

```bash
docker run fredmansky/go-link-checker check-links [URL]
```

For example: `docker run fredmansky/go-link-checker check-links https://fredmansky.com/sitemap.xml`

### Basic Authentication (Docker)

```bash
docker run fredmansky/go-link-checker check-links [URL] -u [USERNAME] -p [PASSWORD]
```

---

## Publish a new docker version

1. Commit and push all changes to gh.
2. Login to Docker cli: `docker login`
3. Build docker image: `docker build -t go-link-checker .`
4. Add new docker tag: `docker tag go-link-checker fredmansky/go-link-checker:latest`
5. Push to docker hub: `docker push fredmansky/go-link-checker:latest`
