import asyncio
import os
from typing import Literal

from fastapi import Depends, FastAPI, HTTPException
from fastapi.responses import Response
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from pydantic import BaseModel

app = FastAPI(title="Deployer", version="1.0.0")
security = HTTPBearer()

DEPLOY_TOKEN = os.environ["DEPLOY_TOKEN"]

# Host-side base paths, passed in as env vars from docker-compose.
# The directories are mounted at the same path inside the container so that
# paths are consistent between the Docker CLI (in the container) and the
# Docker daemon (on the host), which is required for build contexts and volumes.
MICROSERVICES_HOST_DIR = os.environ.get("MICROSERVICES_HOST_DIR", "/home/ubuntu/microservices")
PUTTYKNIFE_HOST_DIR = os.environ.get("PUTTYKNIFE_HOST_DIR", "/home/ubuntu/puttyknife")

PROJECT_CONFIGS: dict[str, dict] = {
    "socialapp": {
        "workdir": f"{MICROSERVICES_HOST_DIR}/socialapp",
        "commands": [
            "git pull",
            "docker compose build",
            "docker compose up -d --remove-orphans",
            "docker builder prune -f",
        ],
    },
    "urlshortener": {
        "workdir": f"{MICROSERVICES_HOST_DIR}/urlshortener",
        "commands": [
            "git pull",
            "docker compose build",
            "docker compose down",
            "docker compose up -d",
            "docker builder prune -f",
        ],
    },
    "puttyknife": {
        "workdir": PUTTYKNIFE_HOST_DIR,
        "commands": [
            "git pull",
            "docker compose build",
            "docker compose up -d --remove-orphans",
            "docker builder prune -f",
        ],
    },
}

PROJECT_LOCKS: dict[str, asyncio.Lock] = {
    "socialapp": asyncio.Lock(),
    "urlshortener": asyncio.Lock(),
    "puttyknife": asyncio.Lock(),
}


class DeployRequest(BaseModel):
    project: Literal["socialapp", "urlshortener", "puttyknife"]
    commit: str


def verify_token(credentials: HTTPAuthorizationCredentials = Depends(security)):
    if credentials.credentials != DEPLOY_TOKEN:
        raise HTTPException(status_code=401, detail="Invalid token")
    return credentials


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/deploy")
async def deploy(
    req: DeployRequest,
    _: HTTPAuthorizationCredentials = Depends(verify_token),
):
    config = PROJECT_CONFIGS[req.project]
    workdir = config["workdir"]
    commands = config["commands"]

    env = os.environ.copy()
    env["GIT_SSH_COMMAND"] = "ssh -o StrictHostKeyChecking=no -i /root/.ssh/id_rsa"

    async with PROJECT_LOCKS[req.project]:
        output: list[str] = []
        output.append(f"=== Deploying {req.project} (commit: {req.commit}) ===\n")

        for cmd in commands:
            output.append(f"\n$ {cmd}\n")
            try:
                proc = await asyncio.create_subprocess_shell(
                    cmd,
                    stdout=asyncio.subprocess.PIPE,
                    stderr=asyncio.subprocess.STDOUT,
                    cwd=workdir,
                    env=env,
                )
                stdout, _ = await proc.communicate()
                output.append(stdout.decode("utf-8", errors="replace"))

                if proc.returncode != 0:
                    output.append(
                        f"\n=== FAILED: '{cmd}' exited with code {proc.returncode} ===\n"
                    )
                    return Response(
                        content="".join(output),
                        status_code=500,
                        media_type="text/plain; charset=utf-8",
                    )
            except Exception as exc:
                output.append(f"\n=== ERROR running '{cmd}': {exc} ===\n")
                return Response(
                    content="".join(output),
                    status_code=500,
                    media_type="text/plain; charset=utf-8",
                )

        output.append(f"\n=== SUCCESS: {req.project} deployed ===\n")
        return Response(
            content="".join(output),
            status_code=200,
            media_type="text/plain; charset=utf-8",
        )
