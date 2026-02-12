Kita fokus hanya ke tabel yang dibutuhkan untuk:

> ✅ Registrasi sekolah
> ✅ Free plan 30 hari
> ✅ Upgrade ke paid
> ✅ Suspend jika trial habis
> ✅ Midtrans integration

DB: **PostgreSQL 15+ (recommended)**
PK: **UUID**
JSON: **JSONB**

---

# 🧱 1️⃣ subscription_plans (GLOBAL)

```sql
CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    price_monthly NUMERIC(15,2) NOT NULL DEFAULT 0,
    price_yearly NUMERIC(15,2) NOT NULL DEFAULT 0,
    duration_days INTEGER NOT NULL, -- 30 untuk Free
    max_students INTEGER,
    max_teachers INTEGER,
    features_json JSONB,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

Contoh data:

- Free → price 0, duration_days 30
- Standard → duration_days 30
- Premium → duration_days 30

---

# 🏫 2️⃣ tenants

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(150) NOT NULL,
    phone VARCHAR(30),
    address TEXT,
    logo_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    trial_ends_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### status allowed values:

- active
- suspended
- cancelled

👉 Sebaiknya pakai CHECK:

```sql
ALTER TABLE tenants
ADD CONSTRAINT tenants_status_check
CHECK (status IN ('active','suspended','cancelled'));
```

---

# 💳 3️⃣ tenant_subscriptions

Ini pusat lifecycle billing.

```sql
CREATE TABLE tenant_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES subscription_plans(id),
    billing_cycle VARCHAR(20) NOT NULL,
    price NUMERIC(15,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at TIMESTAMP WITH TIME ZONE NOT NULL,
    next_billing_at TIMESTAMP WITH TIME ZONE,
    midtrans_subscription_id VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

Constraint:

```sql
ALTER TABLE tenant_subscriptions
ADD CONSTRAINT subscription_status_check
CHECK (status IN ('trialing','active','past_due','cancelled'));

ALTER TABLE tenant_subscriptions
ADD CONSTRAINT billing_cycle_check
CHECK (billing_cycle IN ('monthly','yearly'));
```

Index penting:

```sql
CREATE INDEX idx_subscription_tenant
ON tenant_subscriptions(tenant_id);
```

---

# 👤 4️⃣ users (Admin dibuat saat register)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (tenant_id, email)
);
```

Constraint:

```sql
ALTER TABLE users
ADD CONSTRAINT role_check
CHECK (role IN ('admin','teacher','student','parent','staff'));
```

Index penting:

```sql
CREATE INDEX idx_users_tenant ON users(tenant_id);
```

---

# 🧾 5️⃣ invoices

```sql
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    subscription_id UUID REFERENCES tenant_subscriptions(id),
    code VARCHAR(50) NOT NULL UNIQUE,
    amount NUMERIC(15,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

Constraint:

```sql
ALTER TABLE invoices
ADD CONSTRAINT invoice_status_check
CHECK (status IN ('pending','paid','expired','cancelled'));
```

---

# 💰 6️⃣ payments

```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    midtrans_order_id VARCHAR(100) NOT NULL,
    midtrans_transaction_id VARCHAR(100),
    amount NUMERIC(15,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    paid_at TIMESTAMP WITH TIME ZONE,
    raw_response_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

Constraint:

```sql
ALTER TABLE payments
ADD CONSTRAINT payment_status_check
CHECK (status IN ('pending','settlement','expire','cancel','deny'));
```

Index penting:

```sql
CREATE INDEX idx_payments_tenant ON payments(tenant_id);
CREATE INDEX idx_payments_order ON payments(midtrans_order_id);
```

---

# 🛰 7️⃣ webhook_logs

```sql
CREATE TABLE webhook_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(50) NOT NULL,
    payload_json JSONB NOT NULL,
    signature_valid BOOLEAN,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

# 🧠 FREE PLAN 30 HARI — IMPLEMENTASI NYATA

Saat register:

### 1️⃣ Insert tenant

```sql
status = 'active'
trial_ends_at = NOW() + INTERVAL '30 days'
```

### 2️⃣ Insert tenant_subscription

```sql
status = 'trialing'
started_at = NOW()
ended_at = NOW() + INTERVAL '30 days'
price = 0
billing_cycle = 'monthly'
```

---

# 🕒 CRON JOB (WAJIB)

Daily job:

```sql
UPDATE tenants
SET status = 'suspended'
WHERE id IN (
    SELECT tenant_id
    FROM tenant_subscriptions
    WHERE status = 'trialing'
    AND ended_at < NOW()
);
```

---

# 📌 TOTAL TABEL YANG TERLIBAT

| Table                | Scope  |
| -------------------- | ------ |
| subscription_plans   | global |
| tenants              | global |
| tenant_subscriptions | global |
| users                | tenant |
| invoices             | tenant |
| payments             | tenant |
| webhook_logs         | global |

Total: **7 tabel inti SaaS lifecycle**

---

# 🚀 Ini Sudah Production-Ready?

Untuk MVP SaaS?
✅ YES

Untuk scale 10.000 sekolah?
Masih perlu:

- Usage tracking table
- Feature flag table
- Billing history table

---
