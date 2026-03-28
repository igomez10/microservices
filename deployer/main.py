import asyncio
import logging
import os
from typing import Literal

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    datefmt="%Y-%m-%dT%H:%M:%S",
)
logger = logging.getLogger(__name__)

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

def _git_commands(commit: str) -> list[str]:
    return [
        "git fetch --all",
        f"git checkout {commit}",
    ]


PROJECT_CONFIGS: dict[str, dict] = {
    "socialapp": {
        "workdir": f"{MICROSERVICES_HOST_DIR}/socialapp",
        "commands": lambda commit: [
            *_git_commands(commit),
            "docker compose build socialapp",
            "docker compose up -d --no-deps socialapp",
            "docker compose build frontend",
            "docker compose up -d --no-deps frontend",
            "docker builder prune -f",
        ],
    },
    "urlshortener": {
        "workdir": f"{MICROSERVICES_HOST_DIR}/urlshortener",
        "commands": lambda commit: [
            *_git_commands(commit),
            "docker compose build urlshortener",
            "docker compose down",
            "docker compose up -d",
            "docker builder prune -f",
        ],
    },
    "puttyknife": {
        "workdir": PUTTYKNIFE_HOST_DIR,
        "commands": lambda commit: [
            *_git_commands(commit),
            "docker compose build puttyknife-m2",
            "docker compose up -d --no-deps puttyknife-m2",
            "docker compose build puttyknife-fr-producer",
            "docker compose up -d --no-deps puttyknife-fr-producer",
            "docker compose build puttyknife-fincaraiz-consumer",
            "docker compose up -d --no-deps puttyknife-fincaraiz-consumer",
            "docker compose build puttyknife-fincaraiz-elastic-producer",
            "docker compose up -d --no-deps puttyknife-fincaraiz-elastic-producer",
            "docker compose build prediction",
            "docker compose up -d --no-deps prediction",
            "docker compose build puttyknife-server",
            "docker compose up -d --no-deps puttyknife-server",
            "docker compose build streamlit-app",
            "docker compose up -d --no-deps streamlit-app",
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
    commands = config["commands"](req.commit)

    env = os.environ.copy()
    env["GIT_SSH_COMMAND"] = "ssh -o StrictHostKeyChecking=no -i /root/.ssh/id_rsa -F /dev/null"

    async with PROJECT_LOCKS[req.project]:
        logger.info("deploy started project=%s commit=%s", req.project, req.commit)
        output: list[str] = []
        output.append(f"=== Deploying {req.project} (commit: {req.commit}) ===\n")

        for cmd in commands:
            output.append(f"\n$ {cmd}\n")
            logger.info("running project=%s cmd=%r", req.project, cmd)
            try:
                proc = await asyncio.create_subprocess_shell(
                    cmd,
                    stdout=asyncio.subprocess.PIPE,
                    stderr=asyncio.subprocess.STDOUT,
                    cwd=workdir,
                    env=env,
                )
                try:
                    stdout, _ = await asyncio.wait_for(proc.communicate(), timeout=300)
                except asyncio.TimeoutError:
                    proc.kill()
                    await proc.communicate()
                    logger.error("command timed out project=%s cmd=%r", req.project, cmd)
                    output.append(f"\n=== TIMEOUT: '{cmd}' exceeded 5 minutes ===\n")
                    return Response(
                        content="".join(output),
                        status_code=500,
                        media_type="text/plain; charset=utf-8",
                    )

                output.append(stdout.decode("utf-8", errors="replace"))

                if proc.returncode != 0:
                    logger.error(
                        "command failed project=%s cmd=%r exit_code=%d",
                        req.project, cmd, proc.returncode,
                    )
                    output.append(
                        f"\n=== FAILED: '{cmd}' exited with code {proc.returncode} ===\n"
                    )
                    return Response(
                        content="".join(output),
                        status_code=500,
                        media_type="text/plain; charset=utf-8",
                    )

                logger.info("command ok project=%s cmd=%r", req.project, cmd)
            except Exception as exc:
                logger.exception("exception project=%s cmd=%r error=%s", req.project, cmd, exc)
                output.append(f"\n=== ERROR running '{cmd}': {exc} ===\n")
                return Response(
                    content="".join(output),
                    status_code=500,
                    media_type="text/plain; charset=utf-8",
                )

        logger.info("deploy succeeded project=%s commit=%s", req.project, req.commit)
        output.append(f"\n=== SUCCESS: {req.project} deployed ===\n")
        return Response(
            content="".join(output),
            status_code=200,
            media_type="text/plain; charset=utf-8",
        )
