from crontab import CronTab
cron = CronTab(user='mac')
job = cron.new(command="/home/mac/goServer/24.sh")
job.minute.every(2)
cron.write()

# */2 * * * * /home/mac/goServer/24.sh