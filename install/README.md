# Installation of the Proxy and UI on a web-server

The proxy can be installed as a service on a web-server.

The following instructions can be used to set up systemd services on a Linux server using systemd and Apache or Nginx.

## UI
See [UI deploy.sh script](https://github.com/RoboCup-SSL/ssl-status-board-client/blob/master/deploy.sh) for
instructions on how to build and copy the UI resources to the server. Make sure to adapt the path accordingly in the 
following configs. The default path is `/srv/www/status-board`.

## Proxy
Build the proxy and copy or link the `ssl-status-board-proxy` binary to `/usr/local/bin/ssl-status-board-proxy`.

## Systemd config
1. Create a `status-board` user and a `nogroup` group, or adapt the values in the `.service` files accordingly:
   * `adduser --no-create-home --system --disabled-login --disabled-password status-board`
1. Change the default credentials in the `.yaml` files to something more secure.
1. Review the remaining configs like port and path in the `.yaml` files
1. Copy the files in /etc/systemd and /etc/ssl-status-board on the server.
1. Enable services, as desired: 
   * `systemctl enable ssl-status-board-proxy@referee-field-a.service`
   * `systemctl enable ssl-status-board-proxy@referee-field-b.service`
   * `systemctl enable ssl-status-board-proxy@vision-field-a.service`
   * `systemctl enable ssl-status-board-proxy@vision-field-b.service`
1. Start services:
   * `systemctl start ssl-status-board-proxy@referee-field-a.service`
   * `systemctl start ssl-status-board-proxy@referee-field-b.service`
   * `systemctl start ssl-status-board-proxy@vision-field-a.service`
   * `systemctl start ssl-status-board-proxy@vision-field-b.service`
1. Check if the service has been started successfully with e.g. `systemctl status ssl-status-board-proxy@vision-field-b.service`

## Apache config
1. Adapt the paths and ports as you like in `etc/apache2/conf-available/status-board.conf`
1. Copy the file to the server
1. Enable the config file: a2enconf status-board
1. Reload Apache2: systemctl reload apache2

## Nginx config
1. Adapt paths and ports as you like in `etc/nginx/sites-available/ssl-status-board.conf`
1. Include the config into your existing config
1. Reload nginx: systemctl reload nginx