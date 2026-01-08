# Deployment Checklist - Step by Step

Follow this checklist to deploy your application to AWS.

## 📋 Pre-Deployment Checklist

### 1. AWS Account Setup
- [ ] Create AWS account (if you don't have one)
- [ ] Sign in to AWS Console
- [ ] Note your AWS Account ID
- [ ] Choose your preferred AWS region (e.g., us-east-1)

### 2. Create IAM User for GitHub Actions
- [ ] Go to IAM → Users → Add users
- [ ] Username: `github-actions-deploy`
- [ ] Attach policies:
  - [ ] `AmazonEC2ContainerRegistryFullAccess`
  - [ ] `AmazonECS_FullAccess` (if using ECS)
  - [ ] `AmazonEC2FullAccess` (if using EC2)
  - [ ] `AWSAppRunnerFullAccess` (if using App Runner)
  - [ ] `AmazonRDSFullAccess`
- [ ] Create Access Key
- [ ] **SAVE**: Access Key ID and Secret Access Key (you'll need these!)

### 3. Create RDS Database
- [ ] Go to RDS → Create database
- [ ] Engine: PostgreSQL 16.1
- [ ] Template: Free tier (or Production)
- [ ] DB instance identifier: `campus-ambassador-db`
- [ ] Master username: `postgres`
- [ ] Master password: **SAVE THIS!**
- [ ] DB instance class: `db.t3.micro` (or larger)
- [ ] Storage: 20 GB
- [ ] VPC: Default or your VPC
- [ ] Security group: Create new or use existing
- [ ] **SAVE**: Database endpoint, username, password

### 4. Choose Deployment Method

**Select ONE option:**

#### Option A: ECS (Recommended for Production)
- [ ] Create ECR repository: `campus-ambassador-backend`
- [ ] Create ECS cluster: `campus-ambassador-cluster`
- [ ] Create task definition (see DEPLOYMENT.md)
- [ ] Create ECS service (see DEPLOYMENT.md)
- [ ] Configure load balancer (optional but recommended)

#### Option B: EC2 (Simpler)
- [ ] Launch EC2 instance (Ubuntu 22.04, t3.small+)
- [ ] Configure security group (SSH, HTTP, port 8080)
- [ ] Create/Download SSH key pair
- [ ] Note EC2 public IP or domain
- [ ] SSH into instance and set up (see DEPLOYMENT.md)

#### Option C: App Runner (Easiest)
- [ ] Create ECR repository: `campus-ambassador-backend`
- [ ] Create App Runner service via AWS Console
- [ ] Note App Runner service ARN

## 🔐 GitHub Secrets Configuration

Go to: **GitHub Repo → Settings → Secrets and variables → Actions → New repository secret**

### Required for ALL Options:
- [ ] `AWS_ACCESS_KEY_ID` = (from IAM user)
- [ ] `AWS_SECRET_ACCESS_KEY` = (from IAM user)
- [ ] `DATABASE_URL` = `postgres://postgres:password@db-endpoint:5432/campusambassador`

### ECS Additional:
- [ ] (Workflow uses env vars, but verify they match your resources)

### EC2 Additional:
- [ ] `EC2_HOST` = (EC2 public IP or domain)
- [ ] `EC2_USER` = `ubuntu` (or your username)
- [ ] `EC2_SSH_KEY` = (full content of your .pem file, including BEGIN/END lines)
- [ ] `EC2_PORT` = `22` (optional, defaults to 22)

### App Runner Additional:
- [ ] `APP_RUNNER_SERVICE_ARN` = (from App Runner console)

## ⚙️ Workflow Configuration

### Choose Your Workflow

**Enable ONE workflow file:**
- [ ] `.github/workflows/deploy-ecs.yml` (for ECS)
- [ ] `.github/workflows/deploy-ec2.yml` (for EC2)
- [ ] `.github/workflows/deploy-app-runner.yml` (for App Runner)

**Disable the others** (or delete them if you're sure)

### Update Environment Variables

Edit your chosen workflow file and update these if needed:

```yaml
env:
  AWS_REGION: us-east-1  # Change if different
  ECR_REPOSITORY: campus-ambassador-backend  # Match your ECR repo
  ECS_CLUSTER: campus-ambassador-cluster  # Match your cluster
  # ... etc
```

## 🚀 Deployment Steps

### 1. Commit and Push
- [ ] Commit all changes
- [ ] Push to `main` or `master` branch
- [ ] GitHub Actions will automatically start

### 2. Monitor Deployment
- [ ] Go to GitHub → Actions tab
- [ ] Watch the workflow run
- [ ] Check for any errors

### 3. Verify Deployment
- [ ] Wait for deployment to complete (5-10 minutes)
- [ ] Check AWS Console:
  - [ ] ECS: Service is running
  - [ ] EC2: Application is running
  - [ ] App Runner: Service is active

## 🗄️ Post-Deployment

### 1. Run Database Migrations

**For ECS:**
```bash
# Get task ARN
TASK_ARN=$(aws ecs list-tasks --cluster campus-ambassador-cluster --query 'taskArns[0]' --output text)

# Execute command
aws ecs execute-command \
  --cluster campus-ambassador-cluster \
  --task $TASK_ARN \
  --container campus-ambassador-backend \
  --interactive \
  --command "/bin/sh"
```

**For EC2:**
```bash
ssh -i your-key.pem ubuntu@your-ec2-ip
cd ~/app
# Install migrate tool, then:
# migrate -path migrations -database "$DATABASE_URL" up
```

**For App Runner:**
- Run migrations manually or add to startup script

### 2. Seed Database (Optional)
- [ ] Connect to your deployment
- [ ] Run: `go run cmd/seed/main.go` or `make seed`

### 3. Health Check
- [ ] Test: `curl https://your-domain.com/api/v1/health`
- [ ] Should return: `OK`

### 4. Test API
- [ ] Test login endpoint
- [ ] Test other endpoints
- [ ] Verify database connection

## ✅ Verification

- [ ] Health endpoint returns `OK`
- [ ] API endpoints are accessible
- [ ] Database connection works
- [ ] Logs are visible in CloudWatch
- [ ] Application is stable

## 🔧 Troubleshooting

If something fails:

1. **Check GitHub Actions Logs**
   - [ ] Go to Actions tab
   - [ ] Click on failed workflow
   - [ ] Review error messages

2. **Check AWS Console**
   - [ ] Verify resources exist
   - [ ] Check IAM permissions
   - [ ] Verify security groups

3. **Check Application Logs**
   - [ ] CloudWatch Logs (ECS/App Runner)
   - [ ] EC2: `journalctl -u campus-ambassador -f`

4. **Common Issues**
   - [ ] Wrong AWS credentials → Check GitHub secrets
   - [ ] Database connection failed → Check security groups
   - [ ] Build failed → Check Go version
   - [ ] Deployment timeout → Check resource availability

## 📚 Reference Documents

- **QUICK_DEPLOY.md** - Quick 5-step guide
- **DEPLOYMENT.md** - Complete detailed guide
- **DEPLOYMENT_SUMMARY.md** - Overview of what was created

## 🎉 Success!

Once all checkboxes are checked, your application is deployed! 🚀

---

**Need help?** Check the detailed guides or AWS documentation.
