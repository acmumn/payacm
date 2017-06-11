# payacm

[![Build Status](https://travis-ci.org/acmumn/payacm.svg?branch=develop)](https://travis-ci.org/acmumn/payacm)

A simple API service to accept payment. This site uses Stripe to collect
arbitrary amounts of money from people, emailing them afterwards.

## Configuration

Adjust the API key in `static/main.js` and the secret key in the
`STRIPE_SECRET_KEY` environment variable.

Runs on the port specified in the `PORT` environment variable, or on port 3000
by default.

## Developing

This repository uses Git Flow -- every feature, bugfix, or other change should
be developed in a separate branch, which is then pull-requested against the
`develop` branch. When sufficiently many changes have been collected in
`develop` to merit a new release, `develop` is merged into `master`. A webhook
automatically deploys any changes to `master`, so be sure everything's been
tested before you merge!
