# `dto/responses.go`

Named response envelopes so Swagger documents accurate schemas instead of generic objects.

## Types

`UserEnvelope{Message, User}`, `UsersEnvelope{Message, Users}`, `ProfileEnvelope{Message, Profile}`, `ContactEnvelope{Message, Contact}`, `ContactsEnvelope{Message, Contacts, Meta (PageMeta)}`, `MessageEnvelope{Message}`, `ErrorEnvelope{Error}`.

Referenced by handler `@Success`/`@Failure` annotations (e.g. `{object} dto.AuthResponse`, `{object} dto.ErrorEnvelope`).
