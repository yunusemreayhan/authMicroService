#!/usr/bin/python3

import os

snap_dependencies = [
    "sqlc", "go", "go-swagger"
]

apt_dependencies = [
    "python3", "python3-pip", "python3-requests", "python3-pytest"
]

def isSnapInstalled(package):
    try:
        res = os.system(f"sudo snap list | grep {package}")
        if res == 0:
            print(f"{package} is installed res: {res}")
            return True
        else:
            print(f"{package} is not installed res: {res}")
            return False
        return True
    except:
        print(f"{package} is not installed")
        return False

def install_snap_dependency(package):
    print(f"Installing {package}")
    os.system(f"sudo snap install {package} --classic")

def install_apt_dependency(package):
    print(f"Installing {package}")
    os.system(f"sudo apt install -y {package}")

def is_apt_dependency_installed(package):
    try:
        res = os.system(f"sudo apt list --installed | grep {package}")
        if res == 0:
            print(f"{package} is installed res: {res}")
            return True
        else:
            print(f"{package} is not installed res: {res}")
            return False
        return True
    except:
        print(f"{package} is not installed")
        return False

def check_and_install_snap_dependencies():
    for dependency in snap_dependencies:
        if not isSnapInstalled(dependency):
            install_snap_dependency(dependency)
    for dependency in apt_dependencies:
        install_apt_dependency(dependency)
    os.system("sudo apt autoremove -y")

if __name__ == "__main__":
    check_and_install_snap_dependencies()