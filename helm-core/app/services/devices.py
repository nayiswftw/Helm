from typing import List, Optional
from app.schemas.models import Device
from app.services.system import SystemService


class DeviceService:
    def __init__(self):
        self._devices = {}
        # Register local server initially
        local_hostname = SystemService.get_hostname()
        self._devices["local"] = Device(
            id="local", name=local_hostname, type="server", status="online"
        )

    def list_devices(self) -> List[Device]:
        return list(self._devices.values())

    def get_device(self, device_id: str) -> Optional[Device]:
        return self._devices.get(device_id)

    def count_devices(self) -> int:
        return len(self._devices)


# Singleton instance for now, normally injected
device_service = DeviceService()
