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

## Notes
- Create a new user first before creating a new loan.
- To create a new loan, the user must not be delinquent.
- To make a payment, the user ID and loan ID must match each other.
- The loan term is weekly.
- If a user doesn't manage to pay the loan within 14 days, the user will be considered as delinquent.
- User must pay the loan with the same amount every week. Late payments require users to pay all late payments in the same week (if the bill amount is 5000 per week, then 2 weeks late will require users to pay 10000 in total).
- Delinquent users cannot create new loans before finishing all of their outstanding loans.
- After a user finished all of its loans, the user will be considered as non-delinquent and can create new loans again.

## How to run

1. Install Go - (any version would be fine because this project uses basic modules)
2. Clone the repository
3. Run `make deps` to fetch dependencies
4. Run `make run`

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

Response sample:

```
{
    "message": "success to create user"
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

Response sample:

```
{
    "message": "success to create loan"
}
```

### Get the outstanding amount of a loan

```
GET /get_outstanding?loan_id=1
```

Response sample:

```
{
    "data": {
        "loan_id": 1,
        "outstanding": 5000000
    }
}
```

### Get the next payment date and amount

```
GET /get_payment_schedule?loan_id=1
```

Response sample:

```
{
    "data": {
        "next_payment_date": "2024-06-23",
        "payment_amount": 110000,
        "is_late_payment": false
    }
}
```

### Get the delinquent status of a user

```
GET /is_delinquent?user_id=1
```

Response sample:

```
{
    "data": {
        "user_id": 1,
        "is_delinquent": false
    }
}

```

### Get user list

```
GET /user_list
```

Response sample:

```
{
    "data": [
        {
            "id": 1,
            "name": "Andi",
            "is_delinquent": false
        },
        {
            "id": 2,
            "name": "Budi",
            "is_delinquent": false
        }
    ]
}
```

### Get loan list

```
GET /loan_list
```

Response sample:

```
{
    "data": [
        {
            "id": 1,
            "user_id": 1,
            "amount": 5000000,
            "term": 50,
            "bill_amount": 110000,
            "outstanding": 5000000,
            "create_time": "2024-06-16T15:12:27.073779+07:00",
            "update_time": "2024-06-16T15:12:27.073779+07:00"
        },
        {
            "id": 11,
            "user_id": 2,
            "amount": 1000000,
            "term": 5,
            "bill_amount": 220000,
            "outstanding": 0,
            "create_time": "2024-06-16T16:33:58.228969+07:00",
            "update_time": "2024-10-22T00:00:00Z"
        }
    ]
}
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

Response sample:

```
{
    "data": {
        "loan_id": 1,
        "is_success": true,
        "message": "success to make a payment"
    }
}
```

```
{
    "data": {
        "loan_id": 1,
        "is_success": false,
        "message": "amount is not equal to the bill amount. Expected: 110000.00, Actual: 1000.00"
    }
}
```