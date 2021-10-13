#!/usr/bin/env python3
# -*- coding: utf-8 -*-


# from https://alchemypy.com/2020/05/13/speed-your-python-with-lua/


import sys
import time

from lupa import LuaRuntime


def sum_all_to_range_pure_python(end_range_number):
    start_time = time.time()
    the_sum = 0
    for number in range(1, end_range_number+1):
        the_sum += number
    stop_time = time.time()
    return (end_range_number, stop_time - start_time, the_sum)


def sum_all_to_range_with_lua_python(end_range_number):
    start_time = time.time()

    lua_code =  """
    function(n)
       sum = 0
       for i=1,n,1 do
           sum = sum + i
       end
       return sum
    end
    """
    lua_func = LuaRuntime(encoding=None).eval(lua_code)
    the_sum = lua_func(end_range_number)
    stop_time = time.time()
    return(end_range_number, stop_time - start_time, the_sum)


def main(argv=sys.argv[:]):
    end_range_number = 200000000

    python_result = sum_all_to_range_pure_python(end_range_number)
    lupa_result = sum_all_to_range_with_lua_python(end_range_number)

    print('python - range: %d, time: %g sec, sum: %d' % python_result)
    print('lupa - range: %d, time: %g sec, sum: %d' % lupa_result)


if __name__ == '__main__':
    sys.exit(main())
