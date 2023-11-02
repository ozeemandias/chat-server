-- +goose Up
-- +goose StatementBegin

-- CreateTable
CREATE TABLE "messages"
(
    "id"         SERIAL       NOT NULL,
    "chat_id"    INTEGER      NOT NULL,
    "user_id"    INTEGER      NOT NULL,
    "message"    TEXT         NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "messages_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "messages_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "chats" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX "messages_user_id_idx" ON "messages" ("user_id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "messages";
-- +goose StatementEnd
