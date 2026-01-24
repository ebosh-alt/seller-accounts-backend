\c postgres
CREATE EXTENSION IF NOT EXISTS dblink;
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'accounts') THEN
      PERFORM dblink_exec('dbname=postgres user=' || current_user, 'CREATE DATABASE accounts');
   END IF;
END
$$;

\c accounts
DO
$$
    BEGIN
        CREATE EXTENSION IF NOT EXISTS pgcrypto;

        CREATE TABLE categories
        (
            id   serial PRIMARY KEY,
            name text
        );

        CREATE TABLE users
        (
            id       bigint PRIMARY KEY,
            username text
        );

        CREATE TABLE sellers
        (
            id       bigint PRIMARY KEY,
            username text,
            rating   double precision,
            balance  double precision,
            wallet   text,
            login    text,
            password text
        );

        CREATE TABLE accounts
        (
            id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
            category_id integer REFERENCES categories,
            name        text,
            price       double precision,
            description text,
            accepted    boolean,
            view_type   boolean,
            created_at  timestamp with time zone DEFAULT now() NOT NULL
        );

        CREATE TABLE accounts_data
        (
            id         serial PRIMARY KEY,
            account_id uuid REFERENCES accounts,
            status     boolean,
            value      text
        );

        CREATE TABLE deals
        (
            id             serial PRIMARY KEY,
            buyer_id       bigint REFERENCES users (id),
            data_id        bigint REFERENCES accounts_data (id),
            price          double precision,
            wallet         text,
            payment_status int,
            date           date,
            guarantor      bool
        );

        CREATE TABLE chats
        (
            id        bigint PRIMARY KEY,
            user_id   bigint REFERENCES users (id),
            seller_id bigint REFERENCES sellers (id)
        );

        CREATE TABLE shops
        (
            id          serial PRIMARY KEY,
            name        text,
            description text,
            path_photo  text
        );

        CREATE INDEX IF NOT EXISTS idx_deals_buyer_id ON deals (buyer_id);
        CREATE INDEX IF NOT EXISTS idx_deals_data_id ON deals (data_id);
        CREATE INDEX IF NOT EXISTS idx_accounts_category_id ON accounts (category_id);
        CREATE INDEX IF NOT EXISTS idx_accounts_data_account_id ON accounts_data (account_id);
        CREATE INDEX IF NOT EXISTS idx_chats_user_id ON chats (user_id);
        CREATE INDEX IF NOT EXISTS idx_chats_seller_id ON chats (seller_id);

        RAISE NOTICE 'Таблицы успешно созданы.';
    EXCEPTION
        WHEN OTHERS THEN
            RAISE EXCEPTION 'Ошибка при создании таблиц: %', SQLERRM;
    END
$$;
