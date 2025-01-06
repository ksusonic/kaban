<h1 align="center" id="title">Backend for Kabanted board</h1>

<p align="center"><img src="https://socialify.git.ci/ksusonic/kaban/image?language=1&amp;owner=1&amp;name=1&amp;stargazers=1&amp;theme=Light" alt="project-image"></p>

<p id="description">REST API service for kanban board.</p>

Scheme of database:

```mermaid
graph TD;

board_members  -->  boards;
board_members  -->  users;
lists  -->  boards;
tasks  -->  lists;
```

<h2>🛠️ Installation Steps:</h2>

<p>1. Go 1.23</p>

<p>2. Prepare environment</p>

```bash
docker compose up -d
```

<p>3. Create and edit .env file</p>

```bash
cp .env.example .env
```

<p>4. Launch service</p>

```bash
go run cmd/service/main.go
```

<p>5. Run tests (integration is optional)</p>

```bash
go test ./... --tags integration
```
