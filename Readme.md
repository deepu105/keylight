### Keylight CLI - A simple CLI to control your Elgato Keylights

A simple command line tool written in Golang for controlling Elgato [Key
Lights](https://www.elgato.com/en/gaming/key-light) and [Key Light
Airs](https://www.elgato.com/en/gaming/key-light-air).

## USAGE:

```shell
$ keylight command [command options]
```

## COMMANDS:

### `switch`, `s` : Switch on/off lights

**USAGE**:

```shell
$ keylight switch [command options]
```

**OPTIONS**:

```
   --light value, -l value        ID, example E859, for the light to control. If not provided all lights will be updated (default: "all")
   --on                           Toggle light on (default: false)
   --off                          Toggle light off (default: false)
   --brightness value, -b value   Set brightness of the lights (0 to 100) (default: 10)
   --temperature value, -t value  Set temperature of the lights in kelvin (3000 to 7000) (default: 3000)
   --timeout value                Timeout in seconds (default: 2)
   --help, -h                     show help (default: false)
```

Light id is the last part in the Name when you run `keylight list`

### `list`, `l` : Discover and list available lights

**USAGE**:

```shell
$ keylight list [command options]
```

**OPTIONS**:

```
   --timeout value                Timeout in seconds (default: 2)
   --help, -h                     show help (default: false)
```

**OUTPUT**:

```
+-------------------------------+-------------+------------+--------------+---------------------------------------+
| Name                          | Power State | Brightness | Temperature  | Address                               |
+-------------------------------+-------------+------------+--------------+---------------------------------------+
| Elgato\ Key\ Light\ Air\ E859 | on          |         10 | 331 (3021 K) | elgato-key-light-air-e859.local.:9123 |
+-------------------------------+-------------+------------+--------------+---------------------------------------+
```

### `help`, `h` : Shows a list of commands or help for one command

Run `keylight command --help` for info about flags specific to a command

## Licence

MIT
