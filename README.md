## WindowArranger

WindowArranger is a simple tool to arrange windows when using the Sway window manger.

Support for I3 may come later...

### Install

You need go and gtk libraries.

On Ubuntu this would probably be sufficient:

```
sudo apt install golang libgtk-3-dev libglib2.0-dev libgdk-pixbuf2.0-dev
```

If you're on another distro, please consult the manual.

To build, do:

```
cd to/dir/where/you/want/WindowArranger-dir

git clone https://github.com:surlykke/WindowArranger.git
cd WindowArranger
./install.sh
```

This installs two executables ```WindowArranger``` and ```dummywindow``` into ```$HOME/.local/bin```. 

```$HOME/.local/bin``` has to be in your ```$PATH```

(more on ```dummywindow``` below)

### Run

```
WindowArranger configfile
```

where `configfile` is a file defining how you'd like your windows arranged. 

Assume you have 4 windows open, imaginatively named Window1, Window2, Window3 and Window4. Then `configfile` could look like this:

```
eDP-1: H[V['title=Window1' 'title=Window2'] T['title=Window3' 'title=Window4']]
```

This will place a horizontally split workspace on output 'eDP-1', with two containers. The first _split vertically_ containing Window1 and Window2, the secoond _tabbed_ with Window3 and Window4.

A string like "title=Window1" is used to select a window, and is, in fact, passed on to swaymsg, in constructs like 

```
swaymsg '[title=Window1] focus
``` 

So those strings must be valid swaymsg selectors.

### Dump

WindowArranger works by transforming the config file into a bash script file containing mostly ```swaymsg``` commands, and then run it.

In stead of running the generated script you may have `WindowArranger` write it to stdout or a file. Use the 'dump' option to do that:

```
WindowArranger -dump arrangescript.sh configfile
```

or

```
WindowArranger -dump - configfile
```

The former variant will write to file arrangescript.sh, the latter to stdout.

### Config file syntax

The syntax is somewhat aligned with how ```swaymsg -t get_workspacs``` reports layouts.

A config file consists of a sequence of expressions of the form:

```
<output>: <container>
```

&lt;output&gt; is the name of an output - eg. eDP-1 or DP-1

A container is of the form:
```
<layout>[ <content> ]
```
where &lt;layout&gt; is one of:
```
H   for split horizontal
V   for split vertical
T   for tabbed
S   for stacking
```

and &lt;content&gt; is a (space separated) list of selectorstrings and/or containers. 

Each output expression will create a workspace with the given layout and place it on the output. Workspaces will be numbered in the order they are encountered. 

So:

```
eDP-1: T['title=VPN' V['title=Work' 'title=Log']]
eDP-1: H['title=DbVisualizer' 'app_id=firefox']
DP-1:  H['instance=chromium' title=IntelliJ] V['instance: slack' 'title: "^Microsoft Teams"']]
```

would create 3 workspaces: 1 and 2 placed on eDP-1 and 3 placed on DP-1.

You can also write comments. Anything from a `#` tho end-of-line is ignored:

```
# This is ignored
eDP-1: H[V['title=Window1' 'title=Window2'] T['title=Window3' 'title=Window4']]
```


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

### dummy windows

A slight difficulty with swaymsg is that you can't really create containers. What you can do, is focus on a window, and then call eg. `splitv`. Therefore, in order to create the containers, specified by a configuration, WindowArranger resorts to an ugly hack: It creates a dummy window, focuses it, and calls `splitv` or one of the other layouts on it, to create a container, and then fills specified windows and subcontainers into that.

Once the layout is completed all dummywindows are closed. 

To that end, when you install WindowArranger, you also get a small stupid program `dummywindow` in $HOME/.local/bin. WindowArranger uses that to create the dummy windows. 