# `tests/notifications_test.go` — inbox tests

- `TestNotifyWorkerWritesRow` — enqueue → River worker writes the row with title/body/data intact, unread.
- `TestInboxFlow` — HTTP register + HTTP verify emits the welcome end-to-end (worker path); list/unread-count agree; mark-one-read; cross-user read 404 + empty inbox; read-all idempotent with final count 0.
- `waitForNotification` polls for the async worker (10s timeout).
