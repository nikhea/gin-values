# `tests/auth_test.go` — auth flow integration tests

End-to-end coverage of `service/auth_service.go` over HTTP + services.

- `TestRegisterVerifyLoginFlow` — register → duplicate 409-equivalent → pre-verify login blocked → wrong password → bad token → link verify → login + JWT subject check.
- `TestPasswordResetFlow` — unknown-email success, token issuance, bad token, reset, old-password rejection, single-use token.
- `TestAuthHTTPAndProtectedRoutes` — uses `routes.Setup()`: public register 201 with secret-leakage assertion, 401 without token, login → `/me` 200, protected `/users/` 200 with token / 401 tampered.
- `TestLegacyUserWithoutPasswordCannotLogin` — pre-auth row can't authenticate.
- Helpers `doRequest` (JSON + optional bearer) and `registerVerifiedUser`.
