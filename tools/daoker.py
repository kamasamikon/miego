#!/usr/bin/env python3

import sys
import os
import click
import subprocess
import shlex
import time

#
# - 生成Dockerfile
# - 根据当前的工程Git，搜集一些信息
# - 生成ms/msa.cfg文件
# - copy main ms/
# - docker build -t <msname> .
#


dfTempl = '''
# FROM %s
FROM alpine:3.9

WORKDIR /root
COPY ./PRC /etc/localtime
ENTRYPOINT /root/main

# >>> User part >>>
%s
# <<< User part <<<

COPY ms/ ./
'''


# Run
_foreground = None
_container = None
_shareMode = None
_appendMode = None
_kill = True

# Service
_msName = None
_msUpstream = None
_msKind = None
_msVern = None
_msPort = None
_msDesc = None

# Dockerfile
_dfuser = None

# MSA
_msaBase = None

# Docker environ
_env = None

# extra options, -x '-v aaa:bbb'
_extra = []


def headVersion():
    try:
        cmd = ("git", "log", "-n", "1")
        return subprocess.check_output(cmd).strip().decode().split()[1]
    except:
        return "NA"

def isUpdated():
    try:
        cmd = ("git", "status", "-uno", "-s")
        return not bool(subprocess.check_output(cmd))
    except:
        return False

def createMsaCfg():
    lines = []

    # Service Info
    lines.append("# Service information")
    lines.append("s:/ms/name=%s" % _msName)
    lines.append("s:/ms/upstream=%s" % _msUpstream)
    lines.append("s:/ms/kind=%s" % _msKind)
    lines.append("s:/ms/version=%s" % _msVern)
    lines.append("i:/ms/port=%s" % _msPort)
    lines.append("s:/ms/desc=%s" % _msDesc)
    lines.append("s:/ms/url/path=/%s" % _msName)
    lines.append("")

    # MSB Info
    lines.append("# MSB information")
    lines.append("i:/msb/regWait/ok=5")
    lines.append("i:/msb/regWait/ng=1")
    lines.append("")

    # Project Info
    lines.append("# Project information")
    lines.append("s:/build/dirname=%s" % os.path.basename(os.getcwd()))
    lines.append("s:/build/time=%s" % time.asctime())
    lines.append("s:/build/version=%s" % headVersion())
    lines.append("i:/build/updated=%d" % isUpdated())

    # This is /root/msa.cfg, it will be used fore register
    with open("ms/msa.cfg", "a") as f:
        f.writelines(["\r\n" + l for l in lines])


def createDockerfile():
    text = dfTempl % (_msaBase, _dfuser)
    with open("Dockerfile", "w") as f:
        f.write(text)

def colorprint(s):
    print('\033[1;3{}m{}\033[0m'.format(2, s))

def saferun(cmd, debug=True):
    try:
        if debug:
            color = 3
            cmdline = " ".join(cmd)
            print('\033[1;3{}m{}\033[0m'.format(color, cmdline))

        subprocess.run(cmd)
        return True
    except:
        return False


def callUserScript():
    saferun(["rm", "-fr", "ms"])
    saferun(["mkdir", "-p", "ms"])
    x = subprocess.run(["sh", "-e", "-u", "./.userScript"])
    if x.returncode != 0:
        sys.exit(x.returncode)

def copyMain():
    saferun(["cp", "-fr", "main", "ms"])

def build():
    '''Generate the docker image'''

    callUserScript()
    createMsaCfg()
    createDockerfile()
    copyMain()

    cmd = ["sudo", "docker", "build", "--no-cache", "-t", _msName, "."]
    for e in _extra:
        segs = shlex.split(e)
        cmd.extend(segs)

    saferun(cmd)

def dockerGateway():
    cmd = ("sudo", "docker", "network", "inspect", "bridge", "--format", '{{(index .IPAM.Config 0).Gateway}}')
    print(">>> ", " ".join(cmd))
    return subprocess.check_output(cmd).strip().decode("utf-8")

@click.command()

# Run
@click.option('--foreground', '-f', is_flag=True, help="(False):   Run docker foreground.")
@click.option('--container', '-c', help="($msName): Container name.")
@click.option('--share-mode', '-s', is_flag=True, help="(False):   -v PWD/ms:/root/ms.")
@click.option('--append-mode', '-a', is_flag=True, help="(False):   New contaner.")
@click.option('--kill', '-k', is_flag=True, type=bool, default=False, help="(False):   Kill old container.")

# Service
@click.option('--ms-name', '-n', help="(demo):    Service Name.")
@click.option('--ms-upstream', '-u', help="():        Upstream name in nginx.conf.")
@click.option('--ms-kind', '-t', help="(http):    Service Type, grpc or http.")
@click.option('--ms-vern', '-v', help="(v1):      Service Version.")
@click.option('--ms-port', '-p', help="(8888):    Service Port.")
@click.option('--ms-desc', '-d', help="(null):    Service Description.")
@click.option('--guess', '-g', is_flag=True, help="Guess from daoker.sh and Makefile.")

# MSB
@click.option('--msb-name', '-m', help="(msb):     MSB container name")
@click.option('--msb-ip', '-i', help="(byGuess): MSB ip address.")

# MSA
@click.option('--msa-base', '-b', help="(msa):     MSA base image")

# Docker environ
@click.option('--env', '-e', help="(null):    Environ passed to docker.", multiple=True)

# extra docker options
@click.option('--extra', '-x', help="(null):    extra docker options.", multiple=True)

def main(foreground, container, share_mode, append_mode, kill,
        ms_name, ms_upstream, ms_kind, ms_vern, ms_port, ms_desc,
        guess,
        msb_name, msb_ip,
        msa_base,
        env,
        extra,
        ):
    '''CMDS: build|b=build, run|r=run.'''

    foreground = foreground
    container = container
    shareMode = share_mode
    appendMode = append_mode
    kill = kill
    msName = ms_name
    msUpstream = ms_upstream
    msKind = ms_kind
    msVern = ms_vern
    msPort = ms_port
    msDesc = ms_desc
    guess = guess
    msbName = msb_name
    msaBase = msa_base
    env = env
    extra = extra

    #
    # Set global
    #

    # Run
    global _foreground, _container, _shareMode, _appendMode, _kill
    _foreground = foreground
    _container = container
    _shareMode = shareMode
    _appendMode = appendMode
    _kill = kill

    # Service
    global _msName, _msUpstream, _msKind, _msVern, _msPort, _msDesc
    _msName = msName or "demo"
    _msUpstream = msUpstream or ""
    _msKind = msKind or "http"
    _msVern = msVern or "v1"
    _msPort = msPort or 8888
    _msDesc = msDesc or ""
    if not msName or guess:
        for guessfile in ("./daoker.sh", "Makefile"):
            try:
                for line in open(guessfile).readlines():
                    line = line.strip()
                    if line.startswith("msName="):
                        _msName = line[7:].strip()
                    if line.startswith("msVern="):
                        _msVern = line[7:].strip()
                    if line.startswith("msPort="):
                        _msPort = line[7:].strip()
                    if line.startswith("msDesc="):
                        _msDesc = line[7:].strip()
                f.close()
            except:
                pass


    # Dockerfile
    global _dfuser
    try:
        _dfuser = open(".userDockerCommand", "r").read()
    except:
        _dfuser = ""

    # MSA
    global _msaBase
    _msaBase = msaBase or "msa-alpine"

    # Docker environ
    global _env
    _env = env or []

    # Docker environ
    global _extra
    _extra = extra or []

    #
    # Go
    #
    build()

if __name__ == "__main__":
    sys.exit(main())
