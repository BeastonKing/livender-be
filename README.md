# Instalasi
1. Clone repository
2. Buat file .env di root folder, isi dengan
```
DB_HOST="localhost"
DB_USERNAME=<username_db_anda>
DB_PASSWORD=<password_db_anda>
DB_DATABASENAME="livender"
DB_PORT=5432
```
3. Buat sebuah database dengan nama livender
4. Install package dengan `go mod tidy`
5. `go run main.go` untuk menjalankan aplikasi

# Testing
1. `go test` untuk menjalankan test