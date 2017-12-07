#!/bin/bash
set -x
OUT=$(nsenter --mount=/host/proc/1/ns/mnt -- nsenter -t 1 -n findmnt --version) && \
OUT=$(nsenter --mount=/host/proc/1/ns/mnt -- nsenter -t 1 -n curl --version) && \
OUT=$(nsenter --mount=/host/proc/1/ns/mnt -- nsenter -t 1 -n blkid) && \
OUT=$(nsenter --mount=/host/proc/1/ns/mnt -- nsenter -t 1 -n mkfs.ext4 -V) && \
exit 0
