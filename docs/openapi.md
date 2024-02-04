1. ENDPOINT /withdraw

Request params:
  - sender_wallet_id: uuid
  - receiver_wallet_id: uuid
  - currency: string
  - amount: float
  
Response params:
  - message: string
  - id: uuid


2. ENDPOINT /balance

Request params:
  - receiver_wallet_id: uuid
  - currency: string

Response params:
  - actual: float
  - frozen: float