-- +goose Up
-- +goose StatementBegin

-- CreateTable
CREATE TABLE "chats"
(
    "id"         SERIAL       NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "chats_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "user_chats"
(
    "user_id" INTEGER NOT NULL,
    "chat_id" INTEGER NOT NULL,

    CONSTRAINT "user_chats_pkey" PRIMARY KEY ("user_id", "chat_id"),
    CONSTRAINT "user_chats_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "chats" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user_chats";
DROP TABLE "chats";
-- +goose StatementEnd
