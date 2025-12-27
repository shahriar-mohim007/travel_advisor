# Travel Advisor

## Getting Started

Follow these steps to set up and run the application:

```bash
# 1. Start dependencies (database, cache)
docker compose up

# 2. Install Go dependencies
go mod tidy

# 3. Run database migrations (schema + districts data)
go run . migration

# 4. Start background scheduler
go run . scheduler

# 5. Launch the main server
go run . serve
