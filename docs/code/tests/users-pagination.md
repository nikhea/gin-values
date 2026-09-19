# `tests/users_pagination_test.go` — users list tests

- `TestListUsersSearch` — name/email substring match, case-insensitivity, no-match empty.
- `TestListUsersPages` — page 1/2 slicing with stable totals, out-of-range page returns rows with full total.
- `TestListUsersHTTP` — `GET /users/?search=&page=&page_size=` returns rows + `meta` envelope; non-numeric page → 400.
