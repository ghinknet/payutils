# payutils

A pluggable, framework-agnostic unified payment toolkit for Go.

`payutils` is the **core**: it owns trade orchestration (create dispatch, callback
routing, status / close / refund) and defines the interfaces that drivers
implement. It depends on **no** payment SDK and **no** HTTP framework. Everything
concrete is provided by plug-in drivers:

- **Pay drivers** talk to a payment provider (Alipay, WeChat Pay, ...).
- **HTTP drivers** adapt a web framework's router (GoFiber, Echo, Gin, Iris, ...).

```
                    ┌─────────────────────────┐
                    │     payutils (core)     │
                    └─────────────────────────┘
                        ▲                  ▲
                 depends│           depends│
             ┌──────────┴────────┐  ┌──────┴────────────┐
             │  HTTP drivers     │  │  Pay drivers      │
             │  (fiber/echo/...) │  │  (alipay/wechat)  │
             └───────────────────┘  └───────────────────┘
        ✗  the two driver kinds never depend on each other  ✗
```

The two driver families are glued together only through the Go standard library
`net/http` types, so any *HTTP framework × payment provider* combination works,
and adding either side never touches the core.

## Module

```
go.gh.ink/payutils/v3
```

```bash
go get go.gh.ink/payutils/v3
```

## Design at a glance

| Concern            | Who owns it                                                        |
|--------------------|--------------------------------------------------------------------|
| Order creation     | **The user.** payutils does not fetch order details; you build the trade detail and pass it to `Create`. |
| `Create`           | A method on `*client.Client` that dispatches a driver-specific param struct to the registered pay drivers. |
| `Callback`         | Auto-registered route (when an HTTP driver is provided) **and** a manual `(*Client).Callback` method that accepts a standard `*http.Request`. |
| Status updates     | Pushed to your `Contract.StatusUpdater` (trade ID prefix/suffix stripped for you). |

## Installation & wiring

A working setup needs the core plus at least one pay driver. An HTTP driver is
optional — import one to get callback routes auto-registered, or omit it and
call `Callback` yourself.

```go
import (
    "go.gh.ink/payutils/v3/client"
    "go.gh.ink/payutils/v3/model"

    // Pay drivers self-register via init().
    payAlipay "go.gh.ink/payutils/pay/alipay/v3"
    _ "go.gh.ink/payutils/pay/wechat/v3"

    // Optional: an HTTP driver self-registers via init().
    httpFiber "go.gh.ink/payutils/http/fiber/v3"
)
```

## Quick start

```go
app := fiber.New() // any framework supported by an imported HTTP driver

c, err := client.NewClient(model.Config{
    Endpoint: "https://api.example.com", // public base URL upstreams call back to
    Contract: myContract{},              // your StatusUpdater implementation

    // Register callback routes on these framework instances. Keyed by driver name.
    // Omit entirely to handle callbacks manually.
    Instances: model.I{ // map[string]any
        httpFiber.Name: app,
    },

    // Per-provider credentials. Keyed by driver name; inner keys are
    // driver-defined (see each driver's README).
    Credentials: model.C{ // map[string]map[string]string
        payAlipay.Name: {
            payAlipay.AppID:             "...",
            payAlipay.AppCertPrivateKey: "...",
            payAlipay.AppCert:           "...",
            payAlipay.RootCert:          "...",
            payAlipay.PublicCert:        "...",
            payAlipay.IsProd:            "true",
        },
    },
})
if err != nil {
    panic(err)
}
```

### Creating a trade

You prepare the order yourself and wrap it in the driver's param struct. The
first driver that recognises the concrete type handles it.

```go
result, err := c.Create(ctx, payAlipay.CreateParam{
    TradeID:  "order-123",
    Platform: payAlipay.PlatformPC,
    Detail: model.TradeDetail{
        Subject:  "A nice product",
        Price:    1990,                 // in cents
        Currency: "CNY",
        Expiry:   time.Now().Add(time.Hour),
    },
})
// result is driver-specific (ready to marshal), e.g. {"payUrl": "https://..."}.
// errors.Is(err, errors.ErrNoDriverClaimed) means no driver matched the param.
```

`result` is `any` so each provider can return its natural payload (a pay URL, a
QR `code_url`, a JSAPI sign object, ...). Marshal it straight to your client, or
type-assert when you need the concrete value.

### Handling callbacks

**Option A — auto-registered route.** Import an HTTP driver and pass its instance
in `Instances`. payutils registers `POST /{provider}/callback` for every
configured provider (e.g. `/alipay/callback`).

**Option B — manual.** Route the request yourself and forward the standard
request:

```go
func (h *handler) onAlipayNotify(w http.ResponseWriter, r *http.Request) {
    if err := c.Callback("alipay", w, r); err != nil {
        // e.g. errors.Is(err, errors.ErrUpstreamNotFound)
    }
}
```

Either way the driver verifies the signature, decodes the notification, and
pushes the resolved status to your `Contract.StatusUpdater`.

## The `Contract`

Implement `model.Contract` to receive trade-state updates:

```go
type Contract interface {
    StatusUpdater(
        ctx context.Context, r *http.Request,
        upstream string, tradeID string,
        status TradeState, t time.Time,
    ) error
}
```

`tradeID` is delivered with the configured `TradeIDPrefix` / `TradeIDSuffix`
already stripped. `Contract` is optional: if `nil`, status updates are a no-op
(useful when you only need `Create` / manual handling).

## Configuration reference (`model.Config`)

| Field                 | Meaning |
|-----------------------|---------|
| `Endpoint`            | **Required.** Public base URL; used to build callback (`notify_url`) URLs. |
| `Credentials`         | `map[provider]map[key]value` of provider credentials. |
| `Instances`           | `map[framework]router` for auto callback registration. Optional. |
| `Contract`            | Your `StatusUpdater`. Optional. |
| `Debug`               | Enables driver SDK debug logging. |
| `TradeIDPrefix/Suffix`| Wraps your trade ID when talking to providers; stripped on the way back. |
| `NoNewPaymentWindows` | Minimum time that must remain before expiry to open a payment. Default 30s. |
| `SafetyMargin`        | Pulled off the expiry sent to providers. Default 10s. |
| `ErrorHandler`        | Optional hook invoked when a callback fails. |
| `Marshal` / `Unmarshal` | Custom JSON codecs. Default `encoding/json`. |

## Trade states (`model.TradeState`)

`PENDING` · `SUCCESS` · `CLOSED` · `FINISHED` · `UNKNOWN` — providers map their
own states onto these.

## Errors

Sentinel errors live in `go.gh.ink/payutils/v3/errors` and are matchable with
`errors.Is` even after a driver attaches upstream metadata via the chainable
`With*` helpers (`WithUpstreamName`, `WithUpstreamCode`, ...). Notable ones:
`ErrNoDriverClaimed`, `ErrUpstreamNotFound`, `ErrTradeNotExist`,
`ErrNoEnoughTimeToPay`, `ErrUnsupportedCurrency`, `ErrUnsupportedPlatform`,
`ErrUnsupportedInstance`.

## Writing your own driver

- **Pay driver:** implement `model.PayDriver` (→ `model.PayClient`), define your
  own `CreateParam`, and `driver.RegisterPay(name, Driver{})` in `init()`.
- **HTTP driver:** implement `model.HttpDriver` (→ `model.HttpInstance`, which
  exposes `Get/Post/Put/Patch/Delete/Head/Options/Any`) and
  `driver.RegisterHttp(name, Driver{})` in `init()`.

No core changes are ever required.

## License

See [LICENSE](LICENSE).
