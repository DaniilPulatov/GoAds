
# _Ads Service_

![Go](https://img.shields.io/badge/Go-1.20+-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)
![REST](https://img.shields.io/badge/API-REST-brightgreen.svg)

A REST API service for managing advertisements with user authentication and admin moderation capabilities.

## Features

### User Features
- ✅ Registration and login (JWT authentication)
- ✅ Create, view, edit, and delete personal ads
- ✅ Submit ads for moderation
- ✅ Upload photos for ads
- ✅ View personal ad statistics

### Admin Features
- ✅ View all ads in the system
- ✅ Moderate ads (approve/reject)
- ✅ Delete any ads
- ✅ View system-wide statistics
- ✅ Filter ads by status

## Technical Stack

| Component               | Technology       |
|-------------------------|------------------|
| Language                | Go 1.20+         |
| Database                | PostgreSQL 15+   |
| API Framework           | Gin              |
| Dependency injection    | Dig              |
| Authentication          | JWT              |
| Password Hashing        | bcrypt           |
| Environmental variables | godotenv         |
| File Storage            | Local filesystem |
| Testing                 | testify          |
| Swagger Documentation   | go-swagger       |

## API Endpoints

### Authentication
| Method | Endpoint        | Description          |
|--------|----------------|----------------------|
| POST   | /auth/register | User registration    |
| POST   | /auth/login    | User login           |

### User Ad Endpoints
| Method | Endpoint              | Description                     |
|--------|-----------------------|---------------------------------|
| POST   | /ads                  | Create new ad (draft status)    |
| GET    | /ads                  | List user's ads                 |
| GET    | /ads/:id              | Get specific ad                 |
| PUT    | /ads/:id              | Update ad                       |
| DELETE | /ads/:id              | Delete ad                       |
| PUT    | /ads/:id/submit       | Submit ad for moderation        |
| POST   | /ads/:id/photo        | Upload photo for ad             |
| GET    | /ads/:id/photo        | Get ad photo                    |

### Admin Endpoints
| Method | Endpoint              | Description                     |
|--------|-----------------------|---------------------------------|
| GET    | /ads                  | List all ads (admin view)       |
| PUT    | /ads/:id/status       | Change ad status (moderation)   |
| DELETE | /ads/:id              | Delete any ad                   |
| GET    | /ads/stats            | Get ad statistics               |

## Getting Started

### Prerequisites
- Go 1.20+
- PostgreSQL 15+
- Docker (optional)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ads-service.git
   cd ads-service

2. Edit .env with your configuration
    ```bash
    cp .env.example .env

3. Run locally 
    ```bash
   go run cmd/adsApp/main.go

4. Run tests with make file
    ```bash
   make test

5. Or run  
    ```bash
   go test -v --cover ./...

### Curl requests
1. Registration
    ```bash
   curl -X POST http://localhost:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{
    "phone": "+998000000000",
    "password": "securepassword123",
    "first_name": "your first name", 
    "last_name": "your last name",
    }'

2. Login
    ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{
    "phone": "+998000000000",
    "password": "securepassword123"
    }'

