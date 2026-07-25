import socket
import time
import psutil

# Boot time recorded when module is loaded
BOOT_TIME = psutil.boot_time()


class SystemService:
    @staticmethod
    def get_hostname() -> str:
        return socket.gethostname()

    @staticmethod
    def get_uptime() -> float:
        return time.time() - BOOT_TIME

    @staticmethod
    def get_cpu_info() -> dict:
        return {
            "num_cores": psutil.cpu_count(logical=True),
            "usage_percent": psutil.cpu_percent(interval=None),
        }

    @staticmethod
    def get_memory_info() -> dict:
        mem = psutil.virtual_memory()
        return {
            "total": mem.total,
            "used": mem.used,
            "available": mem.available,
            "usage_percent": mem.percent,
        }

    @staticmethod
    def get_disk_info() -> dict:
        path = "C:\\" if psutil.WINDOWS else "/"
        try:
            disk = psutil.disk_usage(path)
            return {
                "mount_point": path,
                "total": disk.total,
                "used": disk.used,
                "free": disk.free,
                "usage_percent": disk.percent,
            }
        except Exception:
            return {
                "mount_point": path,
                "total": 0,
                "used": 0,
                "free": 0,
                "usage_percent": 0.0,
            }

    @staticmethod
    def get_network_info() -> dict:
        net = psutil.net_io_counters()
        return {"rx_bytes": net.bytes_recv, "tx_bytes": net.bytes_sent}
