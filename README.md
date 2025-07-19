# Judul Proyek Anda (e.g., Djawa API Service)

[![Go Report Card](https://goreportcard.com/badge/github.com/fath-purn/go-learn)](https://goreportcard.com/report/github.com/fath-purn/go-learn)
[![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/fath-purn/go-learn/actions/workflows/go.yml/badge.svg)](https://github.com/fath-purn/go-learn/actions)

Layanan REST API yang dibangun menggunakan Go, Gin, dan GORM untuk mengelola data buku, pengguna, dan tautan pendek.

## Description

Proyek ini adalah sebuah backend service yang menyediakan API untuk operasi CRUD (Create, Read, Update, Delete) pada entitas Buku, Pengguna, dan Tautan Pendek. Dibangun dengan stack teknologi modern di ekosistem Go untuk performa yang cepat dan efisien.

**Teknologi Utama:**
- **Bahasa:** [Go](https://golang.org/)
- **Web Framework:** [Gin](https://gin-gonic.com/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** [MySQL](https://www.mysql.com/)

## Table of Contents

- [Instalasi](#instalasi)
- [Konfigurasi](#konfigurasi)
- [Penggunaan](#penggunaan)
- [Fitur](#fitur)
- [Endpoint API](#endpoint-api)
- [Berkontribusi](#berkontribusi)
- [Lisensi](#lisensi)

## Instalasi

### Prasyarat

- [Go](https://golang.org/dl/) (versi 1.18 atau lebih baru)
- [MySQL](https://dev.mysql.com/downloads/installer/) atau database kompatibel lainnya.

### Langkah-langkah

1.  **Clone repositori ini:**
    ```bash
    git clone https://github.com/fath-purn/go-learn
    ```

2.  **Masuk ke direktori proyek:**
    ```bash
    cd go-learn
    ```

3.  **Install dependensi:**
    Proyek ini menggunakan Go Modules. Dependensi akan diunduh secara otomatis saat Anda membangun atau menjalankan proyek. Untuk mengunduhnya secara manual, jalankan:
    ```bash
    go mod tidy
    ```

## Konfigurasi

1.  Pastikan server MySQL Anda berjalan.
2.  Buat sebuah database baru. Berdasarkan kode `main.go`, nama databasenya adalah `djawa`.
    ```sql
    CREATE DATABASE djawa;
    ```
3.  Koneksi database saat ini di-*hardcode* di dalam `main.go`. Sangat disarankan untuk memindahkannya ke *environment variables* untuk keamanan dan fleksibilitas.

## Penggunaan

Untuk menjalankan server pengembangan, gunakan perintah berikut dari root direktori proyek:
```bash
go run main.go
```

## Features

- ✨ **Book:** Pencatatan daftar buku
- ✅ **User:** Login user dengan jwt bearer
- 🚀 **Short URL:** Memperpendek URL

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.