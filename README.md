[![Build and Test](https://github.com/Toshik1978/adyen-cli/workflows/Build%20and%20Test/badge.svg)](https://github.com/Toshik1978/adyen-cli/actions)
# !!! DISCLAIMER !!!

This is not official Adyen application!
You are taking responsibility for any damage caused to your Adyen account by using this application!

## Adyen CLI

This application was designed to automate some routine operations over Adyen account.

## Usage
### Link split configurations to stores

1. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production or test).
3. Create the CSV file with the information about stores and splits.
   1. CSV should contain 3 mandatory column - 'Account Holder Code', 'Store ID', 'Split ID'.
   2. Account Holder Code should be Balance Account ID in case of Balance Platform linking.
4. Run linking: `adyen-cli link --csv <Path to file> --prod`.
5. Run `adyen-cli -h` if you have questions.

## How to build it?
### Prerequisites

Go 1.20+ (should be built with lower version, but I didn't try it)
GNU Make (tested on 3.81)
Internet connection to download dependencies

### Build

Run `make app.build`.
You will find the binary in `./bin` folder.

## How to improve it?

Feel free to fork it, change it, send me PR.
