# unce
*Why poll when you can unce?*

## Initial Setup
`git clone https://github.com/kofalt/unce.git && ./unce/goad && ./unce/unce`

You'll receive `~/.config/unce/unce.toml`
There's a spot for various consumers; ATM, github is supported.

If `notify-send "lol"` doesn't work on your system, `unce` won't either.

## Setup for GitHub
1. Browse to https://github.com/settings/tokens/new.
2. Select only notifications (limits access).
3. Create, paste token in ur TOMLs.
4. Run command again. `./unce/unce`


## Logs
Bolt DB location for logs: `~/.local/share/unce/`  
Show all the unique JSON you've downloaded in your travels with `./unce/unce whelp log github`
