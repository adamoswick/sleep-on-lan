# Sleep-on-Lan

Built in go. Listens for the "magic packet" and performs a shutdown if received. 

Contains the options:
  --log-path string
      File to log to (default is stdout)
  --port string
      Set the WoL listen port (default is UDP/9) (default "9")
  --test-mode
      Don't poweroff, just log attempts

Example:
  - sleep-on-lan --log-path=/var/log/sleep-on-lan.log --port=9 --test-mode=true
