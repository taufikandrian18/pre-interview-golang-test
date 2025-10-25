# pre-interview-golang-test

# 🧠 Go Programming Test
Repository ini berisi **3 soal programming test** menggunakan **bahasa Go (Golang)**.  
Setiap soal memiliki cara testing yang berbeda, seperti dijelaskan di bawah ini.

---

## 🧩 Struktur Direktori

```
.
├── soal_1/
│   └── main.go
├── soal_2/
│   └── main.go
|── soal_3/
    ├── main.go
    └── (file terkait cache)
```

--- 

## 🧮 Soal 1 — Sum Even Number

**Tujuan:**  
Menampilkan hasil kalkulasi jumlah (sum) dari angka genap pada sistem console.

**Cara Menjalankan Test:**

```bash
cd soal_1
go run main.go
```

**Hasil yang Diharapkan:**  
Program akan menampilkan output berupa hasil penjumlahan angka genap, misalnya:

```
Sum of even numbers: 30
```

---

## 🌐 Soal 2 — REST API Server

**Tujuan:**  
Mengimplementasikan REST API sederhana menggunakan Go.  
Server akan berjalan di port **5000** dan memiliki beberapa **routes** yang harus dites menggunakan Postman.

**Langkah-langkah Testing:**

1. Jalankan server:

   ```bash
   cd soal_2
   go run main.go
   ```

2. Pastikan server berjalan di:

   ```
   http://localhost:5000
   ```

3. Buka **Postman** dan lakukan pengujian pada setiap route yang telah ditentukan di file `main.go`,  
   misalnya:

   - `GET /users`
   - `POST /users`
   - `PUT /users/{id}`
   - `DELETE /users/{id}`

**Hasil yang Diharapkan:**  
Masing-masing route memberikan response JSON sesuai dengan logika yang diimplementasikan.

---

## ⚡ Soal 3 — Simple Cache System

**Tujuan:**  
Mengimplementasikan sistem **simple cache** menggunakan HTTP server di Go.

**Langkah-langkah Testing:**

1. Jalankan semua file dalam direktori:

   ```bash
   cd soal_3
   go run .
   ```

2. Buka browser dan akses:

   ```
   http://localhost:8080
   ```

3. Lakukan testing cache menggunakan **CURL** dari terminal atau menggunakan fitur bawaan browser (misalnya melalui fetch atau Postman):

   ```bash
   curl http://localhost:8080/cache?key=test
   ```

4. Uji apakah sistem cache bekerja sesuai yang diharapkan (misalnya data disimpan sementara dan bisa diambil kembali).

---## 🧮 Soal 1 — Sum Even Number

**Tujuan:**  
Menampilkan hasil kalkulasi jumlah (sum) dari angka genap pada sistem console.

**Cara Menjalankan Test:**

```bash
cd soal_1
go run main.go
```

**Hasil yang Diharapkan:**  
Program akan menampilkan output berupa hasil penjumlahan angka genap, misalnya:

```
Sum of even numbers: 30
```

---

## 🌐 Soal 2 — REST API Server

**Tujuan:**  
Mengimplementasikan REST API sederhana menggunakan Go.  
Server akan berjalan di port **5000** dan memiliki beberapa **routes** yang harus dites menggunakan Postman.

**Langkah-langkah Testing:**

1. Jalankan server:

   ```bash
   cd soal_2
   go run main.go
   ```

2. Pastikan server berjalan di:

   ```
   http://localhost:5000
   ```

3. Buka **Postman** dan lakukan pengujian pada setiap route yang telah ditentukan di file `main.go`,  
   misalnya:

   - `GET /users`
   - `POST /users`
   - `PUT /users/{id}`
   - `DELETE /users/{id}`

**Hasil yang Diharapkan:**  
Masing-masing route memberikan response JSON sesuai dengan logika yang diimplementasikan.

---

## ⚡ Soal 3 — Simple Cache System

**Tujuan:**  
Mengimplementasikan sistem **simple cache** menggunakan HTTP server di Go.

**Langkah-langkah Testing:**

1. Jalankan semua file dalam direktori:

   ```bash
   cd soal_3
   go run .
   ```

2. Buka browser dan akses:

   ```
   http://localhost:8080
   ```

3. Lakukan testing cache menggunakan **CURL** dari terminal atau menggunakan fitur bawaan browser (misalnya melalui fetch atau Postman):

   ```bash
   curl http://localhost:8080/cache?key=test
   ```

4. Uji apakah sistem cache bekerja sesuai yang diharapkan (misalnya data disimpan sementara dan bisa diambil kembali).

---

## ✅ elesai!
