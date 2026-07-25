# pyrefly: ignore [missing-import]
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    PROJECT_NAME: str = "Helm Core API"
    PORT: int = 8080
    HOST: str = "localhost"

    class Config:
        env_prefix = "HELM_"


settings = Settings()
