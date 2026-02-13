# Go URL Shortener (Microservices Architecture)

A high-performance URL shortener service built with Go, featuring a Cache-Aside architecture to minimize database load and maximize speed.

## System Architecture
- **API:** Golang (Standard Library + `pgx` + `go-redis`)
- **Database:** PostgreSQL (Persistent storage for URLs and click analytics)
- **Cache:** Redis (In-memory storage for high-speed redirects)
- **Orchestration:** Docker Compose
- **CI/CD:** GitHub Actions (Automated builds & Docker Hub deployment)

## Tech Stack
- **Go 1.25**
- **PostgreSQL 15**
- **Redis (Alpine)**
- **Docker & Docker Compose**

## Getting Started

### Prerequisites
- Docker & Docker Compose installed

### Installation
1. Clone the repository:
   ```bash
   git clone [https://github.com/YOUR_USERNAME/go-url-shortener.git](https://github.com/YOUR_USERNAME/go-url-shortener.git)
   cd go-url-shortener```
2. Spin up the entire stack
    ```bash 
    sudo docker compose up --build```
3. Shorten the URL
    ```bash
    curl -X POST -d '{"url": "[https://google.com](https://google.com)"}' http://localhost:8080/shorten```
4. Redirect
    Visit http://localhost:8080/{short_code} in your browser.

### Optimization: Cache-Aside Pattern
This project implements a cache-avoidance strategy. When a redirect is requested:

1. The system checks Redis (RAM).

2. If it's a Cache Hit, the user is redirected instantly.

3. If it's a Cache Miss, the system queries PostgreSQL, populates the cache for future requests, and then redirects.

4. Async Analytics: Click counts are updated in the background via Goroutines to ensure zero latency for the user.


### The Final Step
1.  **Save** this `README.md`.
2.  **Commit and Push:**
    ```bash
    git add README.md
    git commit -m "docs: add professional readme and architecture overview"
    git push
    ```