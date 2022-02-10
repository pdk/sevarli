# sevarli

## Set Variables from Lines

The purpose of this helper is to set shell variables from a line selected by a
pattern from a table-like text file.

For example, given this input:

    Customer ID     Customer Name                            Contact Name             Address                        City         Postal Code    Country
    1             Alfreds Futterkiste                       Maria Anders          Obere Str. 57                    Berlin         12209          Germany
    2             Ana Trujillo Emparedados y helados        Ana Trujillo          Avda. de la Constitución 2222    México D.F.    05021          Mexico
    3             Antonio Moreno Taquería                   Antonio Moreno        Mataderos 2312                   México D.F.    05023          Mexico
    4             Around the Horn                           Thomas Hardy          120 Hanover Sq.                  London         WA1 1DP        UK
    5             Berglunds snabbköp                        Christina Berglund    Berguvsvägen 8                   Luleå          S-958 22       Sweden

and the command

    sevarli -data example.data -pattern Antonio

the output is

    CUSTOMER_ID="3"
    CUSTOMER_NAME="Antonio Moreno Taquería"
    CONTACT_NAME="Antonio Moreno"
    ADDRESS="Mataderos 2312"
    CITY="México D.F."
    POSTAL_CODE="05023"
    COUNTRY="Mexico"

With the shell command `eval` this output can be used to set variables to be
used within a shell script.

    $ eval `sevarli -data example.data -pattern Antonio`
    $ echo $CITY
    México D.F.

If the pattern does not match exactly 1 line of the data, then the matching
lines are written to stderr, and the program exits with error code 1. Scripts
using this program should use this pattern: 

    # exit on any error
    set -e

    setvars=`sevarli -data example.data -pattern $PATTERN`
    eval $setvars

(`eval` will not return an error if/when the backtick expression returns an error.)

## Data Format

Columns must be separated by at least 2 spaces, across all lines of input.

Correct:

    First Name  Last Name
    Mickey      Mouse

Incorrect:

    First Name Last Name
    Mickey     Mouse

Comments can be included with either `#` or `//`.

## Options

    usage of sevarli:
    -pattern string   pattern to search for (*required)
    -caps             convert names to caps (default true)
    -data string      path to data file (otherwise read stdin)
    -export           export the vars
    -hide value       hide column(s) when listing
    -prefix string    prefix variable name with given value
    -suffix string    suffix variable name with given value

By default the names are converted to uppercase. To preserve case, use `-caps=false`.

Example with more options:

    $ sevarli -data example.data -pattern Antonio -prefix DATA_ -suffix _VAR -export
    export DATA_CUSTOMER_ID_VAR="3"
    export DATA_CUSTOMER_NAME_VAR="Antonio Moreno Taquería"
    export DATA_CONTACT_NAME_VAR="Antonio Moreno"
    export DATA_ADDRESS_VAR="Mataderos 2312"
    export DATA_CITY_VAR="México D.F."
    export DATA_POSTAL_CODE_VAR="05023"
    export DATA_COUNTRY_VAR="Mexico"

## Install

Assuming you're on a mac, the best path to install is probably

    go install github.com/pdk/sevarli@latest

Binaries (for mac and other platforms) can be downloade from
https://github.com/pdk/sevarli/releases, but for macs you'll need to go into
system prefs and allow the unknown binary. (Click the question mark in the
dialog that says you cannot run it.)