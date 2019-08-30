#!/usr/bin/env python3

from behave import *
from hamcrest import *


@then('检查状态码 res.status_code: {status:int}')
def step_impl(context, status):
    assert_that(context.status, equal_to(status))


@then('检查返回包体 res.body，包含字符串 "{message:str}"')
def step_impl(context, message):
    assert_that(context.body, contains_string(message))
