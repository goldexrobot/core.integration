# Goldex Robot: Callbacks

Goldex sends HTTP requests to the business backend: **callbacks** and custom **UI methods**.

Exact endpoints to call by Goldex should be defined in Goldex dashboard.

Requests are always of method **POST** and carry **application/json; charset=utf-8** with headers:

| Header | Meaning | Example |
| --- | --- | --- |
| X-CBOT-PROJECT-ID | Origin project ID | "1" |
| X-CBOT-BOT-ID | Origin bot ID (uint64) | "42" |

All the requests are [signed](/SIGNATURE.md) by Goldex (check out dashboard for a verification public key) and could be optionally transferred using mutual TLS.

---

## Callbacks

Callbacks are sent on common events defined by Goldex.

For instance: bot status changes, items evaluation, storage interaction, etc.

Business backend have to respond with successful HTTP status (200, 201, or 202) to signalize about callback consumption. Until then Goldex will continue to send callback requests.

Callback models are described in [swagger](https://goldexrobot.github.io/core.integration/swagger/#/backend-callbacks)

### Evaluation callbacks

| Callback | Description | Model | Inside | Expected response |
| --- | --- | --- | --- | --- |
| Started | New item evaluation is created | `EvalStarted` | Evaluation ID | Status 200 |
| Finished | Item evaluation is successfully finished | `EvalFinished` | Detailed evaluation data, photo file IDs | Status 200 |

### Storage callbacks

| Callback | Description | Model | Inside | Expected response |
| --- | --- | --- | --- | --- |
| Cell occupation | Cell is occupied with some item | `StorageCellEvent` | Cell address and *domain* | Status 200 |
| Cell release | Cell is freed and item is released | `StorageCellEvent` | Cell address and *domain* | Status 200 |
| Pre-occupation | Optional, called before occupation to give more flexibility | `StorageCellEvent` | Cell address and *domain* | Status 200 and `{"allowed":true}` |
| Pre-release | Optional, called before release to give more flexibility | `StorageCellEvent` | Cell address and *domain* | Status 200 and `{"allowed":true}` |

*Domain* is an origin of the cell operation:
| Domain | Whats |
| --- | --- |
| buyout | Cell is occupied with bought item |
| shop | Cell is freed during item purchasing |
| pawnshop | Cell is occupied or freed during pawnshop flow |
| collection | Cell is occupied/freed by staff members (on-bot collection dashboard) |
| dashboard | Cell is occupied/freed by staff members (on-bot system dashboard) |
| other | For custom UI flows not defined above |

---

## UI methods

UI methods are custom named callbacks and available to call directly from the physical Goldex terminal.
In this case Goldex backend acts as secure proxy, for instance to authenticate a terminal on your business backend side.

For instance, let's assume you defined a method `auth` bound to your `https://example.com/bot/auth`.
A terminal now calls `auth` with some payload. Goldex backend adds bot/project IDs to the request, signs it and and sends `POST https://example.com/bot/auth`:

```json
{
  "project_id": 1,
  "bot_id": 42,
  "payload": {
    // bot`s request k/v goes here:
    "foo": {
      "bar": "baz"
    }
  }
}
```

Now Goldex backend waits for a response from your backend (status 200, application/json) and re-sends the response back to the terminal.
