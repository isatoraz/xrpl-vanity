# XRPL Vanity Wallet Generator

A high-performance Go tool for generating wallet addresses that start with or end with specific strings — ideal for creating **vanity wallet addresses**.

## ✨ Features

- Match wallet address **prefix** and/or **suffix**
- Supports **case-sensitive** or **case-insensitive** search
- Configurable number of **parallel goroutines** for fast generation
- Stops automatically after finding the desired number of matches
- Displays total attempts and time elapsed

## Usage

```bash
go run cmd/main.go --starts-with rLOL --ends-with xyz --count 2 --threads 8 --is-case-sensitive
```

## Flags

| Flag                  | Description                                                   | Default             |
|-----------------------|---------------------------------------------------------------|---------------------|
| `--starts-with`       | Prefix to match at the beginning of the wallet address        | `""`                |
| `--ends-with`         | Suffix to match at the end of the wallet address              | `""`                |
| `--count`             | Number of matching wallets to generate                        | `1`                 |
| `--threads`           | Number of goroutines to use for parallel generation           | CPU core count      |
| `--is-case-sensitive` | Whether to match case exactly (`true` or `false`)             | `false`             |

> ⚠️ **Note:** `--starts-with` must begin with the letter **`"r"`** to ensure a valid wallet address format.

## Example
Generate 3 addresses that start with rLOL (case-insensitive) using 4 threads:

```bash
go run cmd/main.go --starts-with rLOL --count 3 --threads 4
```

### Output
```bash
🚀 Searching with 4 goroutines for 3 wallet(s) matching:
  🔍 Prefix: "rLOL"
  🔠 Case-sensitive: false

🎯 Wallet 1 found!
  📍 Address: rLOL8uVh...
  🔑 Seed: s██████████████████

🎯 Wallet 2 found!
  📍 Address: rLOL1xeR...
  🔑 Seed: s██████████████████

🎯 Wallet 3 found!
  📍 Address: rLOLx7Qy...
  🔑 Seed: s██████████████████

🔁 Total attempts: 53972  
⏰  Elapsed time: 5s
```

## Build

```bash
go build -o vanitygen cmd/main.go
./vanitygen --starts-with rLO --count 1
```

## Notes
- Either --starts-with or --ends-with must be provided.
- If --count is less than 1, the program will exit with an error.
- The number of threads is capped at 1000 to avoid excessive CPU usage.
- Addresses generated must start with "r" to be valid (hence the --starts-with requirement).
