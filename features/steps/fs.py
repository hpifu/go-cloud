#!/usr/bin/env python3


from behave import *
from hamcrest import *
import requests
import os


@then('fs 检查文件存在 "{file:str}"')
def step_impl(context, file):
    assert_that(os.path.exists(file), is_(True))
