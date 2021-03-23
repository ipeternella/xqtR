# 🛠️ xqtR (executoR) 🛠️

`xqtR` (short for executor) is a command line tool to execute a series of jobs specified by job yaml files. The steps in a job can be run concurrently (by spawning goroutines in a workpool) or in the old fashioned sync way (default).

A quick demo of running a `job.yaml` which contains a job whose steps will be run by a single main goroutine (sync):

![xqtR sync demo](docs/demos/xqtr-sync-demo.gif)

Now, a similar `job.yaml` which runs the same steps but spawns goroutines (async) to run the steps in parallel (when possible) according to `num_workers` in the yaml file:

![xqtR async demo](docs/demos/xqtr-async-demo.gif)

This project is a WIP.

## Table of Contents

- [Introduction](#Introduction)
  - [How does it work](#How-does-it-work)
- [How to Use](#How-to-use)
  - [Job Files](#Job-yaml-files)

## Introduction

`xqtR` (short for executor) is a command line tool that parses `job.yaml` files, which have, perhaps, similar syntax to github action yaml files, in order
to execute a series of jobs and steps on your machine. Optionally, as mentioned below, the jobs can be executed concurrently by spawning goroutines.

Hence, this project is inspired by modern CICD tools such as `Azure DevOps` and `Github Actions` that use `yaml` files to configure a sequence of steps known
as pipelines. Naturally, this is an ultra-and-I-really-mean-it-simplifed-version of these famous yaml parsers to run jobs! And yes, this project was also used
to explore some concurrency in GO.

However, besides being simple, this tool can be useful in cases like configuring new machines in which a sequence of programs must be downloaded and installed
(and this can be done in concurrently here, which can speed up the process!).

## How-does-it-work

`xqtR` parses yaml files looking for **jobs** that are composed of a series of **steps**. Each step contains commands that, in the current version, are read and used
to spawn new OS processes which will invoke the system shell (currently only `bash` -- so Windows users be warned for now) in order to execute them.

The spawned processes also contains specific os pipes that captures their `stdstreams` (`stdout` and `stderr`) which can be used to display errors, warnings, or the
command's `stdout` when the tool's `--log-level` is set to `debug`.

## How-to-use

Using the tool is pretty straightforward, just open up your favorite shell and invoke the `help` command:

```sh
xqtr -h
```

And (hopefully) a nice explanation will tell you how to use it! The main command to use is `xqtr run`, which can also be used to
get some info about how to use it:

```sh
xqtr run -h
```

In the simplest form, in the case you have a `job.yaml` file on the same directory of `xqtr` binary file, simply run:

```sh
xqtr run  # same as xqtr run --file ./job.yaml (looks for a job.yaml file)
```

And that's it!

### Job yaml files

The yaml files are used to describe jobs that will be executed by `xqtR`! The main concepts are as follows:

- A job yaml file is composed of **jobs**
- Jobs are composed of **steps**
- Steps are composed of **run instructions**
- Run instructions are commands that will be executed on your system shell (bash only for now)

Here's a `job.yaml` example:

```yaml
jobs:
  - title: install macOS software providers
    steps:
      - name: brew update
        run: brew update

  - title: Some silly file creation
    steps:
      - name: echoing into a file
        run: echo "hello world" >> hello-world.txt

  - title: install apps with brew concurrently
    num_workers: 3
    steps:
      - name: install httpie
        run: brew install httpie
      - name: install lynx (terminal browser)
        run: brew install lynx
      - name: echo something into a file
        run: echo "transfer this" >> my_text.txt
```

The **jobs** are executed in a sequence (never in parallel). However, a job, when being executed, can run its **steps** in a worker pool (group of goroutines) defined by the param `num_workers`. Here's the `job.yaml` used in the `sync` demonstration gif on the beginning of this video:

```yaml
jobs:
  - title: running sync sleeps
    steps:
      - name: sleep 1
        run: sleep 1 && echo "slept for 1s"
      - name: sleep 2
        run: sleep 2 && echo "slept for 2s"
      - name: sleep 3
        run: sleep 3 && echo "slept for 3s"
```

The same steps can be run concurrently by specifying the key `num_workers` which maps to the number of goroutines that will be spawned (usually 3 is a good starting point):

```yaml
jobs:
  - title: running async sleeps
    num_workers: 3 # xqtR will spawn 3 goroutines to handle these steps!
    steps:
      - name: sleep 1
        run: sleep 1 && echo "slept for 1s"
      - name: sleep 2
        run: sleep 2 && echo "slept for 2s"
      - name: sleep 3
        run: sleep 3 && echo "slept for 3s"
```

Soon, all keys will be documented properly (WIP)!
