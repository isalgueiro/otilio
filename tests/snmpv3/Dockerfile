FROM alpine:3.6
RUN apk add --update net-snmp net-snmp-tools 
RUN sed -i 's/agentAddress  udp:127.0.0.1:161/agentAddress  udp:10161/g' /etc/snmp/snmpd.conf
RUN net-snmp-create-v3-user -ro -A theauthpassword -a SHA -X theprivacyencryptionpassword -x DES theuser
EXPOSE 10161/tcp
RUN rm -rf /var/cache/apk/*
ENTRYPOINT [ "snmpd", "-f" ]
