# payacm

[![Build Status](https://travis-ci.org/acmumn/payacm.svg?branch=develop)](https://travis-ci.org/acmumn/payacm)

A simple API service to accept payment. This site uses Stripe to collect
arbitrary amounts of money from people, emailing them afterwards.

## Configuration

Runs on the port specified in the `PORT` environment variable, or on port 3000
by default.

Uses the Stripe public key in the `STRIPE_PUBLIC_KEY` environment variable and
the secret key in the `STRIPE_SECRET_KEY` environment variable.

Sends mail using the SMTP server stored in the `SMTP_HOST` and `SMTP_PORT`
environment variables, from the email specified in the `SMTP_FROM` environment
variable, authenticating with the username and password in the `SMTP_USER` and
`SMTP_PASS` environment variables, respectively.

If the application is running in a production environment, it is recommended to
also set the `GIN_MODE` environment variable to `release`.

## Developing

This repository uses Git Flow -- every feature, bugfix, or other change should
be developed in a separate branch, which is then pull-requested against the
`develop` branch. When sufficiently many changes have been collected in
`develop` to merit a new release, `develop` is merged into `master`. A webhook
automatically deploys any changes to `master`, so be sure everything's been
tested before you merge!
