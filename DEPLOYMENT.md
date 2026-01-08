# AWS Deployment Guide

This guide covers deploying the Campus Ambassador Backend to AWS using GitHub Actions.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [AWS Setup](#aws-setup)
3. [GitHub Secrets Configuration](#github-secrets-configuration)
4. [Deployment Options](#deployment-options)
5. [Post-Deployment](#post-deployment)
6. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Accounts & Tools

- [ ] AWS Account with appropriate permissions
- [ ] GitHub repository
- [ ] AWS CLI installed locally (for setup)
- [ ] Docker installed (for local testing)

### AWS Services You'll Use

- **AWS ECR** (Elastic Container Registry) - Store Docker images
- **AWS ECS** (Elastic Container Service) - Run containers (Option 1)
- **AWS EC2** - Virtual server (Option 2)
- **AWS App Runner** - Serverless containers (Option 3)
- **AWS RDS** - PostgreSQL database
- **AWS Secrets Manager** (optional) - Store sensitive data

## AWS Setup

### Step 1: Create AWS IAM User

1. Go to AWS Console → IAM → Users → Add users
2. Create a user named `github-actions-deploy`
3. Attach policies:
   - `AmazonEC2ContainerRegistryFullAccess`
   - `AmazonECS_FullAccess` (if using ECS)
   - `AmazonEC2FullAccess` (if using EC2)
   - `AWSAppRunnerFullAccess` (if using App Runner)
   - `AmazonRDSFullAccess` (for database access)
4. Create Access Key:
   - Go to Security credentials tab
   - Create access key
   - **Save the Access Key ID and Secret Access Key** (you'll need these for GitHub secrets)

### Step 2: Create RDS PostgreSQL Database

```bash
# Using AWS CLI
aws rds create-db-instance \
  --db-instance-identifier campus-ambassador-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --engine-version 16.1 \
  --master-username postgres \
  --master-user-password YourSecurePassword123! \
  --allocated-storage 20 \
  --vpc-security-group-ids sg-xxxxxxxxx \
  --db-name campusambassador \
  --backup-retention-period 7 \
  --storage-encrypted
```

**Note**: Replace `sg-xxxxxxxxx` with your security group ID.

**Important**: Save the database endpoint and credentials!

### Step 3: Choose Your Deployment Option

## Deployment Options

### Option 1: AWS ECS (Recommended for Production)

**Best for**: Production applications, auto-scaling, high availability

#### Setup Steps:

1. **Create ECR Repository**:
```bash
aws ecr create-repository \
  --repository-name campus-ambassador-backend \
  --region us-east-1
```

2. **Create ECS Cluster**:
```bash
aws ecs create-cluster \
  --cluster-name campus-ambassador-cluster \
  --region us-east-1
```

3. **Create Task Definition**:
   - Go to ECS Console → Task Definitions → Create new
   - Use the JSON below or create via console

```json
{
  "family": "campus-ambassador-task",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "containerDefinitions": [
    {
      "name": "campus-ambassador-backend",
      "image": "YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/campus-ambassador-backend:latest",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "SERVER_PORT",
          "value": "8080"
        }
      ],
      "secrets": [
        {
          "name": "DATABASE_URL",
          "valueFrom": "arn:aws:secretsmanager:us-east-1:ACCOUNT_ID:secret:db-credentials"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/campus-ambassador",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

4. **Create ECS Service**:
   - Go to ECS Console → Clusters → Your cluster → Services → Create
   - Select your task definition
   - Configure load balancer (Application Load Balancer recommended)
   - Set desired count to 2+ for high availability

5. **Configure GitHub Secrets** (see below)

### Option 2: AWS EC2 (Simpler, More Manual)

**Best for**: Development, testing, simpler deployments

#### Setup Steps:

1. **Launch EC2 Instance**:
   - Go to EC2 Console → Launch Instance
   - Choose Ubuntu 22.04 LTS
   - Instance type: t3.small or larger
   - Configure security group:
     - SSH (22) from your IP
     - HTTP (80) from anywhere
     - Custom TCP (8080) from anywhere
   - Create/select key pair
   - Launch instance

2. **Setup EC2 Instance**:
```bash
# SSH into your instance
ssh -i your-key.pem ubuntu@your-ec2-ip

# Install dependencies
sudo apt update
sudo apt install -y postgresql-client

# Create app directory
mkdir -p ~/app
cd ~/app

# Create systemd service file
sudo nano /etc/systemd/system/campus-ambassador.service
```

3. **Create systemd Service File**:
```ini
[Unit]
Description=Campus Ambassador Backend
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/app
ExecStart=/home/ubuntu/app/bin/server
Restart=always
RestartSec=10
EnvironmentFile=/home/ubuntu/app/.env

[Install]
WantedBy=multi-user.target
```

4. **Setup Nginx Reverse Proxy** (optional but recommended):
```bash
sudo apt install -y nginx
sudo nano /etc/nginx/sites-available/campus-ambassador
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/campus-ambassador /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

5. **Configure GitHub Secrets** (see below)

### Option 3: AWS App Runner (Serverless)

**Best for**: Simple deployments, automatic scaling, managed service

#### Setup Steps:

1. **Create ECR Repository** (same as ECS):
```bash
aws ecr create-repository \
  --repository-name campus-ambassador-backend \
  --region us-east-1
```

2. **Create App Runner Service**:
   - Go to App Runner Console → Create service
   - Source: Container registry → ECR
   - Select your ECR repository
   - Service name: `campus-ambassador-app`
   - Port: 8080
   - Environment variables: Add your required env vars
   - Create service

3. **Configure GitHub Secrets** (see below)

## GitHub Secrets Configuration

Go to your GitHub repository → Settings → Secrets and variables → Actions → New repository secret

### Required Secrets (All Options):

| Secret Name | Description | Example |
|------------|-------------|---------|
| `AWS_ACCESS_KEY_ID` | AWS IAM user access key | `AKIAIOSFODNN7EXAMPLE` |
| `AWS_SECRET_ACCESS_KEY` | AWS IAM user secret key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@host:5432/dbname` |

### ECS-Specific Secrets:

| Secret Name | Description |
|------------|-------------|
| `ECS_CLUSTER` | Your ECS cluster name (or set in workflow env) |
| `ECS_SERVICE` | Your ECS service name (or set in workflow env) |
| `ECS_TASK_DEFINITION` | Your task definition name (or set in workflow env) |

### EC2-Specific Secrets:

| Secret Name | Description |
|------------|-------------|
| `EC2_HOST` | EC2 instance public IP or domain |
| `EC2_USER` | SSH username (usually `ubuntu`) |
| `EC2_SSH_KEY` | Private SSH key content |
| `EC2_PORT` | SSH port (default: 22) |
| `EC2_INSTANCE_ID` | EC2 instance ID (optional) |

### App Runner-Specific Secrets:

| Secret Name | Description |
|------------|-------------|
| `APP_RUNNER_SERVICE_ARN` | App Runner service ARN |

### How to Get SSH Key for EC2:

```bash
# If you don't have the key, you'll need to:
# 1. Create a new key pair in EC2 console
# 2. Download the .pem file
# 3. Convert to format for GitHub (remove headers/footers, keep only key content)
cat your-key.pem
# Copy the entire content (including BEGIN/END lines) to GitHub secret
```

## Environment Variables

Create a `.env` file or set environment variables with:

```env
# Server
SERVER_PORT=8080

# Database
DB_HOST=your-rds-endpoint.region.rds.amazonaws.com
DB_USER=postgres
DB_PASS=YourSecurePassword123!
DB_NAME=campusambassador
DB_PORT=5432

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this
JWT_REFRESH=your-super-secret-refresh-key-change-this

# Optional
ENVIRONMENT=production
LOG_LEVEL=info
```

**For ECS**: Store in AWS Secrets Manager or ECS Task Definition environment variables.

**For EC2**: Create `.env` file on the server.

**For App Runner**: Add in App Runner service configuration.

## Post-Deployment

### 1. Run Database Migrations

**For ECS**:
```bash
# Get running task
TASK_ARN=$(aws ecs list-tasks --cluster campus-ambassador-cluster --query 'taskArns[0]' --output text)

# Execute migration command
aws ecs execute-command \
  --cluster campus-ambassador-cluster \
  --task $TASK_ARN \
  --container campus-ambassador-backend \
  --interactive \
  --command "/bin/sh"
# Then inside container: migrate -path migrations -database "$DATABASE_URL" up
```

**For EC2**:
```bash
ssh -i your-key.pem ubuntu@your-ec2-ip
cd ~/app
# Install migrate tool first
# Then: migrate -path migrations -database "$DATABASE_URL" up
```

**For App Runner**: Run migrations manually or add to startup script.

### 2. Seed Database (Optional)

```bash
# SSH/connect to your deployment
# Run: go run cmd/seed/main.go
# Or: make seed
```

### 3. Health Check

```bash
# Test health endpoint
curl https://your-domain.com/api/v1/health

# Should return: OK
```

### 4. Test API

```bash
# Test login
curl -X POST https://your-domain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@campusambassador.com","password":"password123"}'
```

## Workflow Files

The repository includes three workflow files:

1. **`.github/workflows/deploy-ecs.yml`** - ECS deployment
2. **`.github/workflows/deploy-ec2.yml`** - EC2 deployment  
3. **`.github/workflows/deploy-app-runner.yml`** - App Runner deployment

**To use a specific workflow**:
- Enable the workflow file you want
- Disable or delete the others
- Configure the appropriate secrets

## Troubleshooting

### Build Fails

- Check Go version matches (1.24.3)
- Verify all dependencies in `go.mod`
- Check GitHub Actions logs

### Deployment Fails

**ECS**:
- Verify ECR repository exists
- Check task definition is correct
- Verify IAM permissions
- Check ECS service logs in CloudWatch

**EC2**:
- Verify SSH key is correct
- Check security group allows SSH
- Verify EC2 instance is running
- Check application logs: `journalctl -u campus-ambassador -f`

**App Runner**:
- Verify ECR image exists
- Check App Runner service logs
- Verify environment variables

### Database Connection Issues

- Verify RDS security group allows connections from your deployment
- Check database endpoint is correct
- Verify credentials
- Test connection: `psql -h endpoint -U user -d dbname`

### Application Not Starting

- Check environment variables
- Verify port configuration
- Check application logs
- Test locally with Docker: `docker-compose up`

## Security Best Practices

1. ✅ Use AWS Secrets Manager for sensitive data
2. ✅ Enable SSL/TLS (use AWS Certificate Manager)
3. ✅ Restrict security groups to minimum required access
4. ✅ Use IAM roles instead of access keys when possible
5. ✅ Enable CloudWatch logging
6. ✅ Set up CloudWatch alarms
7. ✅ Enable database encryption
8. ✅ Regular security updates
9. ✅ Use strong passwords and rotate regularly
10. ✅ Enable AWS WAF for API protection

## Cost Optimization

- **ECS Fargate**: Pay per use, good for variable traffic
- **EC2**: Reserved instances for predictable workloads
- **App Runner**: Pay per request, good for low traffic
- **RDS**: Use smaller instance types for development
- **CloudWatch**: Set log retention policies

## Monitoring

Set up CloudWatch:
- Application logs
- CPU/Memory metrics
- Request counts
- Error rates
- Database connections

## Next Steps

1. Set up custom domain with Route 53
2. Configure SSL certificate (ACM)
3. Set up CI/CD for staging environment
4. Configure auto-scaling
5. Set up monitoring and alerts
6. Implement backup strategy
7. Set up disaster recovery plan

## Support

For issues:
1. Check GitHub Actions logs
2. Check CloudWatch logs
3. Review this documentation
4. Check AWS service status
