# Quick Deployment Checklist

## 🚀 Quick Start (5 Steps)

### Step 1: AWS Account Setup (5 minutes)

1. Create AWS account if you don't have one
2. Create IAM user with deployment permissions
3. Save Access Key ID and Secret Access Key

### Step 2: Create Database (10 minutes)

```bash
# Option A: Using AWS Console
# Go to RDS → Create database → PostgreSQL → Configure and create

# Option B: Using AWS CLI
aws rds create-db-instance \
  --db-instance-identifier campus-ambassador-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --master-username postgres \
  --master-user-password YourPassword123! \
  --allocated-storage 20
```

**Save**: Database endpoint, username, password

### Step 3: Choose Deployment Method

**Option A: ECS (Recommended)**
- Best for production
- Auto-scaling
- More setup required

**Option B: EC2 (Simplest)**
- Good for development
- Full control
- Manual scaling

**Option C: App Runner (Easiest)**
- Serverless
- Auto-scaling
- Less control

### Step 4: Configure GitHub Secrets (2 minutes)

Go to: `GitHub Repo → Settings → Secrets → Actions → New secret`

Add these secrets:

```
AWS_ACCESS_KEY_ID = your-access-key
AWS_SECRET_ACCESS_KEY = your-secret-key
DATABASE_URL = postgres://user:pass@host:5432/dbname
```

**For ECS, also add:**
- (Workflow uses environment variables, but you can override)

**For EC2, also add:**
```
EC2_HOST = your-ec2-ip-or-domain
EC2_USER = ubuntu
EC2_SSH_KEY = your-private-key-content
```

**For App Runner, also add:**
```
APP_RUNNER_SERVICE_ARN = arn:aws:apprunner:...
```

### Step 5: Deploy! (Automatic)

1. Push to `main` branch
2. GitHub Actions will automatically:
   - Run tests
   - Build Docker image
   - Deploy to AWS
3. Check deployment status in GitHub Actions tab

## 📋 Detailed Steps by Option

### ECS Deployment

1. **Create ECR Repository**:
```bash
aws ecr create-repository --repository-name campus-ambassador-backend
```

2. **Create ECS Cluster**:
```bash
aws ecs create-cluster --cluster-name campus-ambassador-cluster
```

3. **Create Task Definition** (via AWS Console or CLI)
4. **Create ECS Service** (via AWS Console)
5. **Configure GitHub Secrets**
6. **Push to main branch** → Auto-deploy!

### EC2 Deployment

1. **Launch EC2 Instance**:
   - Ubuntu 22.04
   - t3.small or larger
   - Security group: SSH (22), HTTP (80), Custom (8080)

2. **Setup EC2**:
```bash
ssh -i key.pem ubuntu@your-ec2-ip
sudo apt update
sudo apt install -y postgresql-client
mkdir -p ~/app
```

3. **Create systemd service** (see DEPLOYMENT.md)
4. **Configure GitHub Secrets**
5. **Push to main branch** → Auto-deploy!

### App Runner Deployment

1. **Create ECR Repository**:
```bash
aws ecr create-repository --repository-name campus-ambassador-backend
```

2. **Create App Runner Service** (via AWS Console)
3. **Get Service ARN** from App Runner console
4. **Configure GitHub Secrets**
5. **Push to main branch** → Auto-deploy!

## 🔧 Post-Deployment

### Run Migrations

**ECS**:
```bash
# Get task ARN
TASK_ARN=$(aws ecs list-tasks --cluster campus-ambassador-cluster --query 'taskArns[0]' --output text)

# Execute command
aws ecs execute-command --cluster campus-ambassador-cluster --task $TASK_ARN --container campus-ambassador-backend --interactive --command "/bin/sh"
```

**EC2**:
```bash
ssh -i key.pem ubuntu@your-ec2-ip
cd ~/app
# Install migrate tool, then run migrations
```

### Seed Database (Optional)

```bash
# Connect to your deployment
go run cmd/seed/main.go
```

### Test Deployment

```bash
# Health check
curl https://your-domain.com/api/v1/health

# Should return: OK
```

## 🐛 Troubleshooting

**Build fails?**
- Check Go version (1.24.3)
- Check GitHub Actions logs

**Deployment fails?**
- Check AWS credentials
- Verify IAM permissions
- Check CloudWatch logs

**App not starting?**
- Check environment variables
- Check application logs
- Verify database connection

## 📚 Full Documentation

See `DEPLOYMENT.md` for complete documentation.

## ✅ Verification Checklist

- [ ] AWS account created
- [ ] IAM user with deployment permissions
- [ ] RDS database created
- [ ] ECR/ECS/EC2/App Runner setup complete
- [ ] GitHub secrets configured
- [ ] Workflow file enabled
- [ ] Pushed to main branch
- [ ] Deployment successful
- [ ] Migrations run
- [ ] Health check passes
- [ ] API endpoints working

## 🎉 You're Done!

Your application should now be live on AWS!
