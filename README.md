@author winstark

## Linux's login logs

- `utmp` and `wtmp` is binary file contains logs for logins and logouts in Linux system.  
- `/var/run/utmp` This file contains information about the users who are currently logged onto the system.  
- `who` command users `utmp` file to display logged in users
- `/var/log/wtmp`  This file is like history for utmp file, i.e. it maintains the logs of all logged in and logged out users (in the past).  
- The `last` command users `wtmp` file to display listing of last logged in users.  


## Crontab Linux

- **cron** is a utility that allows tasks to automatically run in the background of the system at regular intervals by use of the **cron deamon**.  
- `crond` deamon is the background service that enables cron funtionality.  
- **Crontab(Cron table)** is a file which contains the schedule of **cron** entries to be run and at what times they are to be run.  
- **Cron** checks these files and directories:  
  - `/etc/crontab`: system crontab.Originally it was usually used to run daily, weekly, monthly jobs.    
  - `/etc/cron.d`: directory that contains systems cronjobs stored for different users.  
  - `/var/spool/cron`: directory that contains user corntables created by the `crontab` command.  
 [more](https://opensource.com/article/17/11/how-use-cron-linux) 

## /proc filesystem  

The `/proc` filesystem contains a illusionary filesystem. It does not exist on a disk. Instead, the kernerl creates it in memory. It is used to provide information about the system(originally about processes, hence the name).  

