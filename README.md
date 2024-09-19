# Project Management Application

## Overview

This application designed to streamline team collaboration and workflow organization. It provides a robust set of features for managing projects, tasks, and team members.

## Key Features

### User Management
- User registration and authentication
- Password reset functionality
- Google OAuth integration for easy sign-in

### Board Management
- Create, edit, and delete project boards
- Organize multiple tasks within boards

### Task Management
- Create, update, and delete tasks within boards
- Assign tasks to owners, managers, and employees
- Set task priorities (Low, Medium, High)
- Define planning and project due dates
- Track task progress with planning description percentages
- Update task statuses (Planning: Approved/Not Approved, Project: Working/Done/Undone)
- Add comments to tasks

### File Management
- Upload and manage planning files, project files, and planning description files for each task
- Delete individual files associated with tasks

### Team Collaboration
- Invite managers and employees to tasks
- Accept or reject task invitations
- View and manage invitations

### Reporting and Overview
- Retrieve all boards and tasks for a comprehensive project overview
- Get detailed information about specific tasks

## Technical Details

- RESTful API structure
- Cookie-based authentication for secure access to protected endpoints
- Designed to handle complex project structures with multiple user roles
- Detailed task tracking capabilities

## Target Users

This platform is ideal for businesses and teams seeking a powerful solution to:
- Manage projects efficiently
- Assign and track tasks
- Monitor progress
- Facilitate collaboration among team members in different roles

## Getting Started

### Environment Setup
This application requires a .env file in the root directory with the following environment variables:

1 . Create a file named `.env` in the root directory of the project.

2 . Copy the following content into the `.env` file:

```bash
# Database connection string (DSN)
# Example: "username:password@tcp(hostname:port)/database_name?charset=utf8mb4&parseTime=true"
DB_URL=""

# Port on which the backend application will run
# Example: "4040"
PORT="4040"

# Random character string for JWT token encryption
# Example: "audgiawudadaidiawyf123"
SECRET="audgiawudadaidiawyf123"

# AWS Region for your services
# Example: "us-west-2"
AWS_REGION=""

# Your AWS Access Key ID
# Example: "AKIAIOSFODNN7EXAMPLE"
AWS_ACCESS_KEY_ID=""

# Your AWS Secret Access Key
# Example: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
AWS_SECRET_ACCESS_KEY=""

# Your Brevo (formerly Sendinblue) username
# Example: "your_username@example.com"
BREVO_USERNAME=

# Your Brevo account password
# Example: "your_brevo_password"
BREVO_PASSWORD=

# SMTP host for Brevo
# Example: "smtp-relay.brevo.com"
SMTP_HOST=

# SMTP port for Brevo
# Example: "587"
SMPTP_PORT=

# Google Cloud Platform Client ID (from API & Services credentials)
# Example: "123456789012-abcdefghijklmnopqrstuvwxyz123456.apps.googleusercontent.com"
CLIENT_ID=""

# Google Cloud Platform Client Secret (from API & Services credentials)
# Example: "GOCSPX-ABCdefGHIjklMNOpqr123456"
CLIENT_SECRET=""

# Google Calendar API credentials in JSON format
# Example: 
GOOGLE_CALENDAR_CREDENTIALS="{\"web\":{\"client_id\":\"123456789012-abcdefghijklmnopqrstuvwxyz123456.apps.googleusercontent.com\",\"project_id\":\"your-project-id\",\"auth_uri\":\"https://accounts.google.com/o/oauth2/auth\",\"token_uri\":\"https://oauth2.googleapis.com/token\",\"auth_provider_x509_cert_url\":\"https://www.googleapis.com/oauth2/v1/certs\",\"client_secret\":\"GOCSPX-ABCdefGHIjklMNOpqr123456\",\"redirect_uris\":[\"https://www.yourdomain.com/auth/callback\"]}}"
```

Remember to replace all example values with your actual credentials and configuration details.

### Run appllication

After setting up all the environment variables, run the application by typing:

```bash
  go run main.go
```

## API Documentation

This system uses Swagger to generate comprehensive API documentation. There are two ways to access the documentation:

### Method 1: Using Local Server

1 . Set up the environment:
   - Create a `.env` file in the root directory of the project.
   - At minimum, add the following environment variable:
     ```bash
     DB_URL="username:password@tcp(hostname:port)/database_name?charset=utf8mb4&parseTime=true"
     ```
   - Replace `DB_URL` with your actual database connection string.

2 . Run the application:
```bash
  go run main.go
```

3 . Access the Swagger documentation:
   - Open your web browser
   - Navigate to: http://localhost:4040/swagger/

This will allow you to view and interact with the API documentation without needing to configure AWS, GCP, or Brevo environment settings. However, some features that depend on AWS, GCP, or Brevo services may not be fully functional when running with this minimal configuration.

### Method 2: Using Swagger Editor

1 . Go to the docs folder in the application files.

2 . Choose either the swagger.json or swagger.yaml file.

3 . Copy the contents of the chosen file.

4 . Paste the contents into the left side of the Swagger Editor at https://editor-next.swagger.io/

5 . You can view the generated API documentation on the right side of the screen.

## CI/CD and Docker
This repository includes CI/CD GitHub Actions workflows and a Dockerfile for easy deployment and containerization.

1 . GitHub Actions

Located in .github/workflows/
Automates building, testing, and deployment processes
May require configuration of GitHub secrets for credentials

2 . Dockerfile

Defines how to build a Docker image for this application

#### Note: Both CI/CD workflows and Dockerfile may need adjustments based on your specific deployment requirements.

## System Architecture

You can view the system architecture diagram of the application at this link:
https://drive.google.com/file/d/13TIeluuvUF3TMFZ44XknJHCH68avRHXu/view?usp=drive_link