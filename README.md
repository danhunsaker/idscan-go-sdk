# IDScan Go SDK

This repo encodes the IDScan web API(s) in a Go package.

## Screening API

The only supported API at the moment is the Screening API. It consists of several services which you can query individually with the minimum required
data, or you can construct your own `ScreeningAPIRequest` and pass in whatever supported data you prefer. Consult
[the documentation for the API](https://docs.idscan.net/screening/screening-services.html) to ensure you know how to pass the required data correctly.
