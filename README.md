# Sleep-on-Lan

Built in go. Listens for the "magic packet" and performs a shutdown if received. 

Contains the options:
  - --log-path string
      - File to log to (default is stdout)
  - --port string
      - Set the WoL listen port (default is UDP/9) (default "9")
  - --test-mode 
      - Don't poweroff, just log attempts


Example:
  - sleep-on-lan --log-path=/var/log/sleep-on-lan.log --port=9 --test-mode=true


To start on boot:
  - Build the binary with `go build` 
  - Move binary with `mv sleep-on-lan /usr/local/bin/sleep-on-lan`
  - Move systemd unit file with `mv sleep-on-lan.service /etc/system/systemd/sleep-on-lan.service`
  - Run `systemctl daemon-reload`
  - Enable on boot with `systemctl enable sleep-on-lan`
  - Start immediately with `systemctl start sleep-on-lan`
