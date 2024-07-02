# Service patrol

An application to check whether certain websites are online or offline. If a limit is reached or exceeded, an email is sent to the provided mailing list. An email is also sent when the connection is recovered.

## Configuration

A `config.yaml` file must be provided in the root dir with the following fields: 
```bash
down_limit: 2                   # if reached or exceded, an email is sent
packet_loss_limit_percent: 10   # if reached or exceeded, ip address is considered unavailable
timeout_s: 5                    # max timeout in seconds when pinging an address
frequency_h: 2                  # amount of hours after which the app will be run again as configured in cron
provider: PROV                  # name of the service(s) provider
services:
- https://www.google.com
mailing_list:
- example@example.org
```

## Email credentials

In order to send an email the following env variables must be provided:
```bash
SP_MAIL_USERNAME
SP_MAIL_PASSWORD
```

### Permissions

Pinging raw ip addresses requires sudo privileges, thus `--preserve-env` (or `-E`) flag must be provided to read the env variables.

```bash
sudo -E ./service-patrol
```

## Notes

- Currently configured with `smtp-mail.outlook.com:587` email server with LOGIN Auth method.