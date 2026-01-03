-- ==================================================================================
-- AFTERLIGHT DATABASE SCHEMA
-- Target Database: SQLite (with WAL mode recommended)
-- ==================================================================================

-- Enforce Foreign Key constraints to ensure data integrity.
-- If a User is deleted, their Vaults/Beneficiaries are automatically wiped (CASCADE).
PRAGMA foreign_keys = ON;

-- ==================================================================================
-- 1. USERS
-- The core account table. Manages authentication and the "Dead Man's Switch" timer.
-- ==================================================================================
CREATE TABLE IF NOT EXISTS users (
    id                   TEXT PRIMARY KEY,  -- UUID v4
    name                 TEXT NOT NULL,
    email                TEXT UNIQUE NOT NULL,
    password_hash        TEXT NOT NULL,     -- Argon2id hash of the master password

    -- Configuration: Liveness
    is_paused            BOOLEAN NOT NULL DEFAULT FALSE, -- If TRUE, timer countdown is suspended (Vacation/Coma mode)
    check_in_interval    INTEGER NOT NULL,      -- Frequency of expected pings (in Seconds)
    trigger_interval_num INTEGER NOT NULL,      -- Consecutive missed intervals required to trigger
    buffer_period        INTEGER NOT NULL,      -- Grace period after the trigger deadline before notifying verifiers
    verifier_quorum      INTEGER DEFAULT 0,     -- Consensus: How many verifiers must confirm death? (M-of-N)

    -- State Machine
    last_check_in        DATETIME NOT NULL,
    current_status       TEXT DEFAULT 'ALIVE',  -- Enum: 'ALIVE', 'WARNING', 'VERIFICATION_REQUIRED', 'CONFIRMED_DEAD'

    created_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==================================================================================
-- 2. CONTACT METHODS
-- Stores communication channels for BOTH Users (for pings) and Beneficiaries (for delivery).
-- Ideally, a row links to EITHER a user_id OR a beneficiary_id (Polymorphic-ish).
-- ==================================================================================
CREATE TABLE IF NOT EXISTS contact_methods (
    id              TEXT PRIMARY KEY,

    -- Owner (One of these must be NULL)
    user_id         TEXT REFERENCES users(id) ON DELETE CASCADE,
    beneficiary_id  TEXT REFERENCES beneficiaries(id) ON DELETE CASCADE,

    -- Channel Details
    channel         TEXT NOT NULL, -- 'EMAIL', 'DISCORD_WEBHOOK', 'TELEGRAM', 'SLACK'
    target          TEXT NOT NULL, -- 'john@doe.com', 'https://discord.com/api/webhooks/...'

    -- MetaData (JSON) for extra provider-specific data (e.g., Discord Bot Token, formatting rules)
    metadata        TEXT,

    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraint: A contact method must belong to someone, but not both.
    CHECK (
        (user_id IS NOT NULL AND beneficiary_id IS NULL) OR
        (user_id IS NULL AND beneficiary_id IS NOT NULL)
    )
);

-- ==================================================================================
-- 2. VAULTS
-- Logical containers for secrets. A user can have multiple vaults (e.g. "Financial", "Socials").
-- ==================================================================================
CREATE TABLE IF NOT EXISTS vaults (
    id          TEXT PRIMARY KEY, -- UUID v4
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,    -- Friendly name (e.g. "My Bitcoin Keys")

    -- Security Metadata
    kdf_salt    TEXT NOT NULL,    -- Random salt used to derive the encryption key on the Client Side.
                                  -- The server NEVER sees the actual key, only this salt.

    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==================================================================================
-- 3. ARTIFACTS
-- The actual encrypted data blobs inside a vault.
-- ==================================================================================
CREATE TABLE IF NOT EXISTS artifacts (
    id              TEXT PRIMARY KEY, -- UUID v4
    vault_id        TEXT NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,

    message_type    TEXT NOT NULL,    -- Enum: 'TEXT_MESSAGE', 'FILE_UPLOAD', 'S3_LINK'
    encrypted_blob  BLOB NOT NULL,    -- The encrypted payload (AES-GCM ciphertext)
    iv              TEXT NOT NULL,    -- Initialization Vector (Required for AES-GCM decryption)

    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==================================================================================
-- 4. BENEFICIARIES (People)
-- Contacts who can either receive data OR verify death (or both).
-- ==================================================================================
CREATE TABLE IF NOT EXISTS beneficiaries (
    id              TEXT PRIMARY KEY, -- UUID v4
    user_id         TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,

    -- Verifier Role
    is_verifier     BOOLEAN DEFAULT FALSE, -- If TRUE, this person will be contacted to confirm death
    has_confirmed   BOOLEAN DEFAULT FALSE,
    confirmed_at    DATETIME,

    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==================================================================================
-- 5. VAULT ACCESS (Junction Table)
-- Controls "Who gets What". Maps Beneficiaries to specific Vaults.
-- ==================================================================================
CREATE TABLE IF NOT EXISTS vault_access (
    vault_id        TEXT NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    beneficiary_id  TEXT NOT NULL REFERENCES beneficiaries(id) ON DELETE CASCADE,

    hint            TEXT, -- Plaintext hint to help the beneficiary decrypt the vault.
                          -- Example: "The password is the city where we met."

    PRIMARY KEY (vault_id, beneficiary_id)
);

-- =================================================================================
-- 6. SESSIONS
-- Manages user sessions for authentication
-- =================================================================================
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
