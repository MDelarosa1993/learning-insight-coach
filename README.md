# Learning Insight Coach

Learning Insight Coach is a full-stack AI reading support app.

The backend is a Go API that helps students understand lesson content by using uploaded instructional material as the source of truth. Teachers or learning platforms can upload lesson content, the app stores the full lesson, chunks it into searchable sections, creates embeddings for retrieval, and lets students ask for explanations, hints, quiz prompts, or vocabulary help.

The project also includes a simple Angular frontend in the `frontend/` folder. The frontend is currently being built as a lightweight demo UI for uploading lesson content and later showing student reader interactions and teacher insights.

---

## What This App Does

Learning Insight Coach is designed for education products that need lesson-aware AI support.

It can:

* Upload instructional documents or lesson content
* Store the full original lesson content in the `documents` table
* Split lesson content into searchable chunks
* Generate embeddings for each document chunk
* Store chunks and embeddings in the database
* Fetch a full document by ID
* Search relevant lesson chunks using cosine similarity
* Generate student-friendly AI responses
* Support multiple reading modes:

  * Simplify
  * Hint
  * Quiz
  * Vocabulary
* Check whether student input is allowed
* Check whether AI responses stay grounded in lesson content
* Detect whether the AI gave away a direct answer during hint or quiz mode
* Log student interactions
* Save teacher insight data internally
* Generate teacher insights from student struggles
* Run basic eval cases for AI behavior testing
* Protect API routes with an API key
* Provide a simple Angular UI for demoing the backend

---

## RAG Workflow

The app uses a Retrieval-Augmented Generation, or RAG, workflow.

```txt
Teacher uploads lesson content
        ↓
App stores the full lesson in the documents table
        ↓
App splits the lesson into chunks
        ↓
App creates embeddings for each chunk
        ↓
Student highlights text or asks a question
        ↓
App embeds the question and highlighted text
        ↓
App finds the most relevant chunks
        ↓
App sends the question and lesson context to OpenAI
        ↓
OpenAI returns a student-friendly response
        ↓
App logs the student interaction
        ↓
Teacher insight endpoints summarize student struggles
```

This makes the app lesson-aware. The AI does not just answer from general knowledge. It answers using relevant content from the uploaded lesson.

---

## Tech Stack

### Backend

* Go 1.22
* Gin
* GORM
* SQLite
* OpenAI API
* go-openai
* godotenv
* Docker
* Docker Compose

### Frontend

* Angular
* TypeScript
* Angular Forms
* Angular HttpClient
* CSS

---

## Requirements

To run this app locally, you need:

* Go 1.22+
* Node.js
* npm
* Angular CLI
* An OpenAI API key
* A terminal
* Optional: Docker Desktop
* Optional: Docker Compose
* Optional: DBeaver, Postman, or another API/database tool

Install Angular CLI if needed:

```bash
npm install -g @angular/cli
```

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
├── data/
│   └── sample_lesson.txt
└── frontend/
    ├── angular.json
    ├── package.json
    ├── package-lock.json
    ├── src/
    └── ...
```

---

## Getting Started

You can run the backend in two ways:

1. Directly with Go
2. With Docker

You can run the frontend separately with Angular.

For local development, you will usually use two terminals:

```txt
Terminal 1: Go backend
Terminal 2: Angular frontend
```

---

## Environment Variables

Create a `.env` file from the example file:

```bash
cp .env.example .env
```

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

### Environment Variable Reference

| Variable          | Description                              |
| ----------------- | ---------------------------------------- |
| `PORT`            | Port where the backend server runs       |
| `API_KEY`         | API key required for protected routes    |
| `OPENAI_API_KEY`  | Your OpenAI API key                      |
| `OPENAI_MODEL`    | Chat model used for student responses    |
| `EMBEDDING_MODEL` | Embedding model used for document search |
| `DATABASE_PATH`   | SQLite database path                     |
| `GIN_MODE`        | Gin mode, usually `debug` or `release`   |

Important:

Do not commit your real `.env` file to GitHub.

---

## Run the Backend without Docker

From the project root:

```bash
go mod tidy
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

## Run the Backend with Docker

Create the data folder:

```bash
mkdir -p data
```

Build and start the app:

```bash
docker compose up --build
```

The API should be available at:

```txt
http://localhost:8080
```

Stop the container:

```bash
docker compose down
```

The SQLite database file will stay in the `data` folder because the app uses a mounted volume.

---

## Run the Frontend

The Angular frontend lives in the `frontend/` folder.

From the project root:

```bash
cd frontend
npm install
ng serve
```

The frontend should be available at:

```txt
http://localhost:4200
```

The backend must also be running at:

```txt
http://localhost:8080
```

---

## Frontend Status

The frontend is intentionally simple for now.

Current goal:

```txt
Teacher uploads lesson content
        ↓
Angular sends the lesson to the Go API
        ↓
Go stores the full lesson, chunks it, and creates embeddings
        ↓
Angular displays the returned document ID
```

Current/planned frontend sections:

```txt
frontend/src/app/
├── features/
│   ├── teacher-upload/
│   ├── student-reader/
│   └── teacher-insights/
└── services/
    └── coach-api.service.ts
```

### Current Feature: Teacher Upload

The teacher upload UI should collect:

* Title
* Subject
* Grade minimum
* Grade maximum
* Product ID
* Lesson content

It calls:

```txt
POST /api/v1/documents
```

and displays the returned `document_id`.

### Planned Feature: Student Reader

The student reader UI will allow a student to:

* View or paste lesson content
* Highlight or paste selected text
* Ask a question
* Choose a mode:

  * `simplify`
  * `quiz`
  * `hint`
  * `vocabulary`

It will call:

```txt
POST /api/v1/reader/respond
```

### Planned Feature: Teacher Insights

The teacher insights UI will show class-level learning patterns.

It will call:

```txt
GET /api/v1/teacher/classes/:class_id/insights
```

---

## CORS for Angular

The Angular app runs at:

```txt
http://localhost:4200
```

The Go API runs at:

```txt
http://localhost:8080
```

Because those are different origins, the backend must allow CORS for local frontend development.

Install Gin CORS:

```bash
go get github.com/gin-contrib/cors
```

In `routes/routes.go`, use CORS middleware near the top of `SetupRouter`:

```go
router.Use(cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:4200"},
	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type", "X-API-Key"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
}))
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

## API Endpoints

### Health Check

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

1. Save the full document content in the `documents` table.
2. Split the content into chunks.
3. Generate embeddings for each chunk.
4. Store the searchable chunks in the `document_chunks` table.
5. Mark the document as `indexed`.
6. Return the `document_id`.

Example request:

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -H "X-API-Key: dev-secret-change-me" \
  -d '{
    "title": "Introduction to PostgreSQL",
    "subject": "Databases",
    "grade_min": 9,
    "grade_max": 12,
    "product_id": "postgres-course",
    "content": "PostgreSQL is an open-source relational database management system. It stores data in tables. A table is made up of rows and columns. Each row represents one record, and each column represents one type of information."
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

Save the `document_id`. You will need it when asking reader questions or fetching the document.

---

## Fetch a Document

```txt
GET /api/v1/documents/:document_id
```

Fetches a document by ID.

This is useful for displaying the full lesson content in a frontend.

Example request:

```bash
curl http://localhost:8080/api/v1/documents/PASTE_DOCUMENT_ID_HERE \
  -H "X-API-Key: dev-secret-change-me"
```

Example response:

```json
{
  "id": "generated-document-id",
  "product_id": "postgres-course",
  "title": "Introduction to PostgreSQL",
  "subject": "Databases",
  "grade_min": 9,
  "grade_max": 12,
  "content": "PostgreSQL is an open-source relational database management system. It stores data in tables...",
  "status": "indexed",
  "created_at": "2026-06-21T12:00:00Z",
  "updated_at": "2026-06-21T12:00:00Z"
}
```

---

## Generate a Student Reader Response

```txt
POST /api/v1/reader/respond
```

Generates a student-friendly response using uploaded lesson content.

The app will:

1. Check the student input.
2. Load the document chunks.
3. Embed the highlighted text and student question.
4. Find the most relevant chunks.
5. Send the question and chunks to OpenAI.
6. Check whether the response is grounded.
7. Check whether the response gave away a direct answer.
8. Save the student interaction.
9. Return the student-facing response.

Teacher insight data is still generated and saved internally, but it should not be returned in the student-facing response.

Example request:

```bash
curl -X POST http://localhost:8080/api/v1/reader/respond \
  -H "Content-Type: application/json" \
  -H "X-API-Key: dev-secret-change-me" \
  -d '{
    "student_id": "student-1",
    "class_id": "database-class",
    "document_id": "PASTE_DOCUMENT_ID_HERE",
    "highlighted_text": "A table is made up of rows and columns.",
    "student_question": "Can you quiz me on this?",
    "grade_level": 10,
    "mode": "quiz"
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
  "response": "Quiz question: In a database table, what are rows and columns used for? Try answering in your own words.",
  "mode": "quiz",
  "grounded_sources": [
    {
      "document_id": "generated-document-id",
      "chunk_id": "generated-chunk-id",
      "excerpt": "PostgreSQL is an open-source relational database management system..."
    }
  ],
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
curl http://localhost:8080/api/v1/teacher/classes/database-class/insights \
  -H "X-API-Key: dev-secret-change-me"
```

Example response:

```json
{
  "class_id": "database-class",
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

## Core Concepts

### Document

A document is an uploaded lesson, article, or instructional text.

The `documents` table stores:

* ID
* Product ID
* Title
* Subject
* Grade range
* Full lesson content
* Status
* Created timestamp
* Updated timestamp

A document can have one of these statuses:

* `pending`
* `indexed`
* `failed`

The `content` field stores the full original lesson text. This is useful for displaying the lesson in the frontend.

---

### Document Chunk

A document chunk is a smaller section of a document.

The app splits large lessons into chunks so it can search smaller pieces of content instead of sending the entire lesson to OpenAI every time.

Each chunk stores:

* Document ID
* Product ID
* Chunk index
* Chunk content
* Embedding

Chunks are used for semantic search and retrieval.

---

### Embedding

An embedding is a numeric representation of text.

The app uses embeddings to compare the meaning of a student question with the meaning of document chunks. This helps the app find relevant lesson content even when the student uses different words from the original lesson.

---

### Student Interaction

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
* Created timestamp

Teacher insight endpoints use these records to find common student struggles.

---

### Teacher Signal

A teacher signal is the internal classification of a student interaction.

Example:

```txt
Concept: vocabulary
Struggle type: unknown word or phrase
Recommended action: review key terms
```

The teacher signal should be saved internally on the student interaction, but it should not be returned in the student-facing reader response.

The teacher can view aggregated insight later through:

```txt
GET /api/v1/teacher/classes/:class_id/insights
```

---

## Database Tables

### documents

Stores the full uploaded lesson and metadata.

```txt
id
product_id
title
subject
grade_min
grade_max
content
status
created_at
updated_at
```

### document_chunks

Stores smaller searchable pieces of each lesson.

```txt
id
document_id
product_id
chunk_index
content
embedding
created_at
```

### student_interactions

Stores student questions and teacher insight signals.

```txt
id
student_id
class_id
document_id
highlighted_text
student_question
mode
concept
struggle_type
grade_level
grounded
created_at
```

If you choose to store generated AI responses, add an `ai_response` column to this table.

---

## Main Application Flow

### 1. Document Upload Flow

```txt
POST /api/v1/documents
        ↓
DocumentHandler.Upload
        ↓
DocumentService.Ingest
        ↓
Save document metadata and full content as pending
        ↓
Split content into chunks
        ↓
Create embeddings for each chunk
        ↓
Save chunks and embeddings
        ↓
Mark document as indexed
        ↓
Return document ID
```

---

### 2. Document Display Flow

```txt
GET /api/v1/documents/:document_id
        ↓
DocumentHandler.Show
        ↓
DocumentService.GetByID
        ↓
Store.GetDocumentByID
        ↓
Return full document including content
```

---

### 3. Student Reader Flow

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
Embed highlighted text and student question
        ↓
Search most relevant chunks
        ↓
Build AI prompt with source context
        ↓
Generate AI response
        ↓
Check groundedness
        ↓
Check direct answer behavior
        ↓
Build teacher signal internally
        ↓
Save student interaction
        ↓
Return student-facing response
```

---

### 4. Teacher Insight Flow

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

## Why Store Content in Both Documents and Chunks?

The app stores lesson content in two places for two different purposes.

### `documents.content`

Stores the full original lesson.

Use this for:

* Displaying the full lesson to students
* Fetching a document by ID
* Preserving the original uploaded content
* Supporting a frontend reading experience

### `document_chunks.content`

Stores smaller searchable pieces of the lesson.

Use this for:

* Embedding search
* Retrieval
* Finding relevant lesson sections
* Sending focused context to OpenAI

This design keeps the app flexible. The full lesson is available for display, while chunks are available for AI search.

---

## Development Commands

### Run the backend

```bash
go run main.go
```

Or:

```bash
make run
```

### Run backend tests

```bash
go test ./... -v -count=1
```

Or:

```bash
make test
```

### Format backend code

```bash
gofmt -w .
```

Or:

```bash
make lint
```

### Build the backend

```bash
go build -o bin/coach main.go
```

Or:

```bash
make build
```

### Run the backend with Docker

```bash
docker compose up --build
```

### Run the frontend

```bash
cd frontend
ng serve
```

---

## Recommended `.gitignore`

Because the Angular app lives inside the Go API repo, make sure the root `.gitignore` does not ignore the whole `frontend/` folder.

Track frontend source files, but ignore dependencies and build output.

Recommended root `.gitignore` additions:

```gitignore
.env
data/*.db
bin/

frontend/node_modules/
frontend/dist/
frontend/.angular/
```

It is okay for `frontend/.gitignore` to exist. Do not delete it.

But if `frontend/.git` exists, remove it so the frontend is not treated as a nested Git repository:

```bash
rm -rf frontend/.git
```

---

## Database

This app currently uses SQLite.

The database path is controlled by this environment variable:

```env
DATABASE_PATH=data/coach.db
```

When running with Docker, the `data` folder is mounted so the SQLite database file can persist outside the container.

This means your data can survive when the container stops.

If you added the `content` column after already creating documents, older documents may have an empty `content` field.

For local development, if you do not need old data, you can reset the database:

```bash
rm data/coach.db
go run main.go
```

Only do this if the existing data is not important.

---

## Cost Notes

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

## Testing with Postman

Most API routes require the API key header:

```txt
X-API-Key: dev-secret-change-me
```

In Postman:

1. Open the request.
2. Go to the `Headers` tab.
3. Add:

```txt
Key: X-API-Key
Value: dev-secret-change-me
```

For JSON requests, also add:

```txt
Key: Content-Type
Value: application/json
```

---

## Troubleshooting

### The server does not start

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

### Docker says an environment variable is missing

Make sure you copied the example file:

```bash
cp .env.example .env
```

Then fill in the required values.

---

### The app starts, but OpenAI requests fail

Check that your OpenAI API key is valid:

```env
OPENAI_API_KEY=sk-your-real-key-here
```

Also make sure your OpenAI account has access to the models in your `.env` file.

---

### The database does not persist

Make sure the `data` folder exists:

```bash
mkdir -p data
```

Also make sure your Docker Compose file mounts the folder:

```yaml
volumes:
  - ./data:/app/data
```

---

### Angular requests fail with CORS errors

Make sure the backend allows requests from:

```txt
http://localhost:4200
```

Use Gin CORS middleware in `routes/routes.go`.

---

### Angular build fails with number/null errors

If a form field starts empty, use:

```ts
gradeMin: number | null = null;
gradeMax: number | null = null;
```

Then validate before sending the request:

```ts
if (this.gradeMin === null || this.gradeMax === null) {
  this.errorMessage = 'Please fill out all fields.';
  return;
}

const gradeMin: number = this.gradeMin;
const gradeMax: number = this.gradeMax;
```

---

### MaxTokens error

If you see an error like:

```txt
this model is not supported MaxTokens, please use MaxCompletionTokens
```

Remove `MaxTokens` from the OpenAI chat request or use `MaxCompletionTokens` for models that require it.

For `gpt-4o-mini`, you can usually omit the token limit while developing.

---

## Future Improvements

Possible next steps:

* Finish the Angular teacher upload UI
* Add the Angular student reader UI
* Add the Angular teacher insights UI
* Add user authentication
* Add teacher accounts
* Add student accounts
* Add class management
* Add document listing endpoint
* Add document update endpoint
* Add document delete endpoint
* Add PostgreSQL for production
* Add Redis and background jobs for document ingestion
* Add streaming responses for students
* Add better evals for groundedness and hint quality
* Add rate limiting
* Add request logging
* Add dashboard charts for teacher insights
* Add tests for handlers and services

---

## Summary

Learning Insight Coach is an AI-powered reading support app.

It helps students understand lesson content while giving teachers visibility into student struggles.

The app combines:

* Full lesson content storage
* Document chunking
* Embeddings
* Vector search
* AI-generated reading support
* Guardrails
* Student interaction logging
* Teacher insight generation
* A simple Angular frontend for demos

It is more than a chatbot because it uses uploaded lesson content as context and gives teachers insight into how students are struggling.
