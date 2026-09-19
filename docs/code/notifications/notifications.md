# `models/notification.go` + inbox slice + River notify worker

User inbox for event notifications, delivered asynchronously.

## Model & storage

`Notification{ID, UserID (cascade), Type, Title, Body, Data (JSONB), ReadAt?, CreatedAt}`. Types: `welcome`, `password_changed`, `import_finished`, `added_to_org`, `role_changed` (`models.Notify*` constants). Repository: `CreateNotification`, `ListNotifications` (self-scoped, `unread_only` filter, paginated, self-normalizing), `UnreadCount`, `GetNotificationByID` (owner-scoped), `MarkNotificationRead` (idempotent), `MarkAllNotificationsRead` (count).

## Delivery (`jobs/jobs.go`)

`NotifyArgs{UserID, Type, Title, Body, Data}` (`"notify"`) + `NotifyWorker` writing rows via the repository; `EnqueueNotify` helper (best-effort, callers log-and-continue). Registered alongside email workers in `Setup`.

## Emitters (all best-effort)

Welcome on HTTP verify (link + OTP), password-changed on reset, import-finished with counts, added-to-org / role-changed on member changes.

## Endpoints (`handlers/notification_handler.go`, self-scoped by JWT)

`GET /notifications/` (paginated, `?unread_only`), `GET /notifications/unread-count`, `PATCH /notifications/:id/read`, `POST /notifications/read-all`. Service layer in `service/notification_service.go` (thin passthroughs).

Tested in `tests/notifications_test.go` (worker delivery incl. data JSON, full inbox flow, cross-user 404/empty, idempotent read-all).
