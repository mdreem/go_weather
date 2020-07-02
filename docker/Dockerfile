FROM jboss/keycloak:10.0.2

COPY ./data/keycloak-export.json /opt/jboss/keycloak-export.json
COPY ./data/keycloak_entry.sh /opt/jboss/keycloak_entry.sh

ENTRYPOINT /opt/jboss/keycloak_entry.sh
