# GroundSeg modules
from log import Log
from setup import Setup
from login import Login
from system_get import SysGet
from system_post import SysPost

# Docker
from wireguard import Wireguard
from urbit import Urbit
from webui import WebUI

class Orchestrator:

    wireguard = None

    def __init__(self, config):
        self.config_object = config
        self.config = config.config

        self.wireguard = Wireguard(config)
        self.urbit = Urbit(config)
        #self.minio
        self.webui = WebUI(config)

        self.config_object.gs_ready = True
        Log.log("GroundSeg: Initialization completed")


    #
    #   Setup
    #


    def handle_setup(self, page, data):
        try:
            if page == "anchor":
                return Setup.handle_anchor(data,self.config_object)

            if page == "password":
                if self.config_object.create_password(data['password']):
                    return 200

        except Exception as e:
            Log.log(f"Setup: {e}")

        return 401


    #
    #   Login
    #


    def handle_login_request(self, data):
        res = Login.handle_login(data, self.config_object)
        if res:
            return Login.make_cookie(self.config_object)
        else:
            return Login.failed()


    #
    #   Urbit Pier
    #


    # List of Urbit Ships in Home Page
    def get_urbits(self):
        return self.urbit.list_ships()

    # Handle POST request relating to Urbit ID
    def urbit_post(self ,urbit_id, data):
        try:
            # Boot new Urbit
            if data['app'] == 'boot-new':
                return self.urbit.create(urbit_id, data.get('data'))

        except Exception as e:
            Log.log(f"Urbit: Post Request failed: {e}")

        return 400

        '''
            # Check if Urbit Pier exists
            urb = self._urbits.get(urbit_id)
            if urb == None:
                return 400

            # Urbit Pier requests
            if data['app'] == 'pier':
                if data['data'] == 'toggle':
                    return self.toggle_pier_power(urb)

                if data['data'] == '+code':
                    return self.get_urbit_code(urbit_id, urb)

                if data['data'] == 's3-update':
                    return self.set_minio_endpoint(urbit_id)

                if data['data'] == 's3-unlink':
                    lens_port = self.get_urbit_loopback_addr(urbit_id)
                    return urb.unlink_minio_endpoint(lens_port)

                if data['data'] == 'schedule-meld':
                    return urb.set_meld_schedule(data['frequency'], data['hour'], data['minute'])

                if data['data'] == 'toggle-meld':
                    x = self.get_urbit_loopback_addr(urb.config['pier_name'])
                    return urb.toggle_meld_status(x)

                if data['data'] == 'do-meld':
                    lens_addr = self.get_urbit_loopback_addr(urbit_id)
                    return urb.send_meld(lens_addr)

                if data['data'] == 'export':
                    return self.export_urbit(urb)

                if data['data'] == 'delete':
                    return self.delete_urbit(urbit_id)

                if data['data'] == 'toggle-autostart':
                    return self.toggle_autostart(urbit_id)

                if data['data'] == 'loom':
                    return urb.set_loom_size(data['size'])

            # Wireguard requests
            if data['app'] == 'wireguard':
                if data['data'] == 'toggle':
                    return self.toggle_pier_network(urb)

            # MinIO requests
            if data['app'] == 'minio':
                pwd = data.get('password')
                if pwd != None:
                    return self.create_minio_admin_account(urbit_id, pwd)

                if data['data'] == 'export':
                    return self.export_minio_bucket(urbit_id)

            return 400

        '''


    #
    #   Anchor Settings
    #


    # Get anchor registration information
    def get_anchor_settings(self):
        lease_end = None
        ongoing = False
        #lease = self.anchor_config['lease']

        #if self.anchor_config['ongoing'] == 1:
        #    ongoing = True

        '''
        if lease != None:
            x = list(map(int,lease.split('-')))
            lease_end = datetime(x[0], x[1], x[2], 0, 0)
            '''

        anchor = {
                "wgReg": self.config['wgRegistered'],
                "wgRunning": False, #TODO:self.wireguard.is_running()
                "lease": lease_end,
                "ongoing": ongoing
                }

        return {'anchor': anchor}


    #
    #   System Settings
    #


    # Get all system information
    def get_system_settings(self):
        is_vm = "vm" == self.config_object.device_mode

        ver = str(self.config_object.version)
        if self.config['updateBranch'] == 'edge':
            settings['gsVersion'] = ver + '-edge'

        required = {
                "vm": is_vm,
                "updateMode": self.config['updateMode'],
                "minio": False, #TODO:  self.minIO_on
                "containers" : SysGet.get_containers(),
                "sessions": len(self.config['sessions']),
                "gsVersion": ver
                }

        optional = {} 
        if not is_vm:
            optional = {
                    "ram": self.config_object._ram,
                    "cpu": self.config_object._cpu,
                    "temp": self.config_object._core_temp,
                    "disk": self.config_object._disk,
                    "connected": SysGet.get_connection_status(),
                    "ethOnly": SysGet.get_ethernet_status()
                    }

        settings = {**optional, **required}
        return {'system': settings}

    # Modify system settings
    def system_post(self, module, data, sessionid):

        # sessions module
        if module == 'session':
            return SysPost.handle_session(data, self.config_object, sessionid)

        # power module
        if module == 'power':
            return SysPost.handle_power(data)

        # binary module
        if module == 'binary':
            return SysPost.handle_binary(data)

        # network connectivity module
        if module == 'network':
            return SysPost.handle_network(data,self.config_object)

        # watchtower module
        if module == 'watchtower':
            return SysPost.handle_updater(data, self.config_object)

        '''
        # logs module
        if module == 'logs':
            if data['action'] == 'view':
                return self.get_log_lines(data['container'], data['haveLine'])

            if data['action'] == 'export':
                return '\n'.join(self.get_log_lines(data['container'], 0))


        # minIO module
        if module == 'minio':
            if data['action'] == 'reload':
                self.toggle_minios_off()
                self.toggle_minios_on()
                time.sleep(1)
                return 200
        

        # anchor module
        if module == 'anchor':
            if data['action'] == 'restart':
                return self.restart_anchor()

            if data['action'] == 'register':
                return self.register_device(data['key']) 

            if data['action'] == 'toggle':
                running = self.wireguard.is_running()
                if running:
                    return self.toggle_anchor_off()
                return self.toggle_anchor_on()

            if data['action'] == 'get-url':
                return self.config['endpointUrl']

            if data['action'] == 'change-url':
                return self.change_wireguard_url(data['url'])

            if data['action'] == 'unsubscribe':
                endpoint = self.config['endpointUrl']
                api_version = self.config['apiVersion']
                url = f'https://{endpoint}/{api_version}'
                x = self.wireguard.cancel_subscription(data['key'],url)
                if x != 400:
                    self.anchor_config = x
                    return 200



        '''
        return module

