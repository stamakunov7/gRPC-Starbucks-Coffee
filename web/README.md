# Coffee Shop Web UI

Next.js UI for the Coffee Shop project.

## Architecture (BFF)

- The browser calls **Next.js API routes** (`/api/*`) using `fetch`.
- Those route handlers call the Go backend over **gRPC** (default `localhost:9001`).
- The Go backend reads menu + orders from **Firestore** (emulator in local dev).

## Local development

In one terminal:

```bash
make emulators
```

Seed the menu once (or whenever you want to reset data):

```bash
make seed
```

Run the Go gRPC server:

```bash
make run-server
```

Run the UI:

```bash
make install-ui   # first time
make run-ui
```

Open `http://localhost:3000`.

## Environment

- **`GRPC_ADDR`**: gRPC server address for the Next.js BFF (default `localhost:9001`).
  - `make run-ui` sets it automatically.

