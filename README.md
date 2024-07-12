## TODO

- remove password from DSN Postgres
- use error code for error
- add email functionality
  - add email field for users
  - c. once invested all investors will receive an email containing link to agreement letter (pdf)
  - test with yopmail
- add S3 functionality for file sharing

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

- implement user authentication
  - store salted and hashed password in users table when sign up
  - use JWT for sign in
  - in loan endpoint requests, get userID from JWT directly instead of from HTTP request
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
- Create end-to-end tests to make sure:
  - http request parsing works correctly
  - integration with database/object storage/email client works correctly