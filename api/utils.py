from datetime import datetime
import sys
import os
import nmcli
import docker

'''
class Utils:
    def remove_urbit_containers():
        client = docker.from_env()

        # Force remove containers
        containers = client.containers.list(all=True)
        for container in containers:
            try:
                if container.image.tags[0] == "tloncorp/urbit:latest":
                    container.remove(force=True)
                if container.image.tags[0] == "tloncorp/vere:latest":
                    container.remove(force=True)
            except:
                pass

        # Check if all have been removed
        containers = client.containers.list(all=True)
        count = 0
        for container in containers:
            try:
                if container.image.tags[0] == "tloncorp/urbit:latest":
                    count = count + 1
                if container.image.tags[0] == "tloncorp/vere:latest":
                    count = count + 1
            except:
                pass
        return count == 0


class Network:

    def list_wifi_ssids():
        return [x.ssid for x in nmcli.device.wifi() if len(x.ssid) > 0]

    def wifi_connect(ssid, pwd):
        try:
            nmcli.device.wifi_connect(ssid, pwd)
            return "success"
        except Exception as e:
            return f"failed: {e}"
'''
