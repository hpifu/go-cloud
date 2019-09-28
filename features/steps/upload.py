#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json
import os


@when('请求 /upload, file: "{file:str}"')
def step_impl(context, file):
    res = requests.post("{}/upload/{}".format(context.config["url"], context.token), files={
        'file': open('features/assets/hatlonely.png', 'rb')
    })

    context.status = res.status_code
    context.body = res.text
    context.cookies = res.cookies
    print(res.text)
    if context.status == 200:
        context.res = json.loads(res.text)
    print({
        "status": context.status,
        "body": context.body,
        "cookies": context.cookies,
    })


@then('检查 data 目录存在文件, file: "{file:str}"')
def step_impl(context, file):
    print("{}/data/{}/{}".format(context.config["prefix"], context.id, file))
    assert_that(os.path.exists(
        "{}/data/{}/{}".format(context.config["prefix"], context.id, file)
    ), is_(True))
