# Afterlight

**Afterlight** is a self-hosted "Dead Man's Switch" application. It ensures your critical information like passwords, documents, instructions and personal messages are securely delivered to trusted beneficiaries only if you fail to check in continuously within a specified timeframe. It also has a verifier system to prevent accidental releases.

Designed for privacy and simplicity, Afterlight runs as a single binary with an embedded database, making it easy to audit and deploy on a low-cost VPS or home server.

---

## Features

- **Secure Vaults:** Encrypt and store sensitive text or files that remain locked until a trigger event.
- **Dead Man's Switch:** Automatic release mechanism based on a custom check-in timer (e.g., 30 days).
- **Beneficiary Management:** Assign different trusted contacts to different vaults.
- **Verifier Quorum:** Require m-of-n verifiers to confirm your inactivity before releasing data.
- **Single Binary:** The React frontend is embedded into the Go binary. No complex web server setup required.
- **Zero-Config Storage:** Uses SQLite and local file storage. No external database server needed.
- **Docker Ready:** Production-ready Docker container and Compose setup included.

---

## Tech Stack

**Backend**
- **Language:** Go (Golang) 1.25+
- **Database:** SQLite (Postgres coming soon)
- **Architecture:** REST API with embedded static assets

**Frontend**
- **Framework:** React + TypeScript (Vite)
- **UI System:** Tailwind CSS + shadcn/ui

---

## Quick Start (Docker)

The easiest way to run Afterlight is using Docker Compose. This ensures your data persists in a local volume.

### 1. Build and Run
```bash
docker compose up --build -d
```

### 2. Access the App
Navigate to `http://localhost:8080` in your browser.

### 3. Data Persistence
Your database and encrypted files are stored in the `./afterlight-data` folder created in your project root. Back up this folder regularly.

---

## Development Setup
If you want to contribute or modify the code, you can run the backend and frontend separately.

**Prerequisites**
- Go 1.25+
- Node.js 24+

### 1. Start the Backend
The backend runs on port `8080`.
```bash
go mod download
go run .
```

### 2. Start the Frontend
The frontend runs on port `5173` and proxies API requests to the backend.
```bash
cd web
npm install
npm run dev
```

---

## Configuration
The application is configured via Environment Variables (automatically handled if using Docker).
| Variable         | Description                          | Default Value          |
|------------------|--------------------------------------|------------------------|
| `DB_PATH`        | Path to the SQLite database file     | `/data/afterlight.db`  |
| `ARTIFACTS_PATH` | Directory to store encrypted files   | `/data/artifacts`      |
| `PORT`           | Port for the backend server          | `8080`                 |

---

## License
Afterlight is open-source software licensed under the [AGPL-3.0 License](LICENSE).
