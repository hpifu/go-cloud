#!/usr/bin/env python3

import pymysql
import redis
import subprocess
import time
import requests
import datetime
import json
import socket
from behave import *


register_type(int=int)
register_type(str=lambda x: x if x != "N/A" else "")
register_type(bool=lambda x: True if x == "true" else False)


config = {
    "prefix": "output/cloud",
    "service": {
        "port": 16061,
        "secure": False,
        "domain": "127.0.0.1"
    },
    "api": {
        "account": "test-go-account:16060",
    },
    "mysqldb": {
        "host": "test-mysql",
        "port": 3306,
        "user": "hatlonely",
        "password": "keaiduo1",
        "db": "hads"
    },
    "redis": {
        "host": "test-redis",
        "port": 6379
    },
    "account": {
        "phone": "13145678901",
        "email": "hatlonely1@foxmail.com",
        "password": "12345678",
        "firstName": "爽",
        "lastName": "郑",
        "birthday": "1992-01-01",
        "gender": 1
    }
}


def wait_for_port(port, host="localhost", timeout=5.0):
    start_time = time.perf_counter()
    while True:
        try:
            with socket.create_connection((host, port), timeout=timeout):
                break
        except OSError as ex:
            time.sleep(0.01)
            if time.perf_counter() - start_time >= timeout:
                raise TimeoutError("Waited too long for the port {} on host {} to start accepting connections.".format(
                    port, host
                )) from ex


def deploy():
    fp = open("{}/configs/cloud.json".format(config["prefix"]))
    cf = json.loads(fp.read())
    fp.close()
    cf["service"]["port"] = ":{}".format(config["service"]["port"])
    cf["service"]["secure"] = config["service"]["secure"]
    cf["service"]["domain"] = config["service"]["domain"]
    cf["api"]["account"] = config["api"]["account"]
    print(cf)
    fp = open("{}/configs/cloud.json".format(config["prefix"]), "w")
    fp.write(json.dumps(cf, indent=4))
    fp.close()


def start():
    subprocess.Popen(
        "cd {} && nohup bin/cloud &".format(config["prefix"]),  shell=True
    )

    wait_for_port(config["service"]["port"], timeout=5)


def stop():
    subprocess.getstatusoutput(
        "ps aux | grep bin/cloud | grep -v grep | awk '{print $2}' | xargs kill"
    )


def before_all(context):
    config["url"] = "http://127.0.0.1:{}".format(config["service"]["port"])
    deploy()
    start()
    context.config = config
    context.mysql_conn = pymysql.connect(
        host=config["mysqldb"]["host"],
        user=config["mysqldb"]["user"],
        port=config["mysqldb"]["port"],
        password=config["mysqldb"]["password"],
        db=config["mysqldb"]["db"],
        charset="utf8",
        cursorclass=pymysql.cursors.DictCursor
    )
    context.redis_client = redis.Redis(
        config["redis"]["host"], port=6379, db=0
    )
    context.cleanup = {
        "sql": "DELETE FROM accounts WHERE phone='{}' OR email='{}'".format(
            config["account"]["phone"], config["account"]["email"]
        )
    }
    account_url = "http://{}".format(config["api"]["account"])
    res = requests.post("{}/signup".format(account_url), json={
        "phone": config["account"]["phone"],
        "email": config["account"]["email"],
        "password": config["account"]["password"],
        "firstName": config["account"]["firstName"],
        "lastName": config["account"]["lastName"],
        "birthday": config["account"]["birthday"],
        "gender": config["account"]["gender"],
    })
    print(res.text)
    res = requests.post("{}/signin".format(account_url), json={
        "username": "13145678901",
        "password": "12345678",
    })
    print(res.text)
    obj = json.loads(res.text)
    context.token = obj["token"]
    res = requests.get("{}/getaccount".format(account_url), params={
        "token": context.token,
    })
    print(res.text)
    obj = json.loads(res.text)
    context.id = obj["account"]["id"]
    print(context.token, context.id)


def after_all(context):
    if "sql" in context.cleanup:
        with context.mysql_conn.cursor() as cursor:
            cursor.execute(context.cleanup["sql"])
        context.mysql_conn.commit()
    stop()
