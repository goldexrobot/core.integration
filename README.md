# Goldex Robot: Integration

Goldex performs [callbacks](/CALLBACK.md) to notify about events or to request some realtime information.

Goldex [API](/API.md) exposes methods to get information on demand.

HTTP statuses threated as successful are: 200, 201, 202

We are using JWT to [sign](/SIGNATURE.md) requests to/from Goldex.

Some terms:
| Term | Description |
| --- | --- |
| API | Goldex API provides details on bots |
| Bot | Goldex Robot vending terminal |
| Callback | Event from Goldex side |
| UI methods | Custom methods for UI flow like customer identification, payments processing, etc. |
| Project | Particular integration with Goldex |

