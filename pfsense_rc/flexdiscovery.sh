#!/bin/sh

rc_start() {
	/usr/local/etc/rc.d/flexdiscovery onestart
}

rc_stop() {
	/usr/local/etc/rc.d/flexdiscovery onestop
}

case $1 in
	start)
		rc_start
		;;
	stop)
		rc_stop
		;;
	restart)
		rc_stop
		rc_start
		;;
esac