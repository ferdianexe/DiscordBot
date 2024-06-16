# billing-engine

This is a simple billing engine built with Go.

## Features

- Create a new user
- Create a new loan
- Get the outstanding amount of a loan
- Get the next payment date and amount
- Get the delinquent status of a user
- Get user list
- Get loan list
- Make a payment

## How to run

1. Install Go - (any version would be fine because this project uses basic modules)
2. Clone the repository
3. Run `make run`

## How to use

### Create a new user

```
POST /create_user
Content-Type: application/json

{
    "name": "John Doe",
    "is_delinquent": false
}
```

### Create a new loan

```
POST /make_loan
Content-Type: application/json

{
    "user_id": 1,
    "amount": 1000,
    "term": 3
}
```

### Get the outstanding amount of a loan

```
GET /get_outstanding?loan_id=1
```

### Get the next payment date and amount

```
GET /get_payment_schedule?loan_id=1
```

### Get the delinquent status of a user

```
GET /is_delinquent?user_id=1
```

### Get user list

```
GET /user_list
```

### Get loan list

```
GET /loan_list
```

### Make a payment

```
POST /make_payment
Content-Type: application/json

{
    "loan_id": 1,
    "user_id": 1,
    "amount": 1000,

    // Optional if you want to override the payment date. Otherwise, the system will use the current date.
    "payment_date": "2020-01-01"
}
``` 