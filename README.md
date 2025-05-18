# TicketHub

A lightweight, real-time ticket management system that you can self-host

## Key Features

- GitHub-like UI/UX
- Real-time updates (WebSocket)
- Rich content expression with Markdown
- Lightweight and easy-to-deploy infrastructure
- Cross-platform compatibility (x86/ARM)

## Technology Stack

- **Frontend**: Nuxt 3 (Vue 3 + Vite), Pinia, TailwindCSS
- **Backend**: Go 1.22 / Gin
- **Real-time Communication**: WebSocket (SockJS + STOMP)
- **Database**: SQLite 3 (FTS5)
- **Cache**: Redis (optional)
- **Container**: Docker

## Quick Start

### Prerequisites

- Docker and Docker Compose installed

### 1. Clone the Repository

```bash
git clone https://github.com/shimauma0312/module-tickethub.git
cd module-tickethub
```

### 2. Set Up Environment Variables

```bash
cp .env.example .env
# Edit the .env file as needed
```

### 3. Launch the Application

```bash
docker compose up
```


## Development Environment

### Frontend Development

```bash
cd frontend
yarn install
yarn dev
```

Or if you're using npm:

```bash
cd frontend
npm install
npm run dev
```

### Backend Development

```bash
cd backend
go mod tidy
go run main.go
```

## Production Environment

```bash
docker compose -f docker-compose.prod.yml up -d
```
