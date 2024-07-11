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

- Use JWT for user authentication
  - get userID from JWT directly instead of from HTTP request
- Cron to cancel loan if investments doesn't reach principal amount after some expiry time.
  - New loan status: expired.
- Use LoanStatusHistory table.
  - From status
  - To status
  - Timestamp
  - Detail json (picture proof, signed agreement, etc.)
- Upload VisitProofBeforeApprovalPictureUrl
  - With cloud object storage:
    - generate presigned URL to upload image and then store the link in database.
    - get image with presigned URL because it's not for public.
- Create end-to-end tests to make sure http request parsing and integration with database/object storage works correctly