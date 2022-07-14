from crontab import CronTab
cron = CronTab(user='mac')
job = cron.new(command="echo 'hello_world'")
job.minute.every(1)
cron.write()