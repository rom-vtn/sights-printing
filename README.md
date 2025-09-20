# Sights printing

A sample program to fetch current sights from my API and print today's remaining ones on a receipt.

## Requirements

- a train API (i have one)
- a font file (i use `lato-regular`)
- an ESC/POS receipt printer reachable over TCP

## Usage

Build and `./sights-printing <config.json>`

## Syntax

Sample `config.json` syntax:

```json
{
    "fetch_url": "https://trainmap.cozytren.ch/api/sights/<lat>/<lon>",
    "font_file": "./path/to/font.ttf",
    "printer_address": "192.168.1.87:9100"
}
```

## Legal

The data returned from the APIs comes from transit feeds often under Creative Commons or ODbL licenses. Please make sure the source licenses are alright with this usage.