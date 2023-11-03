#!/usr/bin/python3.11

import os

if not os.path.exists("./tools/migrate"):
    os.system("curl -L -v https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz -o ./tools/migrate.tar.gz")
    os.system("tar -xvf ./tools/migrate.tar.gz -C ./tools/")
    os.system("rm ./tools/migrate.tar.gz")
    os.system("rm ./tools/README.md")
    os.system("rm ./tools/LICENSE")
