# Link API Documentation

## Link

Used to add a link to persistent storage.

**URL**: `/link`

**Method**: `POST`

**Auth required** YES

**Data Constrains**

```json
{
  "user_id": "must be a valid user",
  "url": "none",
  "is_phishing": "{safe || phishing || indeterminate}"
}
```

**Data Example**

```json
{
  "user_id": 1,
  "url": "example.com/test",
  "is_phishing": "safe"
}
```

**Success Response**

**Code**: 200 OK

```json
{
  "id": 11
}
```

## Links

Used to get all links that are in our database.

**URL**: `/links`

**Method**: `GET`

**Auth required** YES

**Success Response**

**Code**: 200 OK

```json
[
    {
        "link_id": 1,
        "user_id": 1,
        "url": "http://allfinanciercolombiasitematrasacional.replit.app/assets/css",
        "clicked_at": "2023-11-14T01:00:51.528726Z",
        "is_phishing": "phishing"
    },
    ...
]
```

## Link by `UserId`

Fetch a link associated with a specific user.

**URL**: `/links/user/{id}`

**Method**: `GET`

**Auth required** YES

**Success Response**

**Code**: 200 OK

```json
[
    {
        "link_id": 1,
        "user_id": 1,
        "url": "http://allfinanciercolombiasitematrasacional.replit.app/assets/css",
        "clicked_at": "2023-11-14T01:00:51.528726Z",
        "is_phishing": "phishing"
    },
    ...
]
```

## Link by `LinkId`

Fetch a link that has a specific id.

**URL**: `/links/id/{id}`

**Method**: `GET`

**Auth required** YES

**Success Response**

**Code**: 200 OK

```json
{
  "link_id": 11,
  "user_id": 1,
  "url": "example.com/test",
  "clicked_at": "2023-11-15T00:31:19.483548Z",
  "is_phishing": "safe"
}
```

## Check if Link Already Exists

Check to see if a link already exists in our database.

**URL**: `/links/analyze/{link_url}`

**Method**: `GET`

**Auth required** YES

**Success Response**

**Code**: 200 OK

**Exists**

```json
{
  "link_id": 11,
  "user_id": 1,
  "url": "example.com/test",
  "clicked_at": "2023-11-15T00:31:19.483548Z",
  "is_phishing": "safe"
}
```

**Does Not Exist**

```json
{
  "error": "Failed to gather link"
}
```
