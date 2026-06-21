# Learning Insight Coach

Learning Insight Coach is a Go-based AI reading support API.

It helps students understand lesson content by using uploaded instructional material as the source of truth. Teachers or learning platforms can upload lesson content, ask the API to support students with explanations, hints, quiz prompts, and vocabulary help, and then review student struggles through teacher insight endpoints.

The app uses a Retrieval-Augmented Generation, or RAG, workflow:

```txt
Upload lesson content
        ↓
Split the lesson into chunks
        ↓
Create embeddings for each chunk
        ↓
Student asks a question about the lesson
        ↓
Find the most relevant chunks
        ↓
Send the question and lesson context to OpenAI
        ↓
Return a student-friendly response
        ↓
Log the student interaction
        ↓
Generate teacher insights
```

---

## What This App Does

Learning Insight Coach is designed for education products that need AI reading support.

It can:

* Upload instructional documents or lesson content
* Split documents into searchable chunks
* Generate embeddings for each document chunk
* Search relevant lesson chunks using cosine similarity
* Generate student-friendly AI responses
* Support different reading modes:

  * Simplify
  * Hint
  * Quiz
  * Vocabulary
* Check whether student input is allowed
* Check whether AI responses stay grounded in lesson content
* Detect whether the AI gave away a direct answer during hint or quiz mode
* Log student interactions
* Generate teacher insights from student struggles
* Run basic eval cases for AI behavior testing
* Protect API routes with an API key

---

## Tech Stack

* Go 1.22
* Gin
* GORM
* SQLite
* OpenAI API
* go-openai
* godotenv
* Docker
* Docker Compose

---

## Requirements

To run this app locally, you need:

* Go 1.22+
* Docker Desktop
* Docker Compose
* An OpenAI API key
* A terminal
* Optional: DBeaver, Postman, or another API/database tool

---

## Project Structure

```txt
learning-insight-coach/
├── go.mod
├── .env.example
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── main.go
├── config/
│   └── config.go
├── models/
│   └── models.go
├── store/
│   ├── db.go
│   ├── document_store.go
│   ├── interaction_store.go
│   └── vector_store.go
├── services/
│   ├── llm/
│   │   └── client.go
│   ├── document/
│   │   └── service.go
│   ├── guardrail/
│   │   └── service.go
│   ├── reader/
│   │   └── service.go
│   ├── teacher/
│   │   └── service.go
│   └── eval/
│       └── service.go
├── handlers/
│   ├── health.go
│   ├── document.go
│   ├── reader.go
│   ├── teacher.go
│   └── eval.go
├── middleware/
│   └── auth.go
├── routes/
│   └── routes.go
├── evals/
│   └── test_cases.json
└── data/
    └── sample_lesson.txt
```

---

## Getting Started

You can run the app in two ways:

1. With Docker
2. Directly with Go

The Docker setup is the easiest way for someone else to run the project because it uses your `Dockerfile` and `docker-compose.yml`.

---

# Run the App with Docker

## 1. Clone the repository

```bash
git clone YOUR_REPO_URL
cd learning-insight-coach
```

Replace `YOUR_REPO_URL` with the real GitHub repository URL.

---

## 2. Create the environment file

Copy the example environment file:

```bash
cp .env.example .env
```

Open the new `.env` file and add your real values.

Example:

```env
PORT=8080
API_KEY=dev-secret-change-me
OPENAI_API_KEY=sk-your-real-key-here
OPENAI_MODEL=gpt-4o-mini
EMBEDDING_MODEL=text-embedding-3-small
DATABASE_PATH=data/coach.db
GIN_MODE=debug
```

Important:

Do not commit your real `.env` file to GitHub.

---

## 3. Create the data folder

The app uses SQLite. The database file is stored in the `data` folder.

Create the folder if it does not already exist:

```bash
mkdir -p data
```

The database will be created here:

```txt
data/coach.db
```

---

## 4. Build and start the app

Run:

```bash
docker compose up --build
```

The API should be available at:

```txt
http://localhost:8080
```

---

## 5. Test the health endpoint

In another terminal, run:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "message": "Learning Insight Coach API is running",
  "status": "ok"
}
```

---

## 6. Stop the app

To stop the running container:

```bash
docker compose down
```

Your SQLite database file will stay in the `data` folder because the app uses a mounted volume.

---

# Run the App without Docker

You can also run the app directly with Go.

## 1. Install dependencies

```bash
go mod tidy
```

## 2. Create the environment file

```bash
cp .env.example .env
```

Add your OpenAI API key to `.env`.

## 3. Start the app

```bash
go run main.go
```

Or use the Makefile:

```bash
make run
```

The API should be available at:

```txt
http://localhost:8080
```

---

## Environment Variables

The app uses environment variables for configuration.

| Variable          | Description                              |
| ----------------- | ---------------------------------------- |
| `PORT`            | Port where the server runs               |
| `API_KEY`         | API key required for protected routes    |
| `OPENAI_API_KEY`  | Your OpenAI API key                      |
| `OPENAI_MODEL`    | Chat model used for student responses    |
| `EMBEDDING_MODEL` | Embedding model used for document search |
| `DATABASE_PATH`   | SQLite database path                     |
| `GIN_MODE`        | Gin mode, usually `debug` or `release`   |

Example `.env`:

```env
PORT=8080
API_KEY=dev-secret-change-me
OPENAI_API_KEY=sk-your-real-key-here
OPENAI_MODEL=gpt-4o-mini
EMBEDDING_MODEL=text-embedding-3-small
DATABASE_PATH=data/coach.db
GIN_MODE=debug
```

---

## Authentication

Most API routes are protected with an API key.

Send the API key with this header:

```txt
X-API-Key: dev-secret-change-me
```

The value must match the `API_KEY` value in your `.env` file.

The health check route does not require authentication.

---

# API Endpoints

## Health Check

```txt
GET /health
```

Checks whether the API is running.

Example:

```bash
curl http://localhost:8080/health
```

Example response:

```json
{
  "message": "Learning Insight Coach API is running",
  "status": "ok"
}
```

---

## Upload a Document

```txt
POST /api/v1/documents
```

Uploads lesson content and indexes it for retrieval.

The app will:

1. Save the document
2. Split the content into chunks
3. Generate embeddings for each chunk
4. Store the searchable chunks
5. Mark the document as indexed

Example request:

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -H "X-API-Key: dev-secret-change-me" \
  -d '{
    "title": "Photosynthesis Basics",
    "subject": "Science",
    "grade_min": 4,
    "grade_max": 6,
    "product_id": "demo-product",
    "content": "Photosynthesis is the process plants use to make their own food. Plants need sunlight, carbon dioxide, and water."
  }'
```

Example response:

```json
{
  "document_id": "generated-document-id",
  "status": "indexed",
  "chunk_count": 1,
  "message": "document uploaded and indexed"
}
```

Save the `document_id`. You will need it when asking reader questions.

---

## Generate a Student Reader Response

```txt
POST /api/v1/reader/respond
```

Generates a student-friendly response using the uploaded lesson content.

The app will:

1. Check the student input
2. Load the document chunks
3. Embed the student question
4. Find the most relevant chunks
5. Send the question and chunks to OpenAI
6. Check whether the response is grounded
7. Check whether the response gave away a direct answer
8. Save the student interaction
9. Return the response and teacher signal

Example request:

```bash
curl -X POST http://localhost:8080/api/v1/reader/respond \
  -H "Content-Type: application/json" \
  -H "X-API-Key: dev-secret-change-me" \
  -d '{
    "student_id": "student-1",
    "class_id": "class-1",
    "document_id": "PASTE_DOCUMENT_ID_HERE",
    "highlighted_text": "Plants need sunlight, carbon dioxide, and water.",
    "student_question": "Can you explain this in an easier way?",
    "grade_level": 5,
    "mode": "simplify"
  }'
```

Supported reader modes:

| Mode         | Purpose                                                |
| ------------ | ------------------------------------------------------ |
| `simplify`   | Explain the highlighted text in easier language        |
| `hint`       | Give a helpful hint without directly giving the answer |
| `quiz`       | Create a practice question or quiz-style response      |
| `vocabulary` | Explain important words or phrases                     |

Example response:

```json
{
  "response": "This means plants use sunlight, air, and water to help make their own food.",
  "mode": "simplify",
  "grounded_sources": [
    {
      "document_id": "generated-document-id",
      "chunk_id": "generated-chunk-id",
      "excerpt": "Photosynthesis is the process plants use to make their own food..."
    }
  ],
  "teacher_signal": {
    "concept": "reading comprehension",
    "struggle_type": "needs support understanding the text",
    "recommended_action": "Review this concept with a short example and ask the student to explain it in their own words."
  },
  "safety": {
    "grounded": true,
    "gave_direct_answer": false,
    "allowed": true
  }
}
```

---

## Get Teacher Insights

```txt
GET /api/v1/teacher/classes/:class_id/insights
```

Returns class-level insight based on saved student interactions.

Example request:

```bash
curl http://localhost:8080/api/v1/teacher/classes/class-1/insights \
  -H "X-API-Key: dev-secret-change-me"
```

Example response:

```json
{
  "class_id": "class-1",
  "total_interactions": 3,
  "unique_students": 2,
  "top_concepts": [
    {
      "concept": "vocabulary",
      "count": 2,
      "percentage": 66.67,
      "struggle_type": "unknown word or phrase"
    }
  ],
  "recommended_review": [
    "Review vocabulary because 2 student interactions showed difficulty with unknown word or phrase."
  ],
  "generated_summary": "Students are showing difficulty with vocabulary. Consider reviewing key lesson terms before moving forward."
}
```

---

## Run Evals

```txt
POST /api/v1/evals/run
```

Runs AI behavior test cases from the eval file.

Example request:

```bash
curl -X POST http://localhost:8080/api/v1/evals/run \
  -H "X-API-Key: dev-secret-change-me"
```

You can also provide a custom eval path:

```bash
curl -X POST "http://localhost:8080/api/v1/evals/run?path=evals/test_cases.json" \
  -H "X-API-Key: dev-secret-change-me"
```

---

# Core Concepts

## Document

A document is an uploaded lesson, article, or instructional text.

It stores information such as:

* Title
* Subject
* Grade range
* Product ID
* Status

A document can have one of these statuses:

* `pending`
* `indexed`
* `failed`

---

## Document Chunk

A document chunk is a smaller section of a document.

The app splits large lessons into chunks so it can search smaller pieces of content instead of sending the entire lesson to OpenAI every time.

Each chunk stores:

* Document ID
* Product ID
* Chunk index
* Content
* Embedding

---

## Embedding

An embedding is a numeric representation of text.

The app uses embeddings to compare the meaning of a student question with the meaning of document chunks. This helps the app find relevant lesson content even when the student uses different words from the original lesson.

---

## Student Interaction

A student interaction is a saved record of a student question and the AI response process.

It includes:

* Student ID
* Class ID
* Document ID
* Highlighted text
* Student question
* Reader mode
* Concept
* Struggle type
* Grade level
* Whether the response was grounded

Teacher insight endpoints use these records to find common student struggles.

---

## Reader Mode

Reader mode controls how the AI should respond to a student.

| Mode         | Purpose                                                |
| ------------ | ------------------------------------------------------ |
| `simplify`   | Explain the highlighted text in easier language        |
| `hint`       | Give a helpful hint without directly giving the answer |
| `quiz`       | Create a practice question or quiz-style response      |
| `vocabulary` | Explain important words or phrases                     |

---

# Main Application Flow

## 1. Document Upload Flow

```txt
POST /api/v1/documents
        ↓
DocumentHandler.Upload
        ↓
DocumentService.Ingest
        ↓
Save document as pending
        ↓
Split content into chunks
        ↓
Create embeddings for each chunk
        ↓
Save chunks
        ↓
Mark document as indexed
        ↓
Return document ID
```

---

## 2. Student Reader Flow

```txt
POST /api/v1/reader/respond
        ↓
ReaderHandler.Respond
        ↓
ReaderService.Respond
        ↓
Check input safety
        ↓
Load document chunks
        ↓
Embed student question
        ↓
Search most relevant chunks
        ↓
Build AI prompt
        ↓
Generate AI response
        ↓
Check groundedness
        ↓
Check direct answer behavior
        ↓
Save student interaction
        ↓
Return response and teacher signal
```

---

## 3. Teacher Insight Flow

```txt
GET /api/v1/teacher/classes/:class_id/insights
        ↓
TeacherHandler.Insights
        ↓
TeacherService.GetInsights
        ↓
Load class interactions
        ↓
Count concepts and struggles
        ↓
Generate recommendations
        ↓
Return teacher insight summary
```

---

# Development Commands

## Run the app

```bash
go run main.go
```

Or:

```bash
make run
```

---

## Run tests

```bash
go test ./... -v -count=1
```

Or:

```bash
make test
```

---

## Format code

```bash
gofmt -w .
```

Or:

```bash
make lint
```

---

## Build the app

```bash
go build -o bin/coach main.go
```

Or:

```bash
make build
```

---

## Run with Docker

```bash
docker compose up --build
```

---

# Database

This app currently uses SQLite.

The database path is controlled by this environment variable:

```env
DATABASE_PATH=data/coach.db
```

When running with Docker, the `data` folder is mounted so the SQLite database file can persist outside the container.

This means your data can survive when the container stops.

---

# Cost Notes

Each student reader request may call OpenAI multiple times:

* One embedding call for the student question
* One chat call for the student response
* One chat call for groundedness checking
* One chat call for direct answer checking

During local development, you may want to temporarily disable some guardrail checks to reduce API usage.

Recommended development models:

```env
OPENAI_MODEL=gpt-4o-mini
EMBEDDING_MODEL=text-embedding-3-small
```

---

# Troubleshooting

## The server does not start

Check that your `.env` file exists:

```bash
ls -a
```

Make sure it contains:

```env
OPENAI_API_KEY=your_openai_api_key_here
API_KEY=dev-secret-change-me
```

---

## Docker says an environment variable is missing

Make sure you copied the example file:

```bash
cp .env.example .env
```

Then fill in the required values.

---

## The app starts, but OpenAI requests fail

Check that your OpenAI API key is valid:

```env
OPENAI_API_KEY=sk-your-real-key-here
```

Also make sure your OpenAI account has access to the models in your `.env` file.

---

## The database does not persist

Make sure the `data` folder exists:

```bash
mkdir -p data
```

Also make sure your Docker Compose file mounts the folder:

```yml
volumes:
  - ./data:/app/data
```

---

## MaxTokens error

If you see an error like:

```txt
this model is not supported MaxTokens, please use MaxCompletionTokens
```

Remove `MaxTokens` from the OpenAI chat request or use `MaxCompletionTokens` for models that require it.

For `gpt-4o-mini`, you can usually omit the token limit while developing.

---

# Future Improvements

Possible next steps:

* Add user authentication
* Add teacher accounts
* Add student accounts
* Add class management
* Add document listing endpoints
* Add document delete/update endpoints
* Add PostgreSQL for production
* Add Redis and background jobs for document ingestion
* Add streaming responses for students
* Add frontend with React or Angular
* Add better evals for groundedness and hint quality
* Add rate limiting
* Add request logging
* Add dashboard charts for teacher insights
* Add tests for handlers and services

---

# Summary

Learning Insight Coach is an AI-powered reading support API.

It helps students understand lesson content while giving teachers visibility into student struggles.

The app combines:

* Document ingestion
* Embeddings
* Vector search
* AI-generated reading support
* Guardrails
* Student interaction logging
* Teacher insight generation

It is more than a chatbot because it uses uploaded lesson content as context and gives teachers insight into how students are struggling.
