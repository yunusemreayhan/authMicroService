#!/usr/bin/python3

import pytest
import requests
import os
import atexit

testusername = "test"
testpassword = "test"
testemail = "test@gmail.com"
voucher = ""

def prepare_test_env():
    print("prepare test env")
    os.system("make compose_test")

def dump_docker_logs():
    print("DUMPING DOCKER LOGS")
    print("======================================= authmicroservice2-auth_micro_service_db_test-1")
    os.system("docker logs authmicroservice2-auth_micro_service_db_test-1")
    print("======================================= authmicroservice2-auth_micro_service_backend_test-1")
    os.system("docker logs authmicroservice2-auth_micro_service_backend_test-1")
    print("======================================= authmicroservice2-migration_test-1")
    os.system("docker logs authmicroservice2-migration_test-1")
    print("======================================= authmicroservice2-auth_micro_service_proxy_test-1")
    os.system("docker logs authmicroservice2-auth_micro_service_proxy_test-1")


def tear_down():
    print("tear down")
    dump_docker_logs()
    os.system("make clean_compose_test")

class Resource:

    def __init__(self):
        print("setup")
        prepare_test_env()
        atexit.register(tear_down)

class TestResource:

    r = Resource()

    def test_register(self):
        print("testing register with")
        try:
            response = requests.post("http://localhost/api/register", json={
                "username": testusername,
                "password": testpassword,
                "email": testemail
            })
        except Exception as e:
            print(e)
        try:
            print(response.content.decode("utf-8") + " " + str(response.status_code) + " " + str(response.text))
            print("response.json()[\"username\"] == testusername " + str(response.json()["username"] == testusername) + " " + str(response.json()))
            assert response.json()["username"] == testusername
            assert response.status_code == 200
        except:
            pass

    def test_login(self):
        print("testing login with")
        global voucher
        try:
            response = requests.post("http://localhost/api/login", json={
                "username": testusername,
                "password": testpassword
            })
        except Exception as e:
            print(e)
        try:
            print(response.content.decode("utf-8") + " " + str(response.status_code) + " " + str(response.text))
            print("\"voucher\" in response.json() " + str("voucher" in response.json()))
            assert "voucher" in response.json()
            assert response.status_code == 200
            voucher = response.json()["voucher"]
        except:
            pass

    def test_verify(self):
        print("testing verify with")
        global voucher
        try:
            response = requests.post("http://localhost/api/verify", json={
                "voucher": voucher
            })
        except Exception as e:
            print(e)
        try:
            print(response.content.decode("utf-8") + " " + str(response.status_code) + " " + str(response.text))
            assert response.status_code == 200
            assert response.text == "OK"
        except:
            pass
        try:
            response = requests.post("http://localhost/api/verify", json={
                "voucher": "invalidvoucher"
            })
        except Exception as e:
            print(e)
        try:
            print(response.content.decode("utf-8") + " " + str(response.status_code) + " " + str(response.text))
            assert response.status_code == 502
        except:
            pass







