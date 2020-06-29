#!/bin/bash

set -eu

sed -i  '$ a JAVA_OPTS="$JAVA_OPTS -Dkeycloak.migration.action=import -Dkeycloak.migration.provider=singleFile -Dkeycloak.migration.file=/opt/jboss/keycloak-export.json -Dkeycloak.profile.feature.token_exchange=enabled -Dkeycloak.profile.feature.admin_fine_grained_authz=enabled"' /opt/jboss/keycloak/bin/standalone.conf

/opt/jboss/tools/docker-entrypoint.sh
