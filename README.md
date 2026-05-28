# Runbox

A self-hostable AI agent that writes and executes code in a secure sandboxed environment. Give it a task, it figures out how to solve it.

## What it does

You send Runbox a task. It spins up an AI agent that writes Python code, runs it in an isolated container, sees the output, and iterates until the task is done. If it doesn't need code, it just answers directly.

Think of it as a lightweight, self-hostable version of ChatGPT's Code Interpreter — except you own the entire stack.

## How it works

**The Sandbox** — every code execution runs in a fresh Docker container with hard limits: 128MB memory, 50% CPU, no network access, 10 second timeout. The container is destroyed after each run. Nothing persists, nothing escapes.

**The Agent Loop** — the LLM writes code wrapped in `<code>` tags. Runbox runs it in the sandbox, feeds the output back to the LLM, and repeats until it responds with `<done>` and the final answer. Max 5 iterations.

**The API** — two endpoints. One returns the final result as JSON. The other streams the agent's progress in real time using SSE so you can watch it think step by step.

## API

**POST** `/v1/task` — run a task, get back the result

```bash
curl -X POST http://localhost:8080/v1/task \
  -H "Content-Type: application/json" \
  -d '{"task": "Calculate the sum of all even numbers between 1 and 100"}'
```

```json
{"result": "2550", "status": "success"}
```

**POST** `/v1/stream/task` — same thing but streams progress in real time

```bash
curl -X POST http://localhost:8080/v1/stream/task \
  -H "Content-Type: application/json" \
  -d '{"task": "Calculate the sum of all even numbers between 1 and 100"}'
```

```
data: <code>...python code...</code>
data: 2550
{"result": "2550", "status": "success"}
```

## Stack

- **Go** — API server and agent orchestration
- **Docker SDK** — programmatic container lifecycle management
- **Groq (Llama 3.3 70B)** — LLM backend
- **SSE** — real time streaming

## Running locally

**Prerequisites:** Go, Docker

```bash
git clone https://github.com/HavishBoda/runbox
cd runbox
```

Create a `.env` file:
```
GROQ_API_KEY=your_key_here
```

```bash
go run main.go llm.go agent.go api.go
```

## What's next

- Multi-language support (JS, Go)
- Kubernetes orchestration for concurrent executions at scale
- Persistent task history
- Simple UI to submit tasks and watch the agent work
