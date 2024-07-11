##

- remove password from DSN Postgres

## Run in Local

- go mod tidy
- go run main.go

## Run all tests

- go test ./...

## Setup

- go mod init amartha/loan-service
- go get github.com/gin-gonic/gin gorm.io/gorm

Seed:
```
INSERT INTO public.users
(id, "role", "name")
values
('borrower_1', 'BORROWER', 'Borrower 1'),
('investor_1', 'INVESTOR', 'Investor 1'),
('field_officer_1', 'FIELD_OFFICER', 'Field Officer 1');
```

## Possible Improvements

- Use LoanStatusHistory table.
  - From status
  - To status
  - Timestamp
  - Detail json (picture proof, signed agreement, etc.)