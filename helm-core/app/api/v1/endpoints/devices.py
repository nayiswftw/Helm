from typing import List

# pyrefly: ignore [missing-import]
from fastapi import APIRouter, HTTPException

from app.schemas.models import Device
from app.services.devices import device_service

router = APIRouter()


@router.get("/", response_model=List[Device])
def list_devices():
    return device_service.list_devices()


@router.get("/{device_id}", response_model=Device)
def get_device(device_id: str):
    dev = device_service.get_device(device_id)
    if not dev:
        raise HTTPException(status_code=404, detail=f"device not found: {device_id}")
    return dev
