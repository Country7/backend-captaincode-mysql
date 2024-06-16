CREATE TABLE `users` (
  `username` varchar(255) PRIMARY KEY,
  `role` varchar(255) NOT NULL DEFAULT 'depositor',
  `hashed_password` varchar(255) NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `is_email_verified` bool NOT NULL DEFAULT false,
  `password_changed_at` timestamptz NOT NULL DEFAULT '0001-01-01',
  `created_at` timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE `accounts` (
  `id` bigserial PRIMARY KEY,
  `owner` varchar(255) NOT NULL,
  `balance` bigint NOT NULL,
  `currency` varchar(255) NOT NULL,
  `created_at` timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE `entries` (
  `id` bigserial PRIMARY KEY,
  `account_id` bigint NOT NULL,
  `amount` bigint NOT NULL COMMENT 'can be negative or positive',
  `created_at` timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE `transfers` (
  `id` bigserial PRIMARY KEY,
  `from_account_id` bigint NOT NULL,
  `to_account_id` bigint NOT NULL,
  `amount` bigint NOT NULL COMMENT 'must be positive',
  `created_at` timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE `sessions` (
  `id` uuid PRIMARY KEY,
  `username` varchar(255) NOT NULL,
  `refresh_token` varchar(255) NOT NULL,
  `user_agent` varchar(255) NOT NULL,
  `client_ip` varchar(255) NOT NULL,
  `is_blocked` boolean NOT NULL DEFAULT false,
  `expires_at` timestamptz NOT NULL,
  `created_at` timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX `accounts_index_0` ON `accounts` (`owner`);

CREATE UNIQUE INDEX `accounts_index_1` ON `accounts` (`owner`, `currency`);

CREATE INDEX `entries_index_2` ON `entries` (`account_id`);

CREATE INDEX `transfers_index_3` ON `transfers` (`from_account_id`);

CREATE INDEX `transfers_index_4` ON `transfers` (`to_account_id`);

CREATE INDEX `transfers_index_5` ON `transfers` (`from_account_id`, `to_account_id`);

ALTER TABLE `accounts` ADD FOREIGN KEY (`owner`) REFERENCES `users` (`username`);

ALTER TABLE `entries` ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `transfers` ADD FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `transfers` ADD FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `sessions` ADD FOREIGN KEY (`username`) REFERENCES `users` (`username`);
