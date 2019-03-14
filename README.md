@author winstark

## Linux's login logs

- `utmp` and `wtmp` is binary file contains logs for logins and logouts in Linux system.  
- `/var/run/utmp` This file contains information about the users who are currently logged onto the system.  
- `who` command users `utmp` file to display logged in users
- `/var/log/wtmp`  This file is like history for utmp file, i.e. it maintains the logs of all logged in and logged out users (in the past).  
- The `last` command users `wtmp` file to display listing of last logged in users.  
