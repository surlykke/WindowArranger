# WindowArranger

WindowArranger is a simple tool to arrange windows when using Sway or i3. 

In the following we will describe WindowArranger as used on ```sway```. If you're on i3, you should read ```swaymsg``` as ```i3-msg```.

## Install

You need go and gtk libraries.

On Ubuntu this is probably sufficient:

```
sudo apt install golang libgtk-3-dev libglib2.0-dev libgdk-pixbuf2.0-dev
```

If you're on another distro, please consult the manual.

To build, do:

```
cd to/where/you/want/WindowArranger-dir

git clone https://github.com:surlykke/WindowArranger.git
cd WindowArranger
./install.sh
```

This installs two executables ```WindowArranger``` and ```dummywindow``` into ```$HOME/.local/bin```. 
(which has to be in your `$PATH`)

More on `dummywindow` below

## Run

```
WindowArranger configfile
```

where `configfile` is a file defining how you'd like your windows arranged. 

Assume you have 4 windows open, imaginatively named Window1, Window2, Window3 and Window4. Then `configfile` could look like this:

```
eDP-1: H[V['title=Window1' 'title=Window2'] T['title=Window3' 'title=Window4']]
```

This will place a horizontally split workspace on output 'eDP-1', with two containers. The first _split vertically_ containing Window1 and Window2, the second _tabbed_ with Window3 and Window4.

A string like "title=Window1" is used to select a window, and is, in fact, passed on to swaymsg, in constructs like 

```
swaymsg '[title=Window1] focus'
``` 

So those strings must be valid swaymsg criteria. 

Criteria should pick exactly one window. If a criteria matches several windows it's unpredictable how the layout ends up.

If no configfile is given WindowArranger will read the configuration from standard input, so you could also do:

```
WindowArranger <<EOF
    eDP-1: H[V['title=Window1' 'title=Window2'] T['title=Window3' 'title=Window4']]
EOF

```


## Cofiguration syntax

The syntax is somewhat aligned with how ```swaymsg -t get_workspaces``` reports layouts.

Informally, a configuration consists of a sequence of expressions of the form:

```
 output: container
```

`output` is the name of an output - eg. eDP-1 or DP-1 and should match one of your outputs.

A container is of the form:
```
 layout [ content ]
```
where `layout` is one of:
```
H   for split horizontal
V   for split vertical
T   for tabbed
S   for stacking
```

and `content` is a (space separated) list of criteria and/or containers. 


More formally, the syntax, in [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form), is:
```
    configuration     ::= output* ;
    output            ::= output-identifer, ':', container ;
    output-identifier ::= letter, (letter | digit | '-')* ;
    container         ::= ('V' | 'H' | 'T' | 'S'), '[', (criteria | container)+, ']' ;
```

- `letter` and `digit` are as defined by the unicode standard.
- `criteria` is a single quoted string, ie. a sequence of characters enclosed in single quotes (`'`). There is no escape mechanism, so a criteria cannot contain single quotes (but double quotes).
- Whitespace is ignored (except as separator). Anyting from a `#` to end of line (comments) is ignored unless inside a criteria.

Each output expression will create a workspace with the given layout and place it on the output. Workspaces will be numbered in the order they are encountered. 

So:

```
eDP-1: T['title=VPN' V['title=Work' 'title=Log']]
eDP-1: H['title=DbVisualizer' 'app_id=firefox']
DP-1:  H['instance=chromium' title=IntelliJ] V['instance: slack' 'title: "^Microsoft Teams"']]
```

would create 3 workspaces: 1 and 2 placed on eDP-1 and 3 placed on DP-1.

### Usage 

```
  WindowWrapper [option]... [configfile]
```

#### Options
```
    -dump string
    -wait uint
```

##### Dump

WindowArranger works by transforming the configuration into a bash script file containing mostly `swaymsg` commands, and then run it.

In stead of running the generated script, `WindowArranger` can write it to stdout or a file. Use the `dump` option to do that:

```
WindowArranger -dump arrangescript.sh configfile
```

or

```
WindowArranger -dump - configfile
```

The former variant will write to file arrangescript.sh, the latter to stdout.


##### Wait 

With the `wait` option you can instruct `WindowArranger` to wait until the windows you want to arrange are present. Say you have a criteria `title=Window1` in your configuration, then `WindowArranger` will wait until a window titled 'Window1' is present. So with:

```
WindowArranger -wait 20 configfile
```

`WindowArranger` will wait up to 20 seconds for all criteria in the configuration to find a match.

If the 20 seconds pass without all windows appearing, `WindowArranger` exits with a non-zero exit code.

### Shebang

WindowArranger functions as an interpreter of config files, so you could also write your configfile with a shebang:

```
#!/usr/bin/env WindowArranger
eDP-1: H[V['title=Window1' 'title=Window2'] T['title=Window3' 'title=Window4']]
```

and run with just:

```
./configfile
```

### Surplus windows

Open windows that are not mentioned in the configuration will be left in a workspace named `window_arranger_temp_workspace`.

### Dummy windows

A slight difficulty with swaymsg is that you can't really create containers. What you can do, is focus on a window, and then call eg. `splitv`. Therefore, in order to create the containers, specified by a configuration, WindowArranger resorts to an ugly hack: It creates a dummy window, focuses it, and calls `splitv` or one of the other layouts on it, to create a container, and then fills specified windows and subcontainers into that.

Once the layout is completed all dummywindows are closed. 

To that end, when you install WindowArranger, you also get a small stupid program `dummywindow` in $HOME/.local/bin. `dummywindow' can be called like:

```
dummywindow some-title
```

to create a window titled 'some-title'

WindowArranger uses `dummywindow` to create the dummy windows. 
