# Docker Compose Setup Guide

## Quick Start

1. **Make sure Docker is running**
```bash
docker --version
docker-compose --version
```

2. **Start the services**
```bash
docker-compose up
```

Or in detached mode:
```bash
docker-compose up -d
```

3. **Check logs**
```bash
docker-compose logs -f app
docker-compose logs -f db
```

4. **Stop the services**
```bash
docker-compose down
```

## Configuration

### Environment Variables

The `docker-compose.yml` file includes default environment variables:

- **DB_HOST**: `db` (service name - DO NOT change to localhost)
- **DB_USER**: `postgres`
- **DB_PASS**: `postgres`
- **DB_NAME**: `yourapp`
- **DB_PORT**: `5432`
- **SERVER_PORT**: `8080`

### Optional: Custom .env File

You can create a `.env` file in the project root to override defaults:

```env
# JWT Secrets (important for production)
JWT_SECRET=your-production-secret-key
JWT_REFRESH=your-production-refresh-key

# Server
SERVER_PORT=8080

# Other custom variables
```

**Important**: Even if you have a `.env` file, `DB_HOST` in docker-compose will use `db` (the service name).

## Common Issues

### Issue: Connection Refused Error

**Problem**: `dial tcp 127.0.0.1:5432: connect: connection refused`

**Solution**: Make sure `DB_HOST=db` (not `localhost`) in docker-compose.yml

### Issue: App Starts Before Database is Ready

**Solution**: The docker-compose.yml now includes a healthcheck that waits for the database to be ready.

### Issue: Port Already in Use

**Problem**: Port 5432 or 8080 already in use

**Solution**: 
```bash
# Change ports in docker-compose.yml
ports:
  - "5433:5432"  # Use different host port
  - "8081:8080"  # Use different host port
```

## Database Management

### Run Migrations

```bash
# Option 1: Inside the app container
docker-compose exec app migrate -path migrations -database "postgres://postgres:postgres@db:5432/yourapp?sslmode=disable" up

# Option 2: Using psql directly
docker-compose exec db psql -U postgres -d yourapp
```

### Seed Database

```bash
docker-compose exec app go run cmd/seed/main.go
```

### Access PostgreSQL

```bash
docker-compose exec db psql -U postgres -d yourapp
```

### Backup Database

```bash
docker-compose exec db pg_dump -U postgres yourapp > backup.sql
```

### Restore Database

```bash
docker-compose exec -T db psql -U postgres yourapp < backup.sql
```

## Troubleshooting

### View Logs

```bash
# All services
docker-compose logs

# Specific service
docker-compose logs app
docker-compose logs db

# Follow logs
docker-compose logs -f app
```

### Restart Services

```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart app
docker-compose restart db
```

### Rebuild After Code Changes

```bash
# Rebuild and restart
docker-compose up --build

# Rebuild specific service
docker-compose up --build app
```

### Clean Everything

```bash
# Stop and remove containers
docker-compose down

# Remove volumes (WARNING: deletes database data)
docker-compose down -v

# Remove everything including images
docker-compose down --rmi all -v
```

### Check Service Status

```bash
docker-compose ps
```

### Execute Commands in Container

```bash
# App container
docker-compose exec app sh
docker-compose exec app go run cmd/seed/main.go

# DB container
docker-compose exec db psql -U postgres
```

## Health Checks

### Test API Health

```bash
curl http://localhost:8080/api/v1/health
```

Should return: `OK`

### Test Database Connection

```bash
docker-compose exec app sh -c 'psql -h db -U postgres -d yourapp -c "SELECT version();"'
```

## Production Notes

For production:
1. Change default passwords
2. Use strong JWT secrets
3. Enable SSL for database
4. Use environment variables from secure storage
5. Configure proper network isolation
6. Set up monitoring and logging
7. Use secrets management (AWS Secrets Manager, etc.)

## Volumes

The docker-compose.yml creates a named volume `postgres_data` to persist database data. This means your data survives container restarts.

To backup:
```bash
docker run --rm -v campus-ambassador-backend_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/db-backup.tar.gz /data
```
