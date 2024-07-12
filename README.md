## TODO

- add email functionality
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

Seed (password is Test123!):
```
INSERT INTO public.users
(id, "role", "name", email, hashed_password)
VALUES
('borrower_1', 'BORROWER', 'Borrower 1', 'amartha_test_borrower_1@yopmail.com', '$2a$14$jMgHEV/EpcqRKwsFqt/07.U7zWYBCkUL/xt5/AHhgaytt4XAQGjCC')
('investor_1', 'INVESTOR', 'Investor 1', 'amartha_test_investor_1@yopmail.com', '$2a$14$jMgHEV/EpcqRKwsFqt/07.U7zWYBCkUL/xt5/AHhgaytt4XAQGjCC')
('field_officer_1', 'FIELD_OFFICER', 'Field Officer 1', 'amartha_test_field_officer_1@yopmail.com', '$2a$14$jMgHEV/EpcqRKwsFqt/07.U7zWYBCkUL/xt5/AHhgaytt4XAQGjCC');
```

## Possible Improvements

- make S3 bucket private
  - investors can only view loan agreement letter pdf after login and creating presigned url
- use error code for error
  - example: user/not-found, loan/not-found, user/unauthorized-role, etc.
- sign up flow:
  - implement OTP for sign up email verification
  - restrict field officer role to be added not by sign up but by internal operation team.
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