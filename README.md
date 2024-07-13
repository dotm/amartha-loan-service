## TODO

- Create end-to-end tests to make sure:
  - http request parsing works correctly
  - integration with database/object storage/email client works correctly
  - example Golang end-to-end tests that I've created for another project:
    - https://github.com/dotm/diskon-hunter-price-monitoring/tree/master/app-be-e2e-test

## Run in Local

- go mod tidy
- go run main.go

To upload to S3, you can use the upload-to-s3.html file in a browser.

## Run all unit tests

- go test ./...

## Test Manually

- You can use the exported Postman json: amartha-loan-service-postman-v2.1.postman_collection.json
  - Import this and manually test it with Postman

## Setup

- go mod init amartha/loan-service
- go get github.com/gin-gonic/gin gorm.io/gorm

Seed (password is Test123!):
```
INSERT INTO public.users
(id, "role", "name", email, hashed_password)
VALUES
('borrower_1', 'BORROWER', 'Borrower 1', 'amartha_test_borrower_1@yopmail.com', '$2a$14$jMgHEV/EpcqRKwsFqt/07.U7zWYBCkUL/xt5/AHhgaytt4XAQGjCC'),
('investor_1', 'INVESTOR', 'Investor 1', 'amartha_test_investor_1@yopmail.com', '$2a$14$jMgHEV/EpcqRKwsFqt/07.U7zWYBCkUL/xt5/AHhgaytt4XAQGjCC'),
('investor_2', 'INVESTOR', 'Investor 2', 'amartha_test_investor_2@yopmail.com', '$2a$14$jMgHEV/EpcqRKwsFqt/07.U7zWYBCkUL/xt5/AHhgaytt4XAQGjCC'),
('field_officer_1', 'FIELD_OFFICER', 'Field Officer 1', 'amartha_test_field_officer_1@yopmail.com', '$2a$14$jMgHEV/EpcqRKwsFqt/07.U7zWYBCkUL/xt5/AHhgaytt4XAQGjCC');
```

AWS Setup:
- Create AWS account
- Create SES verified identity for sending email
- Create S3 bucket
  - In Permissions, Edit CORS to:
```
[
  {
    "AllowedHeaders": ["*"],
    "AllowedMethods": ["POST","PUT"],
    "AllowedOrigins": ["*"],
    "ExposeHeaders": []
  }
]
```

## Possible Improvements

- allow investor to update investment amount in a loan
- make S3 bucket private
  - investors can only view loan agreement letter pdf after login and creating presigned url
- use error code for error
  - example: user/not-found, loan/not-found, user/unauthorized-role, etc.
- sign up flow:
  - implement OTP for sign up email verification
  - restrict field officer role to be added not by sign up but by internal operation team.
- Cron to cancel loan if investments doesn't reach principal amount after some expiry time.
  - New loan status: expired.
- Use LoanStatusHistory table.
  - From status
  - To status
  - Timestamp
  - Detail json (picture proof, signed agreement, etc.)
- add logging for visibility and track error metrics for alert in PagerDuty