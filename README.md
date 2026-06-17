# BooksLib - Automation Pipeline & Container Security

Repositori ini merupakan hasil *fork* dari proyek mikroservis [Bookslib](https://github.com/sncyber-ops/bookslib.git) yang ditujukan untuk pemenuhan tugas implementasi CI/CD (Track A). Fokus utama pada pengerjaan ini adalah membangun otomatisasi *build* dan *deployment* menggunakan Jenkins, sekaligus mengamankan lingkungan eksekusi (*runner*) melalui taktik **Just-In-Time (JIT) Socket Management**.

Seluruh infrastruktur otomasi ini dibangun secara mandiri (*self-hosted*) menggunakan kontainer Docker terisolasi.

---

## 🏗️ Alur Otomatisasi & Pengamanan Node

Pipeline ini dirancang untuk berjalan secara efisien dalam satu rangkaian *workflow* terpadu. Berbeda dengan pendekatan CI/CD konvensional yang membiarkan Docker Socket terbuka terus-menerus (berisiko tinggi terhadap eksploitasi *Privilege Escalation*), pipeline ini menerapkan sistem buka-tutup akses secara dinamis pada *runner*.

### Tahapan Eksekusi (Stages):

1. **Environment Check**
   Melakukan inspeksi awal untuk memastikan versi Docker dan Docker Compose pada *runner* siap mengeksekusi instruksi build.
   
2. **Static Application Security Testing (SAST)**
   Memanfaatkan **Trivy FS** untuk memindai berkas mentah kode sumber dan manifes dependensi secara lokal sebelum proses kompilasi kontainer dilakukan.
   
3. **Build Microservices Images (JIT Security)**
   Pada tahap ini, hak akses `/var/run/docker.sock` dibuka secara temporer (`777`) tepat saat perintah `docker compose build --no-cache` akan dieksekusi, sehingga proses perakitan *image* berjalan lancar.
   
4. **Verify Images**
   Melakukan verifikasi pasca-build menggunakan instruksi grep lokal untuk memastikan seluruh komponen *image* mikroservis (`bookslib`) telah tercipta sempurna di penyimpanan lokal.
   
5. **Deploy Application**
   Melakukan penyegaran kontainer lama dan meluncurkan arsitektur mikroservis baru ke latar belakang (`docker compose down && docker compose up -d`).

### 🔒 Post-Execution Lockdown (Mekanisme Defensif)
Melalui blok instruksi `post`, sistem dipaksa untuk **selalu** mengunci kembali hak akses Docker Socket ke mode aman (`660`) segera setelah pipeline selesai, baik dalam kondisi build berhasil maupun gagal (*Failure Lockdown*).

---

## 📦 Arsitektur Layanan & Folder Struktur

Aplikasi terintegrasi ini berjalan di atas ekosistem multi-kontainer yang terdiri dari:
*   **`books-service`**: Layanan utama yang dikembangkan dengan **Go (Golang)**.
*   **`reviews-service`**: Layanan ulasan buku yang ditenagai oleh **Python (Django)**.
*   **`frontend`**: Antarmuka web pengguna berbasis **Node.js (React)**.
*   **Database**: Penyimpanan relasional menggunakan **PostgreSQL**.

### Struktur Manajemen Environment Jenkins:
Infrastruktur otomasi dikelola secara terpisah di dalam direktori berikut:
```text
├── jenkins-docker/
│   ├── docker-compose.yml   # Konfigurasi container untuk menjalankan server Jenkins
│   └── Dockerfile.jenkins   # Custom build image Jenkins (terintegrasi CLI & perkakas pemindai)
```

---

## 🚀 Panduan Pengoperasian
**1. Menyiapkan Server Jenkins (Infrastruktur)**
Sebelum menjalankan pipeline, naikkan environment Jenkins kustom Anda terlebih dahulu melalui folder infrastruktur:

```Bash
cd jenkins-docker
docker compose up -d --build
```
Buka akses Jenkins pada peramban melalui port yang telah dikonfigurasi (misal http://localhost:8080), lalu selesaikan penyiapan awal dokumen kredensial.

**2. Menjalankan Pipeline Aplikasi**
Buat Pipeline Job baru di dasbor Jenkins Anda.

Hubungkan repositori Git proyek Bookslib ini dan arahkan Script Path ke Jenkinsfile utama di root folder.

Jalankan Build Now. Seluruh proses pengujian kode hingga deployment aplikasi mikroservis akan berjalan otomatis di dalam runner.

**3. Pengujian Mandiri secara Lokal (Tanpa Jenkins)**
Jika ingin menjalankan atau menguji fungsionalitas aplikasi secara langsung di luar ekosistem Jenkins, Anda cukup mengeksekusi perintah berikut di terminal root proyek:

```Bash
docker compose up -d --build
```

---

## 📝 Catatan Audit & Evaluasi Keamanan

### ⚠️ Temuan Kerentanan (Vulnerability Report)
Dari hasil pemindaian statis menggunakan Trivy pada tahap ke-2, ditemukan **15 celah keamanan** (2 *Critical*, 13 *High*) yang bersarang di dalam manifes dependensi `reviews-service/requirements.txt`, tepatnya pada penggunaan Django versi 4.2.7.

* **Status Saat Ini:** Pipeline sengaja dikonfigurasi dengan parameter `--exit-code 0` agar proses otomatisasi dan demonstrasi deployment aplikasi tetap dapat berlanjut hingga tahap akhir untuk keperluan visualisasi test.
* **Rekomendasi Perbaikan:** Untuk mitigasi jangka panjang, sangat disarankan melakukan pembaruan versi Django pada manifes ke *patch* aman terbaru (misalnya Django `4.2.30` atau langsung melakukan upgrade ke lini versi `6.0.x`).

### 📌 Manajemen Isu (GitHub Issues Tracking)
Sesuai dengan instruksi penanganan celah keamanan pada Track A, temuan kerentanan di atas telah didokumentasikan dan dilaporkan secara resmi melalui fitur **GitHub Issues** pada repositori ini dengan detail sebagai berikut:
* **Issue Title:** `[SECURITY] Vulnerability Found in reviews-service (Django 4.2.7)`
* **Label:** `bug` / `security`
* **Tujuan:** Memberikan visibilitas penuh kepada tim pengembang (*developer*) agar segera melakukan *patching* dependensi tanpa mengganggu kelancaran *pipeline* utama yang sedang berjalan.

### 🚀 Rencana Pengembangan ke Depan (Future Improvements)
* **Isolasi Environment Node:** Memisahkan Jenkins Controller dari Node Runner (Agent eksternal) agar eksekusi perintah Docker runtime tidak menyentuh lingkungan server utama secara langsung.
* **Penyaringan Kredensial (Secret Scanning):** Menambahkan modul pendeteksi rahasia seperti *Trufflehog* di tahap awal untuk mencegah adanya *hardcoded password* atau token API yang tidak sengaja terunggah ke repositori Git.

**Penyaringan Kredensial**
Menambahkan modul Secret Scanning khusus untuk mendeteksi potensi adanya hardcoded password atau token yang tidak sengaja terunggah ke repositori.
