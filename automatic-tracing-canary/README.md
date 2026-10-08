# Automatic tracing canary

A Go 1.25 Gin application with no OpenTelemetry SDK, middleware or exporter.
Hostless compile-time instrumentation must generate the HTTP server and client spans.

Use Railpack with root directory `automatic-tracing-canary`, the default build command,
and automatic tracing. `/health` checks readiness; `/outbound` calls the local health
endpoint; `/chain` calls the configured `CHAIN_TARGET_URL`. PORT defaults to 8000.

For rollback restore the original root directory, source branch and export-only mode,
then rebuild the existing SDK-enabled example. No main-branch change is needed to
run this canary from its feature branch.
