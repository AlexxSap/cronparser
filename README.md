# cronparser

Simple cron parser written in go
you can read about in here [challenge-cron](https://codingchallenges.fyi/challenges/challenge-cron)

it look something like this:
```
0 9 * * 6 <command>
```

Which can be interpreted using the following:
```
# ┌───────────── minute (0–59)
# │ ┌───────────── hour (0–23)
# │ │ ┌───────────── day of the month (1–31)
# │ │ │ ┌───────────── month (1–12)
# │ │ │ │ ┌───────────── day of the week (0–6) (Sunday to Saturday;
# │ │ │ │ │                                   7 is Sunday on some systems)
# │ │ │ │ │
# │ │ │ │ │
# * * * * * <command>
```
For these fields the allowed values are:

| Field | Required | Allowed values | Allowed special characters|
|:------|:------|:------|:------|
| Minutes | Yes | 0–59 |* , -|
| Hours | Yes | 0–23 |* , -|| Day of month | Yes | 1–31 |* , -|
| Month | Yes | 1–12 |* , -|| Day of week | Yes | 0–6 |* , - |

## Examples
`* * * * *` - every minute  
`* 9 * * SAT` - every minute, between 9:00 and 9:59, on Saturday  
`5 14-15 * 1,3 Mon-Fri` - 14:05 and 15:05, on yan and march, from monday to friday  


Input  
string and datetime  
Output  
bool with result - this datetime match with string or not  
