#!/usr/bin/env python3

import sys
import os
import click
import subprocess
import shlex
import time


dfTempl = '''
FROM msa

# >>> User part >>>

%s
# <<< User part <<<

COPY ms/ ./ms
COPY msa.cfg ./ms
'''


# Run
_foreground = None
_container = None
_sharemode = None
_appendmode = None
_kill = True

# Service
_msname = None
_msvern = None
_msport = None
_msdesc = None

# Dockerfile
_dfuser = None

# MSB
_msbname = None
_msbip = None

# Docker environ
_env = None

# build or run
_cmds = None

# extra options, -x '-v aaa:bbb'
_extra = []


def headVersion():
    try:
        cmd = ("git", "show-ref", "HEAD")
        return subprocess.check_output(cmd).strip().decode().split()[0]
    except:
        return "NA"

def isUpdated():
    try:
        cmd = ("git", "status", "-uno")
        return len(subprocess.check_output(cmd).strip().decode().split("\n")) < 5
    except:
        return False


def getMsbIp(msbName=None):
    msbName = msbName or "msb"
    cmd = ("sudo", "docker", "inspect",
            "--format", '{{ .NetworkSettings.IPAddress }}', msbName)
    print(">>> ", " ".join(cmd))
    return subprocess.check_output(cmd).strip().decode("utf-8")


def createMsaCfg():
    lines = []

    # Service Info
    lines.append("s:/ms/name=%s" % _msname)
    lines.append("s:/ms/version=%s" % _msvern)
    lines.append("i:/ms/port=%s" % _msport)
    lines.append("s:/ms/desc=%s" % _msdesc)

    # MSB Info
    msbip = _msbip or getMsbIp(_msbname)
    lines.append("s:/msb/host=%s" % msbip)
    lines.append("i:/msb/regWait/ok=5")
    lines.append("i:/msb/regWait/ng=1")

    # Project Info
    lines.append("s:/build/dirname=%s" % os.path.basename(os.getcwd()))
    lines.append("s:/build/time=%s" % time.asctime())
    lines.append("s:/build/version=%s" % headVersion())
    lines.append("i:/build/updated=%d" % isUpdated())

    # This is /root/msa.cfg, it will be used fore register
    with open("msa.cfg", "w") as f:
        f.writelines([l + "\r\n" for l in lines])


def createDockerfile():
    text = dfTempl % _dfuser
    with open("Dockerfile", "w") as f:
        f.write(text)

def callUserScript():
    try:
        os.removedirs("ms")
    except:
        pass

    try:
        cmd = ["sh", "./userScript"]
        subprocess.run(cmd)
    except:
        pass

def copyMain():
    try:
        cmd = ["cp", "-frv", "main", "ms"]
        subprocess.run(cmd)
    except:
        pass

def build():
    '''Generate the docker image'''

    createMsaCfg()
    createDockerfile()
    callUserScript()
    copyMain()

    cmd = ["sudo", "docker", "build", "-t", _msname, "."]
    for e in _extra:
        segs = shlex.split(e)
        cmd.extend(segs)

    print(">>> ", " ".join(cmd))
    subprocess.run(cmd)


def run():
    '''Run docker image'''

    foreground = _foreground
    container = _container
    sharemode = _sharemode

    msbip = _msbip or getMsbIp(_msbname)
    container = container or _msname

    #
    # Remove old container
    #
    if _kill:
        cmd = ("sudo", "docker", "exec", container, "saybye")
        print(">>> ", " ".join(cmd))
        subprocess.run(cmd)

        cmd = ("sudo", "docker", "rm", "--force", container)
        print(">>> ", " ".join(cmd))
        subprocess.run(cmd)

    #
    # Run container
    #
    if _appendmode:
        index = 0
        name = container
        while True:
            name = container if index == 0 else container + "_%d" % index
            index += 1
            cmd = ["sudo", "docker", "ps", "-aq", "--filter", "name=^/%s$" % name]
            print(">>> ", " ".join(cmd))
            if not subprocess.check_output(cmd):
                break
        container = name
        print(container)

    # -v: ms: conf.Load("/tmp/conf/main.cfg")
    cmd = ["sudo", "docker", "run", "-it", "--name", container, "-v", "/tmp/.conf.%s:/tmp/conf" % container]
    for e in _env:
        cmd.append("-e")
        cmd.append(e)

    for e in _extra:
        segs = shlex.split(e)
        cmd.extend(segs)

    if not _foreground:
        cmd.append("-d")

    if sharemode:
        cmd.extend(("-v", os.getcwd() + "/ms:/root/ms"))

    cmd.extend(("-e", "MSBHOST=%s" % msbip))
    cmd.append("%s:latest" % _msname)
    print(">>> ", " ".join(cmd))
    subprocess.run(cmd)


@click.command()

# Run
@click.option('--foreground', '-f', is_flag=True, help="Run docker foreground.")
@click.option('--container', '-c', help="Container name, default is demo.")
@click.option('--sharemode', '-s', is_flag=True, help="-v PWD/ms:/root/ms.")
@click.option('--appendmode', '-a', is_flag=True, help="New contaner.")
@click.option('--kill', '-k', is_flag=True, type=bool, default=False, help="Kill old container.")

# Service
@click.option('--msname', '-n', help="Service Name.")
@click.option('--msvern', '-v', help="Service Version.")
@click.option('--msport', '-p', help="Service Port.")
@click.option('--msdesc', '-d', help="Service Description.")

# Dockerfile
@click.option('--dfuser', '-D', help="Commands in Dockerfile.")

# MSB
@click.option('--msbname', '-m', help="MSB container name. Default is 'msb'.")
@click.option('--msbip', '-i', help="MSB ip address.")

# Docker environ
@click.option('--env', '-e', help="environ passed to docker.", multiple=True)

# build or run
@click.argument('cmds', nargs=-1)

# extra docker options
@click.option('--extra', '-x', help="extra docker options.", multiple=True)

def main(foreground, container, sharemode, appendmode, kill,
        msname, msvern, msport, msdesc, dfuser,
        msbname, msbip,
        env,
        extra,
        cmds):


    #
    # Set global
    #

    # Run
    global _foreground, _container, _sharemode, _appendmode, _kill
    _foreground = foreground
    _container = container
    _sharemode = sharemode
    _appendmode = appendmode
    _kill = kill

    # Service
    global _msname, _msvern, _msport, _msdesc
    _msname = msname or "demo"
    _msvern = msvern or "v1"
    _msport = msport or 8888
    _msdesc = msdesc or ""

    # Dockerfile
    global _dfuser
    _dfuser = open(dfuser, "r").read() if dfuser else ""

    # MSB
    global _msbname, _msbip
    _msbname = msbname or "msb"
    _msbip = msbip

    # Docker environ
    global _env
    _env = env or []

    # Docker environ
    global _extra
    _extra = extra or []

    # build or run
    global _cmds
    _cmds = cmds or ["build", "run"]

    #
    # Go
    #
    for cmd in _cmds:
        if cmd in ("build", "b"):
            build()
            continue

        if cmd in ("run", "r"):
            run()
            continue

if __name__ == "__main__":
    sys.exit(main())
