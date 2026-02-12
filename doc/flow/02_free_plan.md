Kita desain dengan asumsi:

- ✅ Multi-tenant shared DB
- ✅ Ada Free Plan (1 bulan trial)
- ✅ Bisa upgrade ke Standard/Premium
- ✅ Midtrans untuk pembayaran
- ✅ Tenant tidak langsung “mati”, tapi bisa suspended

---

# 🎯 FLOW YANG KITA SUPPORT

```
Register sekolah
    ↓
Tenant langsung aktif
    ↓
Free plan 30 hari
    ↓
Sebelum habis → bisa upgrade
    ↓
Jika tidak bayar → suspend
```

---

# 📦 1️⃣ MASTER SAAS TABLES (GLOBAL)

## 1. subscription_plans

Daftar semua plan (Free, Standard, Premium)

```sql
subscription_plans
- id (uuid, pk)
- name (Free, Standard, Premium)
- price_monthly
- price_yearly
- duration_days (30 untuk Free)
- max_students
- max_teachers
- features_json (jsonb)
- is_active
- created_at
```

👉 Free plan:

- price = 0
- duration_days = 30

---

## 2️⃣ REGISTRATION PROCESS TABLES

Karena kamu pakai free trial langsung aktif,
kita TIDAK perlu tenant_registrations terpisah.

Kita bisa langsung create tenant + subscription.

---

# 📦 3️⃣ tenants (SEKOLAH)

```sql
tenants
- id (uuid, pk)
- name
- slug (unique)
- email
- phone
- address
- logo_url
- status (active, suspended, cancelled)
- trial_ends_at (nullable)
- created_at
- updated_at
```

### status logic:

| status    | artinya                   |
| --------- | ------------------------- |
| active    | bisa akses                |
| suspended | trial habis / belum bayar |
| cancelled | berhenti total            |

---

# 📦 4️⃣ tenant_subscriptions

Ini penting untuk billing lifecycle.

```sql
tenant_subscriptions
- id (uuid, pk)
- tenant_id (fk)
- plan_id (fk)
- billing_cycle (monthly, yearly)
- price
- status (trialing, active, past_due, cancelled)
- started_at
- ended_at
- next_billing_at
- midtrans_subscription_id (nullable)
- created_at
```

---

## Status lifecycle

| status    | kondisi      |
| --------- | ------------ |
| trialing  | free 30 hari |
| active    | sudah bayar  |
| past_due  | gagal bayar  |
| cancelled | stop         |

---

# 📦 5️⃣ users (Admin dibuat saat register)

```sql
users
- id
- tenant_id
- name
- email
- password_hash
- role (admin)
- is_active
- created_at
```

Saat registrasi:

- create tenant
- create tenant_subscription (Free trial)
- create admin user

---

# 📦 6️⃣ invoices (Untuk upgrade nanti)

```sql
invoices
- id
- tenant_id
- subscription_id
- code
- amount
- status (pending, paid, expired)
- due_date
- created_at
```

---

# 📦 7️⃣ payments (Midtrans tracking)

```sql
payments
- id
- tenant_id
- invoice_id
- midtrans_order_id
- midtrans_transaction_id
- amount
- status (pending, settlement, expire, cancel)
- paid_at
- raw_response_json (jsonb)
```

---

# 📦 8️⃣ webhook_logs (WAJIB untuk debugging)

```sql
webhook_logs
- id
- provider (midtrans)
- payload_json (jsonb)
- signature_valid (boolean)
- processed (boolean)
- created_at
```

Sangat penting untuk production debugging.

---

# 🧠 MEKANISME FREE PLAN 1 BULAN

Saat register:

### 1️⃣ Insert tenant

```
status = active
trial_ends_at = now() + 30 days
```

### 2️⃣ Insert tenant_subscription

```
plan_id = Free
status = trialing
started_at = now()
ended_at = now() + 30 days
next_billing_at = now() + 30 days
```

---

# 🕒 CRON JOB (WAJIB)

Bikin scheduled job harian:

Pseudo logic:

```sql
SELECT * FROM tenant_subscriptions
WHERE status = 'trialing'
AND ended_at < now()
```

Jika ditemukan:

```
UPDATE tenants SET status = 'suspended'
UPDATE tenant_subscriptions SET status = 'past_due'
```

---

# 🔥 Saat Upgrade ke Paid Plan

Flow:

```
Create invoice
→ Generate Midtrans
→ Webhook success
→ Update subscription:
    status = active
    plan_id = premium
    next_billing_at = now() + 1 month
```

---

# 📊 Ringkasan FINAL Table yang Dibutuhkan

Untuk registrasi + free trial 1 bulan:

1. ✅ subscription_plans
2. ✅ tenants
3. ✅ tenant_subscriptions
4. ✅ users
5. ✅ invoices
6. ✅ payments
7. ✅ webhook_logs

Total: **7 tabel inti untuk SaaS lifecycle**

---

# 🧠 Kenapa Tidak Perlu tenant_registrations?

Karena kamu pakai:

> Free plan langsung aktif

Kalau model bayar dulu baru aktif → baru perlu tenant_registrations.

---

# 🚀 Rekomendasi Profesional untuk Kelasmirai

Karena kamu mau ini jadi SaaS serius:

Tambahkan juga:

```sql
tenant_usage_logs
- tenant_id
- student_count
- teacher_count
- checked_at
```

Untuk enforce limit plan.

---

Kalau kamu mau, next kita bisa bahas:

- 🔥 Desain limit enforcement (misal max 500 siswa)
- 🔥 Billing lifecycle lengkap (upgrade, downgrade, prorate)
- 🔥 Cara suspend tenant tapi tetap simpan data
- 🔥 Atau desain sistem auto email reminder trial

Kita bikin Kelasmirai SaaS level enterprise dari awal 👌
