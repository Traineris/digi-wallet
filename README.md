# How to run this repo

1. Clone Repo

   - git clone github
   - cd digi-wallet

2. Install dependencies

   - go mod tidy

3. Buat Database digidb

   - gunakan postgresql

4. Jalankan migrate

   - go run main.go migrate

5. Jalankan server
   - go run main.go

# Endpoint API

- PPOST /register
- POST /login
- GET /list-user
- PUT /add-balance
- DELETE /delete-user/:id
