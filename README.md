# WindowArranger

WindowArranger is a simple tool to arrange windows when using Sway 

## What can it do ?

Here's an example from my own setup. When working, I have (at least) google-chrome, intellij, slack, microsoft teams and a couple of terminals (titled 'Build' and 'Log') open.
I use a laptop with a 49'' ultrawide external monitor, and I use this yaml file to arrange my windows: 

```
layout:
- name: eDP-1 
  workspaces: 
  - T['title=Build' 'title=Log'] 
  posx: 0       #
  posy: 360
- name: DP-2 
  workspaces: 
  - H[T['app_id=google-chrome'] T['instance="jetbrains-idea"'] V['instance=slack' 'title=".*Microsoft Teams.*"']]
  posx: 1920 
  posy: 0
postcommands:
- 'swaymsg [instance=slack] resize set width 20ppt'
- 'swaymsg [app_id=google-chrome] resize set width 35ppt'
```

The file is called `layout.yaml` and running:

```
WindowArranger layout.yaml
```

gives me:

```
                               DP-2
                             -------------------------------------------------------------------------------------
                             | Chrome                      . Intellij                          . Slack           |  
                             |                             .                                   .                 |
                             |                             .                                   .                 |
eDP-1                        |                             .                                   .                 |
--------------------------   |                             .                                   .                 |
| Build      | Log       |   |                             .                                   ..................|
|                        |   |                             .                                   . Teams           |
|                        |   |                             .                                   .                 |
|                        |   |                             .                                   .                 |
|                        |   |                             .                                   .                 |
|                        |   |                             .                                   .                 |
--------------------------   -------------------------------------------------------------------------------------

```

Chrome and Intellij are put in a tabbed node for if I add more windows there later.


## Install

You need go installed.

On Debian or a Debian-derived distro this is probably sufficient:

```
sudo apt-get install golang 
```

If you're on another distro, please consult the manual.

To build, do:

```
cd to/where/you/want/WindowArranger-dir

git clone https://github.com:surlykke/WindowArranger.git
cd WindowArranger
./install.sh
```

This installs the executable ```WindowArranger``` into ```$HOME/.local/bin```. 
(which has to be in your `$PATH`)

## Run

```
WindowArranger <yaml file> 
```

where `yaml file` is a file defining how you'd like your windows arranged. 


## Yaml file format

The config file must be a valid yaml file defining a map with 2 entries: `monitors` and `postcommands`. 

### monitors

`monitors` must contain a list of maps, each defining a monitor setup. A monitor setup has the following keys:

* `name`: Name of the monitor. eg. `eDP-1` or `DP-2` as reported by `swaymsg -t get_outputs`
* `make`: Make of the moniter as reported by swaymsg
* `model`: Model of the moniter as reported by swaymsg
* `serial`: Serial of the moniter as reported by swaymsg

  At least one of `name`, `make`, `model` and `serial` must be given. They should uniquely identify the monitor.
* `workspaces`: a list of the workspaces you want to have on the monitor. Each workspace is given by a _node definition_. 
  The syntax of a node definition is 
  ```
   layout [ <content> ]
  ```
  where `layout` is one of:
  ```
  H   for split horizontal
  V   for split vertical
  T   for tabbed
  S   for stacking
  ```

  and `<content>` is a (space separated) list of criteria and/or nodes. 


  More formally, the syntax, in [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form), is:
  ```
      node ::= ('V' | 'H' | 'T' | 'S'), '[', (criteria | node)+, ']' ;
  ```

  - `criteria` is a single quoted string, ie. a sequence of characters enclosed in single quotes (`'`). There is no escape mechanism, so a criteria cannot contain single quotes (but double quotes).
  - Whitespace (outside of strings) is ignored 

* posx, posy: Optional. Sets the position of the monitor. Both or none of posx, posy must be given
* scale: Optional. Sets the scale of the monitor

### postcommands

postcommands defines a list of commands that will be sent to sway after the layouts have been established


### Usage 

```
  WindowArranger [option]... [configfile]
```

If no configfile is given, `WindowArranger` reads the configuration from standard input.

#### Options
```
    -dump 
    -wait uint
    -debug
```

##### Dump

WindowArranger works by sending commands to sway over ipc. With `-dump` given, ie:

```
WindowArranger -dump config.yaml
```

commands will be printed to standard out, rather than sent to sway 


##### Wait 

With the `wait` option you can instruct `WindowArranger` to wait until the windows you want to arrange are present. Say you have a criteria `title=Window1` in your configuration, then `WindowArranger` will wait until a window titled 'Window1' is present. So with:

```
WindowArranger -wait 20 configfile
```

`WindowArranger` will wait up to 20 seconds for all criteria in the configuration to find a match.

If the 20 seconds pass without all windows appearing, `WindowArranger` exits with a non-zero exit code.

##### Debug

Given the `debug` option, WindowArranger will print a bit more information when exitting with error

### Shebang

WindowArranger functions as an interpreter of config files, so you could also write your configfile with a shebang:

```
#!/usr/bin/env WindowArranger
monitors:
- name: eDP-1
  workspaces:
  - H[V['title=T1' 'title=T2'] H['title=T3' 'title=T4']]
postcommands:
- '[title=T1] resize set width 33ppt'
```

and run with just:

```
./configfile
```

### Limitations

There are certain layouts that Sway does not seem to accept. For example

```
H[V[V['title=MyWindow']]
```
will become 
```
H[V['title=MyWindow']
```

YMMV.

### Surplus windows

Open windows that are not mentioned in the configuration will be left in a workspace numbered one higher than the number of 
workspaces you've defined. So if you've defined 3 workspaces, surplus windows will end up in workspace `4`.
