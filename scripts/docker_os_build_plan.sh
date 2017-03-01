#!/bin/bash

echo "This script will pre-configure the server to run docker"
DOCKER_USER_INPUT=$1
DOCKER_PUBKEY=$2
DOCKER_HOSTNAME=$3
DOCKER_PROXY=$4
PROXY_ENABLE=$5

if [ -z "${DOCKER_PUBKEY}" ]; then
  echo "ERROR : this script requires a public key for docker user!"
  echo "USAGE: $0 <docker user> '<public key>'"
  exit 1
fi

DOCKER_USER=${DOCKER_USER_INPUT:-"docker"}

# optionally set some persistent proxy server configuration
if [ "${PROXY_ENABLE}" = "true" ]; then
cat > "/etc/environment" << EOF
${DOCKER_PROXY}
EOF
echo "Completed update to /etc/environment, $?"
fi

# create a service account
if [ "$DOCKER_USER" = "root" ]; then
  echo "WARNING : docker-engine user should not be configured as root on bare metal systems, ${DOCKER_USER}."
else
  grep "^${DOCKER_USER}:" /etc/passwd || useradd "${DOCKER_USER}" -d "/home/${DOCKER_USER}"
  echo "Completed adding user account ${DOCKER_USER}, $?"
fi

# setup .ssh folderls -ak
if [ ! -d "/home/${DOCKER_USER}/.ssh" ]; then
  mkdir -p "/home/${DOCKER_USER}/.ssh"
  chmod 700 "/home/${DOCKER_USER}/.ssh"
  chown "${DOCKER_USER}:${DOCKER_USER}" "/home/${DOCKER_USER}/.ssh"
  echo "Completed updating permissions and folders for /home/${DOCKER_USER}/.ssh, $?"
fi
if [ ! -f "/home/${DOCKER_USER}/.ssh/authorized_keys" ] ; then
  touch "/home/${DOCKER_USER}/.ssh/authorized_keys"
  chmod 600 "/home/${DOCKER_USER}/.ssh/authorized_keys"
  chown "${DOCKER_USER}:${DOCKER_USER}" "/home/${DOCKER_USER}/.ssh/authorized_keys"
  echo "Completed updating permissions for /home/${DOCKER_USER}/.ssh/authorized_keys, $?"
fi
grep "${DOCKER_PUBKEY}" "/home/${DOCKER_USER}/.ssh/authorized_keys" || cat >> "/home/${DOCKER_USER}/.ssh/authorized_keys" << EOF
${DOCKER_PUBKEY}
EOF

# modify /home/{user}/.bash_profile to set a persistent proxy
if [ "${PROXY_ENABLE}" = "true" ]; then
cat >> "/home/${DOCKER_USER}/.bash_profile" << EOF
${DOCKER_PROXY}
EOF
fi

# give ansible ssh sudo ability to execute without tty
sed -i 's/.*Defaults.*requiretty/\#Defaults requiretty/g' /etc/sudoers

# give sudoers access
cat >> "/etc/sudoers.d/90-${DOCKER_USER}" << SUDOERS_EOF
# User rules for icsp docker user
${DOCKER_USER} ALL=(ALL) NOPASSWD:ALL
SUDOERS_EOF
echo "Completed updating permissions for sudoers on user ${DOCKER_USER}, $?"

sed -i "s/localhost.localdomain/${DOCKER_HOSTNAME}/g" /etc/hostname
echo "Completed hostname update : $(cat /etc/hostname), $?"

echo "docker host provisioned by docker-machine oneview driver" >> /etc/motd
echo "docker customizations complete"

exit 0
