# Travel Advisor

## Getting Started

Follow these steps to set up and run the application:

```bash
# 1. Start dependencies (database, cache)
docker compose up

# 2. Run database migrations to initialize schema and districts data migration
go run . migration

# 3. Start background  scheduler
go run . scheduler

# 4. Launch the main server
go run . serve
