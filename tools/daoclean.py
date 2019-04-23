#!/usr/bin/env python3

import sys
import re
import click
import subprocess
from collections import namedtuple

Container = namedtuple("Container", ("id", "name"))
Containers = []


def dockerPs():
    lines = subprocess.check_output(("docker", "ps")).decode().split('\n')[1:]
    for line in lines:
        segs = line.split()
        if not segs:
            continue
        Containers.append(Container(id=segs[0], name=segs[-1]))


def containerFind(name, count):
    pat = r"\b%s_\d+\b" % name

    lis = []
    for c in Containers:
        if c.name == name or re.match(pat, c.name):
            lis.append(c)

    print(" SURVIVE :", [x.name for x in lis[:count]])
    print("  KILLED :", [x.name for x in lis[count:]])

    return lis[count:]


def containerRemove(clist):
    for c in clist:
        # unregister
        try:
            subprocess.check_output(("docker", "exec", c.id, "saybye"))
            print("SayBye: OK: ", c.name)
        except:
            print("SayBye: NG: ", c.name)

        # remove
        try:
            subprocess.check_output(("docker", "rm", "-f", c.id))
            print("Remove: OK: ", c.name)
        except:
            print("Remove: NG: ", c.name)

@click.command()
@click.argument('names', nargs=-1)
@click.option('--count', '-c', type=int, default=1, help="How many left.")
@click.option('--guess', '-g', is_flag=True, help="Guess from daoker.sh.")
def main(count, names, guess):
    names = names or []

    msName = None
    if guess:
        try:
            for line in open("./daoker.sh").readlines():
                line = line.strip()
                if line.startswith("msName="):
                    msName = line[7:].strip()
            f.close()
        except:
            pass

    if msName:
        names.append(msName)

    for name in names:
        lis = containerFind(name, count)
        containerRemove(lis)

if __name__ == "__main__":
    dockerPs()
    sys.exit(main())
