# Goldex Robot: Callbacks

Goldex sends HTTP requests to the partner backend: **callbacks** and custom **UI methods**.

Exact endpoints to call should be defined in Goldex dashboard.

Requests are always of method **POST** and carry **application/json**.

Requests are [signed](/SIGNATURE.md) by Goldex (check out dashboard for a verification public key) and could be optionally transferred using mutual TLS.

## Callbacks

Callbacks are sent on common events defined by Goldex.

For instance: items evaluation, storage interaction, etc.

Partner backend have to respond with successful HTTP status (200, 201, or 202) to signalize about callback consumption. Until then Goldex will continue to send callback requests.

Callback models are described in [swagger](https://goldexrobot.github.io/core.integration/swagger/#/backend-callbacks)

## UI methods

UI methods are sort of callbacks defined by partner to access own resources from physical Goldex bot through Goldex backend. In this flow Goldex backend acts as secure proxy.

Requests from physical bot are wrapped into next structure:

```json
{
  "project_id": 1,
  "bot_id": 42,
  "payload": {
    // k/v from bot
  }
}
```

## Headers

Some headers are sent with callbacks:

| Header | Meaning | Example |
| --- | --- | --- |
| X-CBOT-PROJECT-ID | Origin project ID | "1" |
| X-CBOT-BOT-ID | Origin bot ID (uint64) | "42" |
