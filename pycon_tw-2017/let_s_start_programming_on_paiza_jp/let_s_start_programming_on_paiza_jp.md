Let's start programming on paiza.jp
===

這些都是 online judgement system
https://uva.onlinejudge.org/
https://leetcode.com/
https://paiza.jp/logic_summoner
https://paiza.jp/poh/hatsukoi

大guy4這樣,認同請分享

```python=
# coding: utf-8
# 自分の得意な言語で
# Let's チャレンジ！！
s = int(input())
t = int(input())
print(''.join(('-', '+')[i+1==t] for i in range(s)))
```

```python=
# coding: utf-8
# 自分の得意な言語で
# Let's チャレンジ！！

import sys

DECREASE_TEMP = 1
DECREASE_COST = 2
HOLD_COST = 1
INCREASED_TEMP_TABLE = {
    'in': 5,
    'out': 3,
}

def update_status(prev_hour, current_hour, prev_temperature, prev_cost):
    duration = current_hour-prev_hour
    duration_to_zero = prev_temperature//DECREASE_TEMP

    temperature = prev_temperature-DECREASE_TEMP*min(duration, duration_to_zero)
    cost = prev_cost+DECREASE_COST*min(duration, duration_to_zero)+HOLD_COST*max(duration-duration_to_zero, 0)
    #print((prev_hour, current_hour, duration, duration_to_zero, temperature, prev_cost, cost))
    return (current_hour, temperature, cost)

def main(argv=sys.argv[:]):
    n = int(input())

    prev_hour = 0
    temperature = 0
    cost = 0
    for i in range(n):
        (hour_string, inout) = input().split()
        hour = int(hour_string)

        (prev_hour, temperature, cost) = update_status(prev_hour, hour, temperature, cost)
        temperature += INCREASED_TEMP_TABLE[inout]

    (prev_hour, temperature, cost) = update_status(prev_hour, 24, temperature, cost)
    print(cost)
    return 0

if __name__ == '__main__':
    sys.exit(main())
```

