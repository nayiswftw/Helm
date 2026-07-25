# pyrefly: ignore [missing-import]
from fastapi import APIRouter
from app.api.v1.endpoints import health, dashboard, devices

api_router = APIRouter()

api_router.include_router(health.router, tags=["health"])
api_router.include_router(dashboard.router, tags=["dashboard"])
api_router.include_router(devices.router, prefix="/devices", tags=["devices"])
