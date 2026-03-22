<p align="center" style="font-size: 24px; font-weight: bold;">
 Authula Playground
</p>

<p align="center">
  <img src="./project-logo.png" height="100" width="250" alt="Authula Logo"/>
</p>

This repository is a playground that demonstrates how to integrate Authula into a full-stack application. It includes a Next.js (SSR) + React Router (SPA) frontends and a Go backend, showcasing how to implement authentication and authorization using Authula.

## Tech Stack

- Frontend: Next.js + React Router
- Backend: Go (Golang), Echo (for this example)
- Library: [Authula](https://github.com/Authula/authula)

## Getting Started

1. **Clone the repository**
2. **Docker Compose**
   - Ensure Docker and Docker Compose are installed.
   - Copy `docker-compose.env.example` to `docker-compose.env`.
   - Update environment variables in `docker-compose.env` as needed.
   - Start services: `docker compose down -v && docker compose --env-file=docker-compose.env up -d`
3. **Backend Setup**
   - Copy `backend/.env.example` to `backend/.env`.
   - Update environment variables in `backend/.env`.
   - Install dependencies: `go mod tidy`
   - Start the backend: `go run main.go`
4. **Frontend Setup**
   - Install dependencies: `pnpm install`
   - Copy `frontend/.env.local.example` to `frontend/.env.local`.
   - Update environment variables in `frontend/.env.local`.
   - Start the frontend: `pnpm dev`

This repository will be updated as Authula evolves, with additional examples demonstrating integration in various scenarios.

For comprehensive documentation and usage instructions, visit the [Authula Docs](https://authula.vercel.app/docs).

---
