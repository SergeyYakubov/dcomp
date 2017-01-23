#!/usr/bin/env bash

#SBATCH --ntasks=${DCOMP_NCPUS}
#SBATCH --nodes=${DCOMP_NNODES}
#SBATCH --cpus-per-task=1
##SBATCH --partition=all
#SBATCH -t 00:40:00
#SBATCH --workdir=${DCOMP_WORKDIR}
#SBATCH --uid=${DCOMP_UID}
#SBATCH --gid=${DCOMP_GID}

if [ "0${DCOMP_NNODES}" -gt "1" ] || [ "0${DCOMP_NCPUS}" -gt "1" ] ; then
    dockercluster -u ${DCOMP_IMAGE_NAME} ${DCOMP_DOCKER_ARGS}
    dockerexec ${DCOMP_SCRIPT}
else
    dockerrun ${DCOMP_DOCKER_ARGS} ${DCOMP_IMAGE_NAME} ${DCOMP_SCRIPT}
fi
