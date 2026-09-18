# `tests/users_test.go` — user service tests

`TestUserCRUDService`: create (ID assigned) → list count → get → update (name/age) → delete → 404 on deleted and missing IDs (`gorm.ErrRecordNotFound`).
