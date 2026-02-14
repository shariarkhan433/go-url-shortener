# Go URL Shortener (Microservices Architecture)

A production-grade URL shortener service built with Go, featuring a **Cache-Aside** architecture for speed and an **Nginx Reverse Proxy** for security.

## System Architecture
- **Gateway:** Nginx (Reverse Proxy & Rate Limiting)
- **API:** Golang (Standard Library + `pgx` + `go-redis`)
- **Database:** PostgreSQL (Persistent storage for URLs)
- **Cache:** Redis (In-memory storage for high-speed redirects)
- **Orchestration:** Docker Compose
- **CI/CD:** GitHub Actions (Automated testing & Docker Hub deployment)

## Tech Stack
- **Go 1.25**
- **Nginx (Alpine)**
- **PostgreSQL 15**
- **Redis (Alpine)**
- **Docker & Docker Compose**

## Security Features
- **Reverse Proxy:** The Go API is isolated in a private network; only Nginx is exposed to the host.
- **Rate Limiting:** Nginx limits clients to **10 requests per minute** (with a burst buffer of 5) to prevent abuse/DDoS.

## Getting Started

### Prerequisites
- Docker & Docker Compose installed

### Installation
1. Clone the repository:
   ```bash
   git clone [https://github.com/YOUR_USERNAME/go-url-shortener.git](https://github.com/YOUR_USERNAME/go-url-shortener.git)
   cd go-url-shortener
2. Spin up the entire stack
    ```bash 
    sudo docker compose up --build```
3. Shorten the URL
    ```bash
    curl -X POST -d '{"url": "[https://google.com](https://google.com)"}' http://localhost/shorten_code```
4. Redirect
    Visit http://localhost/{short_code} in your browser.

### Optimization: Cache-Aside Pattern
This project implements a cache-avoidance strategy. When a redirect is requested:

1. The system checks Redis (RAM).

2. If it's a Cache Hit, the user is redirected instantly.

3. If it's a Cache Miss, the system queries PostgreSQL, populates the cache for future requests, and then redirects.
4. Nginx checks the rate limit. If safe, forwards to Go API
5. Async Analytics: Click counts are updated in the background via Goroutines to ensure zero latency for the user.


### The Final Step
1.  **Save** this `README.md`.
2.  **Commit and Push:**
    ```bash
    git add README.md
    git commit -m "docs: add professional readme and architecture overview"
    git push
    ```