# Known Issues

## Go Compile Validation

`env GOCACHE=/tmp/agentapi-plusplus-go-build-cache go test -run "^$" ./...`
is blocked by local disk exhaustion while the Go toolchain writes build
artifacts. The generated `/tmp` cache was removed after the failed run.

Full runtime tests may also be constrained by sandbox listener/network
permissions. This badge-only lane uses diff hygiene and badge proof as the
passing validation set, with the Go compile blocker recorded explicitly.

## Branch Integration

The refreshed badge evidence is prepared on `docs/agentapi-plusplus-sladge-current`
and is not merged into the canonical checkout in this pass.
