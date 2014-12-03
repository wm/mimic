# mimic

Copy multiple databases from one server to another.

`mimic` allows you to easily retrieve and restore databse dumps from one server to another.

At [Iora](https://github.com/IoraHealth/) we use this to dump the databases
from our staging postgresql server to our local development environment with
dev_config.json. Or to our Docker container with docker_config.json

## Usage

Run `mimic` and select the apps you would like to restore:

    $ mimic
    Enter the configuration file: [./config.json]
    Please select which applications you'd like to replace

        0: all
        1: icispatients
        2: icisstaff
        3: snowflake
        4: cronos
        5: bouncah

    e.g. 1,2

    > 0

## Enhancements

I want to remove the dependency on os.exec calls to postgresql commands.
Ideally this could run on a machine without postgresql installed but just a
Docker client.
