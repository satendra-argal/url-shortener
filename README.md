
# URL Shortener

A lightweight URL shortening service written in Go, designed to generate and manage short URLs efficiently. The project uses multi-stage Docker builds to create a minimal, production-ready container image.

---

## Features

- Shorten long URLs into 6-character short URLs.
- Redirect users from short URLs to their original URLs.
- Track top domains based on the number of URLs shortened.

---

## Requirements

- [Go 1.23.4](https://go.dev/dl/) or later
- [Docker](https://www.docker.com/) (optional for containerization)

---

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/sate9584/url-shortener.git
cd url-shortener
```

# Run Locally Without Docker
```bash
go mod download
```

```bash
go build -o url-shortener
./url-shortener
```

# Access the Application
```bash
http://localhost:8080
```

# Dockerized Setup:

1. build docker image
```bash
docker build -t url-shortener:1.0 .
```
2. Run the container:
  ```bash
docker run -p 8080:8080 url-shortener:1.0
```

