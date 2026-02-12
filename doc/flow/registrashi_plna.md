Kita bedakan dulu 2 jenis registrasi:

1. 🏢 **Registrasi Sekolah (Tenant Signup)** → dilakukan Admin Sekolah
2. 👤 Registrasi internal user (guru/siswa/parent) → dibuat oleh admin sekolah

Yang kita bahas sekarang: **registrasi tenant + pilih plan + bayar**.

---

# 🎯 Goal Flow yang Ideal (Production Grade)

Target akhir:

```
Admin daftar → pilih plan → bayar → tenant dibuat → subdomain aktif → login
```

Bukan:

```
Tenant langsung dibuat sebelum bayar ❌
```

Karena itu bisa bikin banyak tenant "sampah".

---

# 🧠 High-Level Flow

## STEP 1 — Admin Isi Form Registrasi

Endpoint:

```
POST /public/tenant/register
```

Payload:

```json
{
  "school_name": "SMK Negeri 1 Jakarta",
  "slug": "smkn1-jkt",
  "admin_name": "Budi",
  "admin_email": "admin@smkn1.sch.id",
  "password": "secret",
  "plan_id": "uuid-plan-premium",
  "billing_cycle": "monthly"
}
```

---

## STEP 2 — Backend Validasi

Validasi:

- slug unik
- email belum terpakai
- plan valid
- password kuat

Belum buat tenant permanen dulu.

Buat:

```
tenant_registrations (temporary)
```

---

# 📦 Table: tenant_registrations

```sql
- id
- school_name
- slug
- admin_name
- admin_email
- password_hash
- plan_id
- billing_cycle
- status (pending_payment, paid, expired)
- created_at
```

---

## STEP 3 — Generate Midtrans Transaction

Backend:

- hitung harga dari plan
- generate Midtrans Snap transaction
- simpan `midtrans_transaction_id`
- return `snap_token`

Response:

```json
{
  "payment_url": "...",
  "snap_token": "xxx"
}
```

Frontend redirect ke Snap.

---

# 💰 STEP 4 — Midtrans Webhook

Endpoint:

```
POST /webhooks/midtrans
```

Ketika payment SUCCESS:

Backend:

1. Verify signature key
2. Ambil registration berdasarkan order_id
3. Jika paid → lanjut proses provisioning tenant

---

# 🚀 STEP 5 — Provision Tenant (Ini Critical)

Setelah payment sukses:

### 1️⃣ Insert ke `tenants`

```sql
tenants:
- id
- name
- slug
- status = active
```

### 2️⃣ Insert ke `tenant_subscriptions`

### 3️⃣ Create Admin User

Insert ke:

```
users (role = admin)
```

### 4️⃣ (Optional) Seed Default Data

- Academic year default
- Roles default
- Permissions
- Fee types default

---

# 🔥 Baru di sini tenant resmi aktif

User sekarang bisa akses:

```
https://smkn1-jkt.kelasmirai.com
```

---

# 🧠 Flow Diagram (Ringkas)

```
[Register Form]
        ↓
[Create tenant_registration]
        ↓
[Create Midtrans transaction]
        ↓
[User bayar]
        ↓
[Midtrans webhook]
        ↓
[Provision Tenant]
        ↓
[Send Welcome Email]
        ↓
[Login]
```

---

# 🛡️ Important Security Points

## 1️⃣ Jangan aktifkan tenant sebelum webhook confirm

Jangan percaya redirect frontend.

Wajib tunggu webhook.

---

## 2️⃣ Expire pending registration

Misalnya:

```
if created_at > 1 hour AND belum bayar
→ status = expired
```

Biar database bersih.

---

## 3️⃣ Slug locking

Saat registration pending:

Slug harus dianggap "reserved".

Jangan boleh dipakai user lain.

---

# 💡 Alternatif Flow (Trial Mode)

Kalau kamu mau freemium:

Flow bisa jadi:

```
Register → Tenant langsung aktif → Free plan 14 hari → upgrade nanti
```

Table tambahan:

```sql
tenants.trial_ends_at
```

---

# 🧠 Mana yang Lebih Cocok untuk Kelasmirai?

Karena target kamu:

> All-in-one School Operating System (komersial)

Aku sarankan:

### ✅ Trial 14 hari tanpa kartu

Lebih gampang closing sekolah.

Flow:

```
Register → langsung aktif → trial 14 hari
→ kirim reminder email H-3
→ jika tidak bayar → suspend
```

Lebih conversion friendly.

---

# ⚙️ Endpoint Summary

Public:

```
POST   /public/register-school
POST   /public/check-slug
GET    /public/plans
```

Webhook:

```
POST   /webhooks/midtrans
```

Internal:

```
POST   /internal/provision-tenant
```

---

# 🧠 Di Kode Golang (Struktur Service)

Pisahkan service:

```
AuthService
TenantService
BillingService
ProvisionService
WebhookService
```

Jangan campur semua logic di handler.

---
