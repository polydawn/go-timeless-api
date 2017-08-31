Timeless Stack APIs for Go(lang)
================================

Go(lang) APIs for the Timeless Stack ([:book: docs here](https://github.com/polydawn/timeless/).

These APIs include all message type definitions, message and error enums, and serialization glue for the APIs for projects like
[rio](https://github.com/polydawn/rio),
[repeatr](https://github.com/polydawn/repeatr),
and [hitch](https://github.com/polydawn/hitch).

Function definitions for some RPC APIs are also exported.
Generally you will be able to find implementations of these func APIs in each project,
and also a package which exports the same funcs again, but works via an exec layer and RPC protocol.
