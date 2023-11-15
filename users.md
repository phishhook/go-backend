# Users API Documentation

## User

Used to register a user.

**URL**: `/user`

**Method**: `POST`

**Auth required** YES

**Data Constrains**

```json
{
  "username": "no constraints",
  "phone_number": "no dashes or spaces, add in the country calling code (i.e, prefix of 1 for USA)."
}
```

**Data Example**

```json
{
  "username": "test_user",
  "phone_number": "12527778888"
}
```

**Success Response**

**Code**: 200 OK

```json
{
  "id": 7
}
```

## Users

Used to get all users who are registered.

**URL**: `/users`

**Method**: `GET`

**Auth required** YES

**Success Response**

**Code**: 200 OK

```json
[
    {
        "id": 1,
        "username": "test_hunter",
        "phone_number": "16512529620",
        "created_at": "2023-11-13T22:58:30.191423Z"
    },
    ...
]
```

## User by `PhoneNumber`

Used to get information on a specific user from their phone number.

**URL**: `/users/{phone_number}`

**Method**: `GET`

**Auth required** YES

**Success Response**

**Code**: 200 OK

```json
{
  "id": 1,
  "username": "test_hunter",
  "phone_number": "16512529620",
  "created_at": "2023-11-13T22:58:30.191423Z"
}
```
