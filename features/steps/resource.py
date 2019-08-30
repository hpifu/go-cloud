#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json
import os


@when('请求 /resource, file: "{file:str}"')
def step_impl(context, file):
    res = requests.get("{}/resource".format(context.config["url"]), params={
        "token": context.token,
        "name": file,
    })
    print(res.content)

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
