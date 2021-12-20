# Goldex Robot: Integration

Integration with Goldex consists of two major parts: backend integration and terminal integration.

Documentation is [here](https://goldexrobot.github.io/core.integration/).

---

## Backend

### Callbacks

Goldex sends [callbacks](/CALLBACK.md) to notify your backend about a new events or to request some information for the vending terminal in real time.

Navigate to Goldex dashbaord to setup callbacks, keys etc.

[Documentation](https://goldexrobot.github.io/core.integration/swagger/#/backend-callbacks)

### API

Goldex exposes an API to provide some extended information like photos, storage access history etc.
Moreover the API allows you to control a vending terminal.

Calls to HTTP API must be properly [signed](/SIGNATURE.md) with per-project private key. You can get the key in Goldex dashboard.

GRPC API is also available.

---

## Terminal

TODO