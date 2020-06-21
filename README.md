[![CircleCI](https://circleci.com/gh/sudachen/smwlt.svg?style=svg)](https://circleci.com/gh/sudachen/smwlt)
[![Maintainability](https://api.codeclimate.com/v1/badges/45a94df6ce0f10650766/maintainability)](https://codeclimate.com/github/sudachen/smwlt/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/45a94df6ce0f10650766/test_coverage)](https://codeclimate.com/github/sudachen/smwlt/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/sudachen/smwlt)](https://goreportcard.com/report/github.com/sudachen/smwlt)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

```

Usage:
  smwlt [command]

Available Commands:
  coinbase    Set the account as coinbase account in the node
  export      Export account key pair as a hex string
  help        Help about any command
  import      Import account key pair from the hex string
  info        Display accounts info (address, balance, and nonce)
  net         Display the node status
  new         Create a new account (key pair)
  send        Transfer coins from one to another account
  signhex     Sign a hex message with the account private key
  signtext    Sign a text message with the account private key
  txs         List transactions (outgoing and incoming) for the account

Flags:
  -e, --endpoint string      host:port to connect mesh node (default "localhost:9090")
  -h, --help                 help for info
  -l, --legacy               use legacy unencrypted file format
  -p, --password string      wallet unlock password
  -x, --trace                backtrace on panic
  -v, --verbose              be verbose
  -d, --wallet-dir string    use wallet dir (default "/home/monster/.config/Spacemesh")
  -f, --wallet-file string   use wallet filename
  -n, --wallet-name string   select wallet by name
  -y, --yes                  auto confirm

Use "smwlt [command] --help" for more information about a command.
```
