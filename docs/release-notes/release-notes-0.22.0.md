# Release Notes
- [Bug Fixes](#bug-fixes)
- [New Features](#new-features)
    - [Functional Enhancements](#functional-enhancements)
    - [RPC Additions](#rpc-additions)
    - [lncli Additions](#lncli-additions)
- [Improvements](#improvements)
    - [Functional Updates](#functional-updates)
    - [RPC Updates](#rpc-updates)
    - [lncli Updates](#lncli-updates)
    - [Breaking Changes](#breaking-changes)
    - [Performance Improvements](#performance-improvements)
    - [Deprecations](#deprecations)
- [Technical and Architectural Updates](#technical-and-architectural-updates)
    - [BOLT Spec Updates](#bolt-spec-updates)
    - [Testing](#testing)
    - [Database](#database)
    - [Code Health](#code-health)
    - [Tooling and Documentation](#tooling-and-documentation)

# Bug Fixes

# New Features

## Functional Enhancements

## RPC Additions

## lncli Additions

# Improvements

## Functional Updates

## RPC Updates

* [`AddHoldInvoice` now supports optional preimage/hash
  generation](https://github.com/lightningnetwork/lnd/pull/10685). The `hash`
  field is no longer required — when omitted, the server auto-generates a
  cryptographically random preimage and derives the payment hash. A new
  `preimage` field on the request allows callers to supply their own preimage
  and let the server derive the hash, eliminating the risk of hash/preimage
  mismatches. The response now includes `payment_preimage` (when the server
  knows the preimage) and `payment_hash` (always populated).

## lncli Updates

* The [`addholdinvoice` command now accepts `--hash` and `--preimage`
  flags](https://github.com/lightningnetwork/lnd/pull/10685). When neither is
  provided, the server generates both automatically. The legacy positional hash
  argument is still supported for backward compatibility.

## Code Health

## Breaking Changes

## Performance Improvements

## Deprecations

# Technical and Architectural Updates

## BOLT Spec Updates

## Testing

## Database

## Code Health

## Tooling and Documentation

# Contributors (Alphabetical Order)

* Suheb
