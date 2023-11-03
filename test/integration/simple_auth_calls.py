#/usr/bin/python3.10

import requests
import os

testusername = "test"
testpassword = "test"
testemail = "test@gmail.com"
voucher = ""

def prepare_test_env():
    os.system("make compose_test")

def tear_down():
    os.system("make clean_compose_test")

def testregister():
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

def testlogin():
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

def testverify():
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
    
def test():
    prepare_test_env()
    testregister()
    testlogin()
    testverify()
    tear_down()

if __name__ == "__main__":
    test()