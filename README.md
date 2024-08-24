[![goreleaser](https://github.com/lvstb/dyn-dns/actions/workflows/gobuild-action.yaml/badge.svg)](https://github.com/lvstb/dyn-dns/actions/workflows/gobuild-action.yaml)

# Dynamic DNS Updater for Route53

This project is a Go application that updates DNS records in Amazon Route53 based on the current WAN IP address. It's useful for maintaining a dynamic DNS setup when your public IP address changes frequently.

## Features

- Retrieves the current WAN IP address
- Updates a specified DNS record in Amazon Route53
- Supports command-line arguments for domain name and hosted zone ID
- Uses AWS SDK for Go v2

## Prerequisites

- Go 1.21 or later
- AWS account with Route53 access
- AWS credentials configured (either through environment variables, AWS credentials file, or IAM role)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/dyndns.git
   cd dyndns
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Usage

Run the program with the following command:
```
go run main.go -domain <domain> -zone <zone>
```


Replace `your.domain.com` with the domain you want to update and `YOURHOSTEDZONEID` with your Route53 hosted zone ID.

### Command-line Arguments

- `-domain`: The domain name to update (required)
- `-zone`: The Route53 Hosted Zone ID (required)

## Configuration

The application uses the default AWS credential chain. Make sure your AWS credentials are properly configured in one of the following ways:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. AWS credentials file (`~/.aws/credentials`)

