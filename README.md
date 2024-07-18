[![Build and Test](https://github.com/Toshik1978/adyen-cli/workflows/Build%20and%20Test/badge.svg)](https://github.com/Toshik1978/adyen-cli/actions)
# !!! DISCLAIMER !!!

This is not official Adyen application!
You are taking responsibility for any damage caused to your Adyen account by using this application!

## Adyen CLI

This application was designed to automate some routine operations over Adyen account.

## Usage
### Link split configurations to stores

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about stores and splits.
   1. CSV should contain 4 columns - 'Account Holder Code', 'Store ID', 'Split ID' and 'Merchant ID'.
   2. 'Account Holder Code' should actually contain 'Balance Account ID' in case of Balance Platform linking.
   3. 'Merchant ID' is not required for VIAS, only for Balance.
4. Run linking: `adyen-cli link --csv <Path to file> --prod`.
5. Run `adyen-cli -h` if you have questions.

### Close unused stores

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about stores.
   1. CSV should contain 2 columns - 'Account Holder Code', 'Store ID'.
4. Run the process: `adyen-cli close --csv <Path to file> --prod`.
5. Run `adyen-cli -h` if you have questions.

### Add payment methods to the store

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about stores.
   1. CSV should contain 3 columns - 'Store ID', 'Payment Methods', 'Currency'.
   2. Payment methods is the list of payment methods you want to add to the store. Methods can be separated with |.
   3. The possible methods are: "visa|discover|diners|mc|jcb|cup|amex|interac_card".
   4. Currency can be EUR, CAD or GBP.
4. Run the process: `adyen-cli methods --csv <Path to file> --prod` if you want to add payment methods.
5. Run `adyen-cli -h` if you have questions.

### Fix sweep configuration

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about balance account holders.
   1. CSV should contain one of 2 columns - 'Account Holder ID' or 'Balance ID'.
4. Run the process: `adyen-cli sweep --csv <Path to file> --prod` if you want to fix broken sweep configurations.
5. Run `adyen-cli -h` if you have questions.

### Change store sales close time

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about stores.
   1. CSV should contain 4 columns - 'Account Holder ID' or 'Balance ID', 'Close Time', 'TimeZone' and 'Delays'.
   2. 'Close Time' must be in format "HH:MM" and "HH" must be between 0 and 7.
   3. Delays defines you settlement delay in days.
4. Run the process: `adyen-cli time --csv <Path to file> --prod` if you want to change sales close time.
5. Run `adyen-cli -h` if you have questions.

### Re-assign terminals

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about terminals.
   1. CSV can contain 3 columns - 'Serial', 'Terminal ID', 'Merchant ID', 'Store ID'. You can use either Serial or Terminal ID. You can use either Merchant or Store ID. If you use Merchant, then the terminal will be assigned to inventory.
4. Run assignment: `adyen-cli reassign --csv <Path to file> --prod`.
5. Run `adyen-cli -h` if you have questions.

### Enable/disable cellular on the terminal

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about terminals.
   1. CSV can contain 2 columns - 'Serial', 'Terminal ID'. You can use either Serial or Terminal ID.
4. Run the process: `adyen-cli cellular --csv <Path to file> --prod` if you want to enable cellular and add `--disable` flag if you want to disable it.
5. Run `adyen-cli -h` if you have questions.

### Disable offline payments on the terminal

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about terminals.
   1. CSV can contain 2 columns - 'Serial', 'Terminal ID'. You can use either Serial or Terminal ID.
4. Run the process: `adyen-cli offline --csv <Path to file> --prod` if you want to disable offline payments.
5. Run `adyen-cli -h` if you have questions.

### Install Android application on the supported terminal

1. Download the version relevant to your PC from the [Releases](https://github.com/Toshik1978/adyen-cli/releases) page.
   1. You can check the signature of the downloaded file using the relevant minisign file and public key: `RWTDIoCdDlDV5mgrQt3IK1D3ZOVZMtMpxrOO+yZgFvBP1Sv/D1BXhEkE`.
2. Create the copy of `.env.dist` file locally.
2. Fill it with the actual keys from Adyen (production / test / etc).
3. Create the CSV file with the information about application.
   1. CSV can contain several columns - 'Company ID', 'Store ID', 'Terminal ID', 'Filter', 'Package Name', 'Version Name' and 'Date'.
   2. You should use 'Store ID' or 'Terminal ID'. If 'Store ID' defined, the tool will try to find all terminals under the store with the given 'Filter'.
   3. 'Company ID', 'Package Name' and 'Version Name' will be used to find an application in the list of all available applications.
   4. 'Date' can be empty, the tool will use NOW() + 2 minutes to schedule an installation.
   5. The 'Date' format is defined in the Adyen documentation: https://docs.adyen.com/api-explorer/Management/latest/post/terminals/scheduleActions#request-scheduledAt
4. Run installation: `adyen-cli install --csv <Path to file> --prod`.
5. Run `adyen-cli -h` if you have questions.

## How to build it?
### Prerequisites

- Go 1.20+ (should be built with lower version, but I didn't try it)
- GNU Make (tested on 3.81)
- Internet connection to download dependencies

### Build

Run `make app.build`.

You will find the binary for Intel macOS in `./bin` folder.

You can build the version relevant for your PC using the following build configurations:

- `app.build.darwin.amd64` - for Intel macOS.
- `app.build.darwin.arm64` - for Apple Mx macOS.
- `app.build.linux.amd64` - for Intel Linux x64.
- `app.build.windows.amd64` - for Intel Windows x64.

## How to improve it?

Feel free to fork it, change it, send me PR.
