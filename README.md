# Goldex Robot: Integration

Integration with Goldex consists of two major parts: backend integration and terminal integration.

[Documentation](https://goldexrobot.github.io/core.integration/).

![Goldex environment](/docs/images/goldex_environment.png)

---

## Backend

Goldex backend communicates with a Goldex vending terminal directly over secured communication channel.

### Backend callbacks

Goldex sends [callbacks](/CALLBACK.md) to notify your backend about a new events or to request some information for the vending terminal in real time.

Navigate to Goldex dashbaord to setup callbacks, keys etc.

[Documentation](https://goldexrobot.github.io/core.integration/swagger/#/backend-callbacks)

### Backend API

Goldex exposes an API to provide some extended information like photos, storage access history etc.
Moreover the API allows you to control a vending terminal.

Calls to HTTP API must be properly [signed](/SIGNATURE.md) with per-project private key. You can get the key in Goldex dashboard.

GRPC API is also available.

---

## Terminal

Goldex terminal displays customer UI on the screen. UI is an HTML SPA (single page application) and served locally on the terminal.

The terminal exposes an API that allows a developer to access the terminal hardware, interact with a business backend, etc.

### Terminal UI

UI files must contain `index.html` and `ui-config.yaml`. Index is an entry point for the UI. Config file contains a settings for the UI (see below).

WebKit engine is used to serve HTML. There are some limitations:

- The terminal have a touchscreen, so please keep in mind double-taps and mistaps. See details below;
- Utilize all the required resources locally, i.e. JS, CSS, icons, etc., except videos;
- Do not embed huge resources like video into the UI package;
- Do not use transparent video;
- Database is unavailable, use local storage instead;
- PDF rendering is not supported;
- WebGL is not available;
- Java is not available;

#### Touchscreen

The terminal have a touchscreen, so keep in mind touchscreen mistaps.

Scenario:

- a user taps a button;
- the page changes;
- a new button appears on the same place on the screen;
- the user still holds his finger on the touchscreen;
- suddenly another tap event is fired;
- the new button now is also clicked and changes the page again;

On page update keep buttons disabled for some time (about 200-300ms).

#### UI package delivery

Delivery of the UI is done by uploading packed (zip) UI files to Goldex dashboard.

Current zip size limit is 30MiB.

#### UI config

UI config is `ui-config.yaml` inside UI zip package that defines externally allowed domains and should provide emergency contacts (hardware critical failures) to show to a customer:

```yaml
# Multiline text to show to a customer (along with "Please contact support team:") in case of critical terminal failure
emergency_contacts:
  - 'Phone: <some phone number here>'
  - 'Whatever: <whatever>'
# List of allowed domains/ports to perform fetch, XMLHttpRequests, images loading, etc. (localhost[:80] is allowed by default)
host_whitelist:
  - example.com
  - example.com:8080
```

### Terminal API

API is served locally on the terminal. It exposes **methods** to control the terminal from UI and sends **events** to notify UI (for instance, optional hardware could send events).

API is a [JSONRPC 2](https://www.jsonrpc.org/specification) API over [Websocket](https://en.wikipedia.org/wiki/WebSocket) connection (`http://localhost:80/ws`).

[JSONRPC 2 batch](https://www.jsonrpc.org/specification#batch) requests are not supported. Moreover, hardware-related methods should be called sequently, error will be returned otherwise.

[Documentation](https://goldexrobot.github.io/core.integration/swagger/#/terminal-api-v1).

### Terminal API emulator

Goldex provides terminal API emulator to simplify UI development. You'll find binaries [here](https://github.com/goldexrobot/core.integration/releases).

The emulator serves terminal API on your local machine. The emulator does not emulate optional hardware, but can be connected to the Goldex sandbox environment, therefore can communicate with your testing backend. You'll need TLS certificate issued by Goldex to connect the emulator to the sandbox environment, so please contact Goldex team to get it.

Emulator accepts commands to simulate connectivity, hardware and functional errors which could be occurred during real terminal usage. Please pay attention developing your UI.

#### Terminal API flow

There is also a [sequence diagram](/docs/images/terminal_interaction_diagram.png)

![Terminal API flow](/docs/images/terminal_api.png)
