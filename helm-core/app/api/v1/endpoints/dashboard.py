# pyrefly: ignore [missing-import]
from fastapi import APIRouter
from app.schemas.models import Dashboard
from app.services.system import SystemService
from app.services.devices import device_service

router = APIRouter()


@router.get("/dashboard", response_model=Dashboard)
def get_dashboard():
    return Dashboard(
        hostname=SystemService.get_hostname(),
        uptime=SystemService.get_uptime(),
        cpu=SystemService.get_cpu_info(),
        memory=SystemService.get_memory_info(),
        disk=SystemService.get_disk_info(),
        network=SystemService.get_network_info(),
        devices=device_service.count_devices(),
    )
