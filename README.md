# Learning Insight Coach

Learning Insight Coach is a Go-based AI reading support API that helps students understand instructional content and gives teachers insight into what students are struggling with.

The app allows teachers or products to upload lesson content, stores that content as searchable chunks, uses OpenAI embeddings to find relevant lesson sections, and generates student-friendly responses such as simplified explanations, hints, quiz prompts, and vocabulary support.

It also logs student interactions so teachers can review common struggles by class.

---

## Features

* Upload instructional documents or lesson content
* Split documents into smaller searchable chunks
* Generate embeddings for document chunks
* Search relevant lesson chunks using cosine similarity
* Provide AI-powered student support in different reading modes
* Support modes such as:

  * Simplify
  * Hint
  * Quiz
  * Vocabulary
* Check whether student input is allowed
* Check whether AI responses stay grounded in lesson content
* Detect whether the AI gave away a direct answer in hint or quiz mode
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

## How the App Works

At a high level, the app works like this:

```txt
Teacher uploads lesson content
        ↓
App stores the document
        ↓
App splits the document into chunks
        ↓
App creates embeddings for each chunk
        ↓
Student highlights text and asks a question
        ↓
App embeds the student question
        ↓
App finds the most relevant document chunks
        ↓
App sends the chunks and question to OpenAI
        ↓
AI returns a student-friendly response
        ↓
App saves the student interaction
        ↓
Teacher can view class insights
```

This pattern is known as Retrieval-Augmented Generation, or RAG.

Instead of letting the AI answer from general memory, the app first retrieves relevant lesson content and asks the AI to answer using that content.

---

## Core Concepts

### Document

A document is an uploaded lesson, article, or instructional text.

It contains metadata such as:

* Title
* Subject
* Grade range
* Product ID
* Status

A document can have the status:

* `pending`
* `indexed`
* `failed`

---

### Document Chunk

A document chunk is a smaller section of a document.

The app splits large documents into chunks so it can search smaller pieces of content instead of sending an entire lesson to the AI every time.

Each chunk stores:

* Document ID
* Product ID
* Chunk index
* Content
* Embedding

---

### Embedding

An embedding is a numeric representation of text.

The app uses embeddings to compare the meaning of a student question against the meaning of document chunks.

This allows the app to find relevant lesson content even when the student does not use the exact same words as the lesson.

---

### Student Interaction

A student interaction is a saved record of a student question.

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

These records are used to generate teacher insights.

---

### Reader Mode

Reader mode controls how the AI should respond.

Supported modes:

| Mode         | Purpose                                                |
| ------------ | ------------------------------------------------------ |
| `simplify`   | Explain the highlighted text in easier language        |
| `hint`       | Give a helpful hint without directly giving the answer |
| `quiz`       | Create a practice question or quiz-style response      |
| `vocabulary` | Explain important words or phrases                     |

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

### Required Variables

| Variable          | Description                              |
| ----------------- | ---------------------------------------- |
| `PORT`            | Port where the server runs               |
| `API_KEY`         | API key required for protected routes    |
| `OPENAI_API_KEY`  | Your OpenAI API key                      |
| `OPENAI_MODEL`    | Chat model used for student responses    |
| `EMBEDDING_MODEL` | Embedding model used for document search |
| `DATABASE_PATH`   | SQLite database path                     |
| `GIN_MODE`        | Gin mode, usually `debug` or `release`   |

---

## Installation

Clone or create the project folder:

```bash
mkdir learning-insight-coach
cd learning-insight-coach
```

Install dependencies:

```bash
go mod tidy
```

Create your `.env` file:

```bash
cp .env.example .env
```

Add your OpenAI API key to `.env`.

---

## Running the App

Run the API locally:

```bash
go run main.go
```

Or use the Makefile:

```bash
make run
```

The server should start on:

```txt
http://localhost:8080
```

---

## Health Check

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

## API Authentication

Protected routes require an API key.

Send the API key using the `X-API-Key` header:

```txt
X-API-Key: dev-secret-change-me
```

The value must match the `API_KEY` value in your `.env` file.

---

## API Endpoints

### Health Check

```txt
GET /health
```

Checks whether the server is running.

---

### Upload Document

```txt
POST /api/v1/documents
```

Uploads and indexes lesson content.

#### Request

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

#### Example Response

```json
{
  "document_id": "generated-document-id",
  "status": "indexed",
  "chunk_count": 1,
  "message": "document uploaded and indexed"
}
```

Save the returned `document_id`. You will need it when asking reader questions.

---

### Reader Response

```txt
POST /api/v1/reader/respond
```

Generates a student-friendly response using the uploaded lesson content.

#### Request

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

#### Example Response

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

### Teacher Insights

```txt
GET /api/v1/teacher/classes/:class_id/insights
```

Returns class-level insight based on saved student interactions.

#### Request

```bash
curl http://localhost:8080/api/v1/teacher/classes/class-1/insights \
  -H "X-API-Key: dev-secret-change-me"
```

#### Example Response

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

### Run Evals

```txt
POST /api/v1/evals/run
```

Runs AI behavior test cases from the eval file.

#### Request

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

## Running Tests

Run all tests:

```bash
go test ./... -v -count=1
```

Or use the Makefile:

```bash
make test
```

---

## Formatting Code

```bash
gofmt -w .
```

Or:

```bash
make lint
```

---

## Building the App

```bash
go build -o bin/coach main.go
```

Or:

```bash
make build
```

---

## Running with Docker

Build and run with Docker Compose:

```bash
docker compose up --build
```

The app will be available at:

```txt
http://localhost:8080
```

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

### 2. Student Reader Flow

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

### 3. Teacher Insight Flow

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

## Development Notes

### Cost Control

Each student reader request may call OpenAI multiple times:

* One embedding call for the student question
* One chat call for the student response
* One chat call for groundedness checking
* One chat call for direct answer checking

During local development, you may choose to temporarily disable some guardrail checks to reduce API usage.

---

### Recommended Models

For a cost-friendly development setup:

```env
OPENAI_MODEL=gpt-4o-mini
EMBEDDING_MODEL=text-embedding-3-small
```

---

### Common Issue: MaxTokens Error

If you see an error like:

```txt
this model is not supported MaxTokens, please use MaxCompletionTokens
```

Remove `MaxTokens` from the OpenAI chat request or use `MaxCompletionTokens` for models that require it.

For `gpt-4o-mini`, you can usually omit the token limit while developing.

---

## Future Improvements

Possible improvements for this app:

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

## Summary

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

This makes it more than a chatbot. It is a learning support system that can help both students and teachers.
