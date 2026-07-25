from pydantic import BaseModel


class Device(BaseModel):
    id: str
    name: str
    type: str
    status: str


class CPUInfo(BaseModel):
    num_cores: int
    usage_percent: float


class MemoryInfo(BaseModel):
    total: int
    used: int
    available: int
    usage_percent: float


class DiskInfo(BaseModel):
    mount_point: str
    total: int
    used: int
    free: int
    usage_percent: float


class NetworkInfo(BaseModel):
    rx_bytes: int
    tx_bytes: int


class Dashboard(BaseModel):
    hostname: str
    uptime: float
    cpu: CPUInfo
    memory: MemoryInfo
    disk: DiskInfo
    network: NetworkInfo
    devices: int
