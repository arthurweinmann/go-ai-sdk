#!/usr/bin/env python3

import tiktoken

import argparse
import logging
import os
import pdb
import sys
import traceback
import typing
from collections.abc import Callable


logging.basicConfig(
    level=logging.INFO, format="[%(asctime)s] [%(name)s] [%(levelname)s] %(message)s"
)
logger = logging.getLogger(__name__)


def idb_excepthook(type, value, tb):
    """Call an interactive debugger in post-mortem mode

    If you do "sys.excepthook = idb_excepthook", then an interactive debugger
    will be spawned at an unhandled exception
    """
    if hasattr(sys, "ps1") or not sys.stderr.isatty():
        sys.__excepthook__(type, value, tb)
    else:
        traceback.print_exception(type, value, tb)
        print
        pdb.pm()


def broken_pipe_handler(
    func: Callable[[typing.List[str]], int], *arguments: typing.List[str]
) -> int:
    """Handler for broken pipes

    Wrap the main() function in this to properly handle broken pipes
    without a giant nastsy backtrace.
    The EPIPE signal is sent if you run e.g. `script.py | head`.
    Wrapping the main function with this one exits cleanly if that happens.

    See <https://docs.python.org/3/library/signal.html#note-on-sigpipe>
    """
    try:
        returncode = func(*arguments)
        sys.stdout.flush()
    except BrokenPipeError:
        devnull = os.open(os.devnull, os.O_WRONLY)
        os.dup2(devnull, sys.stdout.fileno())
        # Convention is 128 + whatever the return code would otherwise be
        returncode = 128 + 1
    return returncode


def resolvepath(path):
    return os.path.realpath(os.path.normpath(os.path.expanduser(path)))


def parseargs(arguments: typing.List[str]):
    """Parse program arguments"""
    parser = argparse.ArgumentParser(description="Python command line script template")
    parser.add_argument(
        "--debug",
        "-d",
        action="store_true",
        help="Launch a debugger on unhandled exception",
    )
    # cl100k_base for gpt-4, gpt-3.5-turbo, text-embedding-ada-002; p50k_base for Codex models, text-davinci-002, text-davinci-003
    parser.add_argument("--encoding")
    parser.add_argument("--prompt")
    parsed = parser.parse_args(arguments)
    return parsed


def main(*arguments):
    """Main program"""
    parsed = parseargs(arguments[1:])
    if parsed.debug:
        sys.excepthook = idb_excepthook
    enc = tiktoken.get_encoding(parsed.encoding)
    print(f"{len(enc.encode(parsed.prompt))}")


if __name__ == "__main__":
    exitcode = broken_pipe_handler(main, *sys.argv)
    sys.exit(exitcode)