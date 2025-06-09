# Task management

## Authors

- [@nguyenducloc](https://github.com/nguyen-duc-loc)

## Run locally with docker compose

Clone the project

```bash
  git clone https://github.com/nguyen-duc-loc/task-management
```

Go to the project directory

```bash
  cd task-management
```

Go to `backend` directory, change environment variables example file `.env.example` to `.env`. Then edit your environment variables like this:

```
  SERVER_PORT=8080
  SERVER_ENV=dev
  DB_HOST=postgres
  DB_PORT=5432
  DB_DATABASE=task-management
  DB_USERNAME=root
  DB_PASSWORD=secret
  DB_SCHEMA=public
  JWT_SECRET_KEY=5f3a670f4f875d9a918310e1a54dce1d7ba1407dd40c23c5200c8c639900cb75
  JWT_ACCESS_TOKEN_DURATION=720h
```

Copy this `.env` file

After that, go to the root directory, paste the copied `.env` file to the root directory and run the following command to run docker compose:

```bash
  docker compose up
```

Open [http://localhost](http://localhost) with your browser to see the result.

## Deploy to AWS with Terraform

Clone the project

```bash
  git clone https://github.com/nguyen-duc-loc/task-management
```

Go to the project directory

```bash
  cd task-management
```

Build the frontend image

```bash
  docker build -t task-management/frontend .
```

Then, change directory to `backend` folder and build the backend image

```bash
  cd backend
  docker build -t task-management/backend .
```

Go to AWS Console, create two repository, name them `task-management/frontend` and `task-management/backend`

Tag the created images and push them to ECR:

```bash
  docker tag task-management/frontend:latest <YOUR_AWS_ACCOUNT_ID>.dkr.ecr.<YOUR_AWS_ACCOUNT_REGION>.amazonaws.com/task-management/frontend:latest
  docker push <YOUR_AWS_ACCOUNT_ID>.dkr.ecr.<YOUR_AWS_ACCOUNT_REGION>.amazonaws.com/task-management/frontend:latest

  docker tag task-management/backend:latest <YOUR_AWS_ACCOUNT_ID>.dkr.ecr.<YOUR_AWS_ACCOUNT_REGION>.amazonaws.com/task-management/backend:latest
  docker push <YOUR_AWS_ACCOUNT_ID>.dkr.ecr.<YOUR_AWS_ACCOUNT_REGION>.amazonaws.com/task-management/backend:latest
```

Go to `infras` folder, create `terraform.tfvars`. Add the variables to this file

```terraform
  app_name = "task-management"
  jwt_secret_key = "<YOUR_JWT_SECRET_KEY>"
  backend_image = "<YOUR_AWS_ACCOUNT_ID>.dkr.ecr.<YOUR_AWS_ACCOUNT_REGION>.amazonaws.com/task-management/backend:latest"
  frontend_image = "<YOUR_AWS_ACCOUNT_ID>.dkr.ecr.<YOUR_AWS_ACCOUNT_REGION>.amazonaws.com/task-management/frontend:latest"
```

Init provider, then plan and apply terraform:

```bash
  terraform init -upgrade
  terraform plan
  terraform apply --auto-approve
```

Go to the EC2 section > Load Balancers. Find the load balancer of your frontend service and try to access its DNS.
