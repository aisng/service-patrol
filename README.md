# Service patrol

An application to check whether certain websites are online or offline. If a limit is reached or exceeded, an email is sent to the provided mailing list. An email is also sent when the connection is recovered.

## Configuration

A `config.yaml` file must be provided in the root dir with the following fields: 
```bash
down_limit: 2       # if the limit is reached or exceded, an email is sent
timeout_s: 5        # max timeout in seconds when pinging an url
frequency_h: 2      # amount of hours after the check will be performed again as configured in cron
services:
- https://www.google.com
mailing_list:
- example@example.org
```

## Email credentials

In order to send an email the following env variables must be provided:
```bash
SPMAILUSERNAME
SPMAILTOKEN
```

