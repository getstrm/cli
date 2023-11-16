# Developing the PACE Command Line Interface

## Development

In order to run tests locally while developing, copy `.env` file to `test` directory. The contents are loaded
by `godotenv` to ensure all test properties are set. Contents of the file can be found in 1Password.

## Protos Version

To update the PACE protos version to the latest release, run:

`make update-pace-protos-version`
