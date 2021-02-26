### Keylight CLI - A simple CLI to control your Elgato Keylights

A simple command line tool written in Golang for controlling Elgato [Key
Lights](https://www.elgato.com/en/gaming/key-light) and [Key Light
Airs](https://www.elgato.com/en/gaming/key-light-air).

## Install

```
$ go get github.com/deepu105/keylight
```

OR

Clone this repo first

```shell
$ git clone https://github.com/deepu105/keylight.git
$ cd keylight
```

Run Go Install

```shell
$ go install
```

or Build and move the binary to your PATH

```shell
$ go build -v .

$ mv keylight /usr/local/bin/
```

## Usage:

```shell
$ keylight command [command options]
```

## Commands:

### `switch`, `s` : Toggle the light switch

**Usage**:

```shell
$ keylight switch [command options]
# or
$ keylight s [command options]

# Examples

$ keylight s # Toggle all lights ON or OFF
$ keylight s --on # Switch on all lights. Same as `keylight s -o`
$ keylight s -o -b 15 # Switch on all lights, and set brightness to 15
$ keylight s -o -b 100 -t 4000 # Switch on all lights, and set brightness to 100 and temperature to 4000 kelvin
$ keylight s -o -b 25 -t 7000 --timeout 5 # Switch on all lights, and set brightness to 25 and temperature to 7000 kelvin. Wait 5 seconds before timing out
$ keylight s -o -l E855 -b 25 -t 7000 # Switch on a specific light, and set brightness to 25 and temperature to 7000 kelvin.
$ keylight s -o -l E855 # Switch on a specific light
$ keylight s -p warm # Switch on a specific preset
```

**Options**:

```
   --light value, -l value        ID, example E859, for the light to control. If not provided all lights will be updated (default: "all")
   --on, -o                       Switch light on. If not provided the light power state will be toggled based on last state (default: false)
   --brightness value, -b value   Set brightness of the lights (0 to 100) (default: -1)
   --temperature value, -t value  Set temperature of the lights in kelvin (3000 to 7000) (default: -1)
   --preset value, -p value       Switch on and set a preset temperature and brigtness. 
                                  Values: [warm warm-50 warm-100 cool cool-50 cool-100 normal normal-50 normal-100]
   --timeout value                Timeout for light discovery in seconds (default: 2)
   --help, -h                     show help (default: false)

```

Light id is the last part in the Name when you run `keylight list`

### `list`, `l` : Discover and list available lights

**Usage**:

```shell
$ keylight list [command options]
# or
$ keylight l [command options]
```

**Options**:

```
   --timeout value                Timeout in seconds (default: 2)
   --help, -h                     show help (default: false)
```

**Output**:

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
