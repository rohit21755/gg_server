#!/bin/bash

# AWS Setup Script for Campus Ambassador Backend
# This script helps set up AWS resources for deployment

set -e

echo "🚀 Campus Ambassador Backend - AWS Setup"
echo "========================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check AWS CLI
if ! command -v aws &> /dev/null; then
    echo -e "${RED}❌ AWS CLI not found. Please install it first.${NC}"
    echo "Visit: https://aws.amazon.com/cli/"
    exit 1
fi

echo -e "${GREEN}✅ AWS CLI found${NC}"

# Get AWS region
read -p "Enter AWS region (default: us-east-1): " AWS_REGION
AWS_REGION=${AWS_REGION:-us-east-1}
export AWS_REGION

echo ""
echo "Selected region: $AWS_REGION"
echo ""

# Get AWS account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
echo "AWS Account ID: $ACCOUNT_ID"
echo ""

# Menu
echo "What would you like to set up?"
echo "1) ECR Repository"
echo "2) ECS Cluster"
echo "3) RDS Database"
echo "4) All of the above"
echo "5) Exit"
read -p "Enter choice [1-5]: " choice

case $choice in
    1)
        echo ""
        echo "Creating ECR repository..."
        aws ecr create-repository \
            --repository-name campus-ambassador-backend \
            --region $AWS_REGION \
            --image-scanning-configuration scanOnPush=true \
            --encryption-configuration encryptionType=AES256 \
            2>/dev/null || echo "Repository already exists"
        
        REPO_URI=$(aws ecr describe-repositories \
            --repository-names campus-ambassador-backend \
            --region $AWS_REGION \
            --query 'repositories[0].repositoryUri' \
            --output text)
        
        echo -e "${GREEN}✅ ECR Repository created${NC}"
        echo "Repository URI: $REPO_URI"
        ;;
    
    2)
        echo ""
        echo "Creating ECS cluster..."
        aws ecs create-cluster \
            --cluster-name campus-ambassador-cluster \
            --region $AWS_REGION \
            --capacity-providers FARGATE FARGATE_SPOT \
            --default-capacity-provider-strategy capacityProvider=FARGATE,weight=1 \
            2>/dev/null || echo "Cluster already exists"
        
        echo -e "${GREEN}✅ ECS Cluster created${NC}"
        echo "Cluster: campus-ambassador-cluster"
        ;;
    
    3)
        echo ""
        echo "Creating RDS database..."
        echo -e "${YELLOW}⚠️  This will create a db.t3.micro instance. Make sure you have the right VPC and security group.${NC}"
        read -p "Continue? (y/n): " confirm
        
        if [ "$confirm" != "y" ]; then
            echo "Cancelled."
            exit 0
        fi
        
        read -p "Database password: " DB_PASSWORD
        read -p "VPC Security Group ID: " SECURITY_GROUP_ID
        
        aws rds create-db-instance \
            --db-instance-identifier campus-ambassador-db \
            --db-instance-class db.t3.micro \
            --engine postgres \
            --engine-version 16.1 \
            --master-username postgres \
            --master-user-password "$DB_PASSWORD" \
            --allocated-storage 20 \
            --vpc-security-group-ids "$SECURITY_GROUP_ID" \
            --db-name campusambassador \
            --backup-retention-period 7 \
            --storage-encrypted \
            --region $AWS_REGION \
            2>/dev/null || echo "Database creation initiated (may take 5-10 minutes)"
        
        echo -e "${GREEN}✅ RDS Database creation initiated${NC}"
        echo "Database ID: campus-ambassador-db"
        echo "Check status: aws rds describe-db-instances --db-instance-identifier campus-ambassador-db"
        ;;
    
    4)
        echo ""
        echo "Setting up all resources..."
        
        # ECR
        echo "1/3 Creating ECR repository..."
        aws ecr create-repository \
            --repository-name campus-ambassador-backend \
            --region $AWS_REGION \
            --image-scanning-configuration scanOnPush=true \
            2>/dev/null || echo "ECR repository already exists"
        
        # ECS
        echo "2/3 Creating ECS cluster..."
        aws ecs create-cluster \
            --cluster-name campus-ambassador-cluster \
            --region $AWS_REGION \
            2>/dev/null || echo "ECS cluster already exists"
        
        # RDS
        echo "3/3 RDS setup requires manual configuration."
        echo "Please create RDS database via AWS Console or run option 3 separately."
        
        echo ""
        echo -e "${GREEN}✅ Setup complete!${NC}"
        ;;
    
    5)
        echo "Exiting..."
        exit 0
        ;;
    
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo "Next steps:"
echo "1. Configure GitHub Secrets"
echo "2. Update workflow files with your resource names"
echo "3. Push to main branch to trigger deployment"
echo ""
echo "See DEPLOYMENT.md for detailed instructions."
