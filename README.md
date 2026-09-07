\# Ticket System API



A simple REST API built in Go for managing user support tickets.



\## Features



\- User registration

\- Secure password hashing using bcrypt

\- User login with JWT authentication

\- Create tickets

\- View only your own tickets

\- Get a specific ticket owned by you

\- Update ticket status

\- Enforced ticket status flow:

&#x20; - open -> in\_progress

&#x20; - in\_progress -> closed

&#x20; - closed tickets cannot be reopened



\## Tech Stack



\- Go

\- net/http

\- JWT

\- bcrypt

\- In-memory storage

\- Docker



\## API Endpoints



| Method | Endpoint | Authentication |

|---|---|---|

| GET | `/health` | No |

| POST | `/auth/register` | No |

| POST | `/auth/login` | No |

| POST | `/tickets` | JWT |

| GET | `/tickets` | JWT |

| GET | `/tickets/{id}` | JWT |

| PATCH | `/tickets/{id}/status` | JWT |



\## Authentication



Protected endpoints require:



```text

Authorization: Bearer <JWT>

