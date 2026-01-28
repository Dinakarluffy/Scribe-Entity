# Entity & Classification Service

This repository contains the **Entity & Classification Service**, which analyzes transcript content and extracts entities, tone, style, and safety classifications.

The system consists of three main parts:

- **Backend** (Go) – REST API + database interaction
- **Frontend** (React + Vite) – UI to analyze and view results
- **Python Worker** – Core analysis logic (entities, tone, style, safety)

---

## 1. Prerequisites

Ensure the following are installed on your machine:

- **Go** 1.21+
- **Node.js** 18+
- **npm** 
- **Python** 3.10+
- **PostgreSQL** 
- **pgAdmin 4**
- **Git**

---

## 2. Project Structure

```
PROJECT 1/
│
├── backend/              # Go backend (API + DB)
│   ├── cmd/server        # main.go (entry point)
│   ├── config            # config & DB connection
│   ├── handlers          # API handlers
│   ├── repository        # DB logic
│   ├── routes            # API routes
│   └── go.mod
│
├── frontend/             # React frontend (Vite)
│   └── src/
│
├── Python-Worker/        # Python analysis logic
│
├── .env                  # Environment variables
└── README.md
```

---

## 3. Environment Configuration (.env)

Create a `.env` file in the **project root**:

```
PORT=8080
GEMINI_API_KEY=your_real_gemini_api_key_here
GEMINI_MODEL=gemini-2.5-flash

# Backend (base URL or full endpoint)
BACKEND_URL=https://your-backend.example.com      # or https://host/api/entity-classification/analyze
BACKEND_API_KEY=optional_backend_bearer_token


DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=entity_classification
DB_TABLE=tabel_name
DB_SSLMODE=disable


# Misc
RESULTS_DIR=results
MAX_CHUNK_CHARS=16000

PYTHON_WORKER_PATH=./Python-Worker/main.py
# Python executable (Windows-safe)
PYTHON_EXEC=python
# macOS / Linux: python3
PYTHON_EXEC=python3





```

> The same `.env` file is used by **Go backend** and **Python worker**.

---

## 4. Database Setup (pgAdmin 4)

### 4.1 Connect PostgreSQL in pgAdmin

1. Open **pgAdmin 4**
2. Click **Register → Server**
3. Enter:
   - **Host**: `DB_HOST`
   - **Port**: `5432`
   - **Maintenance DB**: `postgres` or `incubrix`
   - **Username**: `dinakaran_dev`
   - **Password**: `DB_PASSWORD`
   - **SSL Mode**: Require

4. Save and connect

---

### 4.2 Select Database & Role

Open Query Tool and run:

```sql
SET ROLE dinakaran_dev;
```

---

### 4.3 Table Structure

```sql
CREATE TABLE scribe_entity_classification_dev (
  analysis_id UUID PRIMARY KEY,
  transcript_id UUID,
  creator_id UUID,
  entities JSONB,
  tone JSONB,
  style JSONB,
  safety_flags JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 5. Running the Backend (Go)

### 5.1 Install Dependencies

```bash
cd backend
go mod tidy
```

### 5.2 Run Backend Server

```bash
cd backend
go run ./cmd/server
```

Expected output:

```
PostgreSQL connected successfully
Server running on port 8080
```

### 5.3 Backend Endpoints

- **POST** `/api/entity-classification/analyze`
- **GET** `/api/entity-classification/results`
- **GET** `/api/entity-classification/results/{analysis_id}`
- **GET** `/health`

---

## 6. Running the Frontend (React)

### 6.1 Install Dependencies

```bash
cd frontend
npm install
```

### 6.2 Start Frontend

```bash
npm run dev
```

Frontend runs at:

```
http://localhost:5173
```

### 6.3 Frontend Pages


- **Results** – View all analysis results
- **Get Result by ID** – Fetch single record by `analysis_id`

---

## 7. Running the Python Worker

### 7.1 Create Virtual Environment

```bash
cd Python-Worker
python -m venv venv
```

### 7.2 Activate Virtual Environment

**Windows**:
```bash
venv\Scripts\activate
```

**Linux / macOS**:
```bash
source venv/bin/activate
```

### 7.3 Install Dependencies

```bash
pip install -r requirements.txt
```

### 7.4 Run Python Worker

```bash
python main.py
```

The worker:
- Reads transcript input
- Extracts entities, tone, style, safety
- Can insert results into PostgreSQL

---


## 8. Common Issues & Fixes

### ❌ `no pg_hba.conf entry`
✔ Ensure `DB_SSLMODE=require`

### ❌ Frontend not fetching data
✔ Enable CORS in Go backend
✔ Restart backend after changes

### ❌ Database not visible in pgAdmin
✔ Verify correct database name
✔ Run `SET ROLE dinakaran_dev;`

---

## 10. Notes

- Backend controls all DB writes
- Frontend is read-only for results
- Python Worker handles analysis logic
- Designed as per TRD requirements

---


