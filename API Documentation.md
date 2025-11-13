# Auction House API Documentation

## Authentication

This API uses session-based authentication with CSRF protection. Most endpoints require:

- A CSRF token in the `X-CSRF-Token` header
- A valid session cookie

---

## Endpoints

### CSRF Token

#### Get CSRF Token

Retrieve a CSRF token required for protected endpoints. Store this token and include it in the `X-CSRF-Token` header for subsequent requests.

Postman post-response script example to save the token:

```js
pm.environment.set("csrf_token", pm.response.json().csrf_token);
```

**Endpoint:** `GET /v1/csrftoken`

**Success Response:** `200 OK`

```json
{
  "csrf_token": "string"
}
```

---

### Users

#### Sign Up

Create a new user account.

**Endpoint:** `POST /v1/users/signup`

**Headers:**

- `X-CSRF-Token`: Required

**Request Body:**

```json
{
  "user_name": "string",
  "password": "string",
  "email": "string",
  "bio": "string"
}
```

**Example:**

```json
{
  "user_name": "admin",
  "password": "admin123",
  "email": "admin@email.com",
  "bio": "admin"
}
```

**Success Response:** `201 Created`

```json
{
  "user_id": "string"
}
```

---

#### Login

Authenticate a user and create a session.

- Sets a `session` cookie for authenticated requests

**Endpoint:** `POST /v1/users/login`

**Headers:**

- `X-CSRF-Token`: Required

**Request Body:**

```json
{
  "email": "string",
  "password": "string"
}
```

**Example:**

```json
{
  "email": "admin@email.com",
  "password": "admin123"
}
```

**Success Response:** `200 OK`

```json
{
  "message": "logged in successfully"
}
```

---

#### Logout

End the current user session.

**Endpoint:** `POST /v1/users/logout`

**Headers:**

- `X-CSRF-Token`: Required
- `Cookie: session`: Required

**Success Response:** `200 OK`

```json
{
  "message": "logged out successfully"
}
```

---

### Products

#### Create Product

Create a new product for auction.

**Endpoint:** `POST /v1/products`

**Headers:**

- `X-CSRF-Token`: Required
- `Cookie: session`: Required

**Request Body:**

```json
{
  "product_name": "string",
  "description": "string",
  "base_price": number,
  "auction_end": "ISO 8601 datetime string"
}
```

**Example:**

```json
{
  "product_name": "Vintage Camera",
  "description": "Classic film camera in excellent condition",
  "base_price": 100.49,
  "auction_end": "2025-11-24T15:30:00.123Z"
}
```

**Success Response:** `201 Created`

```json
{
  "product_id": "string"
}
```

**Notes:**

- `auction_end` must be in ISO 8601 format with timezone
- `base_price` is the starting bid amount

---

### WebSocket

#### Subscribe to Product Auction

Real-time bidding updates for a specific product. Send messages to place bids and receive messages regarding other bids being placed.

**Protocol:** WebSocket

**Endpoint:** `WS /v1/products/ws/subscribe/{{product_id}}`

**Path Parameters:**

- `product_id`: The ID of the product to subscribe to

**Headers:**

- `X-CSRF-Token`: Required
- `Cookie: session`: Required

**Message Kind:**

- Request: PlaceBid = 0
- Ok/Success: SucceededToPlaceBid = 1
- Info: NewBidPlaced = 2, AuctionFinished = 3
- Error: FailedToPlaceBid = 4, InvalidJSON = 5

**Message Format:**

**Place bid:**

```json
{
  "kind": 0,
  "amount": 10.01
}
```

###### message to server

**Bid placed:**

```json
{
  "kind": 1,
  "amount": "10.01",
  "message": "Your bid was successfully placed",
  "user_id": "ea2fe440-7a22-4567-9ff1-9296fa76f2de"
}
```

###### message to user that placed the bid

**Bid placement broadcast:**

```json
{
  "kind": 2,
  "amount": "10.06",
  "message": "A new bid was placed",
  "user_id": "ea2fe440-7a22-4567-9ff1-9296fa76f2de"
}
```

###### message to other users

**Auction finished broadcast:**

```json
{
  "kind": 3,
  "amount": "0",
  "message": "Auction has been finished"
}
```

###### message to all users

**Bid placement errors:**

```json
{
  "kind": 4,
  "amount": "1",
  "message": "the bid value should be higher than the minimum value",
  "user_id": "ea2fe440-7a22-4567-9ff1-9296fa76f2de"
}
```

```json
{
  "kind": 4,
  "amount": "10.01",
  "message": "the bid should be higher than the previous bid",
  "user_id": "ea2fe440-7a22-4567-9ff1-9296fa76f2de"
}
```

```json
{
  "kind": 5,
  "amount": "0",
  "message": "Client not found in hashmap",
  "user_id": "ea2fe440-7a22-4567-9ff1-9296fa76f2de"
}
```

```json
{
  "kind": 5,
  "amount": "0",
  "message": "This message should be a valid JSON",
  "user_id": "ea2fe440-7a22-4567-9ff1-9296fa76f2de"
}
```

###### message to user that placed the bid

---

## Typical Workflow

1. **Get CSRF Token**

   ```
   GET /v1/csrftoken
   ```

2. **Sign Up** (if new user)

   ```
   POST /v1/users/signup
   Headers: X-CSRF-Token
   ```

3. **Login**

   ```
   POST /v1/users/login
   Headers: X-CSRF-Token
   ```

4. **Create Product** (authenticated)

   ```
   POST /v1/products
   Headers: X-CSRF-Token, Cookie: session
   ```

5. **Logout** (when done)
   ```
   POST /v1/users/logout
   Headers: X-CSRF-Token, Cookie: session
   ```

---

## Error Responses

The API returns appropriate HTTP status codes:

- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Invalid CSRF token
- `500 Internal Server Error` - Server error

---

## Notes

- Session cookies are httpOnly
- CSRF tokens must be refreshed periodically
- All POST, PUT, DELETE requests require CSRF protection
- All timestamps should be in ISO 8601 format
