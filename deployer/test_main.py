import asyncio
import os
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
import pytest_asyncio
from httpx import ASGITransport, AsyncClient

os.environ.setdefault("DEPLOY_TOKEN", "test-token")
os.environ.setdefault("MICROSERVICES_HOST_DIR", "/test/microservices")
os.environ.setdefault("PUTTYKNIFE_HOST_DIR", "/test/puttyknife")

from main import PROJECT_LOCKS, app  # noqa: E402

TOKEN = "test-token"
BAD_TOKEN = "wrong-token"


@pytest_asyncio.fixture
async def client():
    async with AsyncClient(
        transport=ASGITransport(app=app), base_url="http://test"
    ) as c:
        yield c


@pytest_asyncio.fixture(autouse=True)
async def reset_locks():
    """Replace any locked locks with fresh ones between tests."""
    for key in list(PROJECT_LOCKS.keys()):
        PROJECT_LOCKS[key] = asyncio.Lock()
    yield


def make_proc(returncode: int, stdout: bytes = b""):
    proc = MagicMock()
    proc.returncode = returncode
    proc.communicate = AsyncMock(return_value=(stdout, b""))
    return proc


# ── Health ────────────────────────────────────────────────────────────────────


@pytest.mark.asyncio
async def test_health(client):
    resp = await client.get("/health")
    assert resp.status_code == 200
    assert resp.json() == {"status": "ok"}


# ── Auth ──────────────────────────────────────────────────────────────────────


@pytest.mark.asyncio
async def test_deploy_missing_token(client):
    resp = await client.post(
        "/deploy", json={"project": "urlshortener", "commit": "abc"}
    )
    assert resp.status_code == 403


@pytest.mark.asyncio
async def test_deploy_wrong_token(client):
    resp = await client.post(
        "/deploy",
        json={"project": "urlshortener", "commit": "abc"},
        headers={"Authorization": f"Bearer {BAD_TOKEN}"},
    )
    assert resp.status_code == 401


# ── Validation ────────────────────────────────────────────────────────────────


@pytest.mark.asyncio
async def test_deploy_invalid_project(client):
    resp = await client.post(
        "/deploy",
        json={"project": "nonexistent", "commit": "abc"},
        headers={"Authorization": f"Bearer {TOKEN}"},
    )
    assert resp.status_code == 422


# ── Successful deploy ─────────────────────────────────────────────────────────


@pytest.mark.asyncio
@pytest.mark.parametrize("project", ["socialapp", "urlshortener", "puttyknife"])
async def test_deploy_success(client, project):
    proc = make_proc(returncode=0, stdout=b"ok output")
    with patch("main.asyncio.create_subprocess_shell", return_value=proc):
        resp = await client.post(
            "/deploy",
            json={"project": project, "commit": "deadbeef"},
            headers={"Authorization": f"Bearer {TOKEN}"},
        )
    assert resp.status_code == 200
    assert project in resp.text
    assert "deadbeef" in resp.text
    assert "SUCCESS" in resp.text


# ── Failed deploy ─────────────────────────────────────────────────────────────


@pytest.mark.asyncio
async def test_deploy_command_failure(client):
    proc = make_proc(returncode=1, stdout=b"build error")
    with patch("main.asyncio.create_subprocess_shell", return_value=proc):
        resp = await client.post(
            "/deploy",
            json={"project": "urlshortener", "commit": "abc123"},
            headers={"Authorization": f"Bearer {TOKEN}"},
        )
    assert resp.status_code == 500
    assert "FAILED" in resp.text
    assert "build error" in resp.text


@pytest.mark.asyncio
async def test_deploy_exception(client):
    with patch(
        "main.asyncio.create_subprocess_shell", side_effect=OSError("no docker")
    ):
        resp = await client.post(
            "/deploy",
            json={"project": "urlshortener", "commit": "abc123"},
            headers={"Authorization": f"Bearer {TOKEN}"},
        )
    assert resp.status_code == 500
    assert "ERROR" in resp.text
    assert "no docker" in resp.text


# ── Concurrency ───────────────────────────────────────────────────────────────


@pytest.mark.asyncio
async def test_same_project_serialized(client):
    """Two concurrent deploys for the same project must not overlap."""
    order: list[str] = []

    async def slow_proc(*_args, **_kwargs):
        order.append("start")
        proc = MagicMock()
        proc.returncode = 0

        async def communicate():
            await asyncio.sleep(0.05)
            order.append("end")
            return b"done", b""

        proc.communicate = communicate
        return proc

    with patch("main.asyncio.create_subprocess_shell", side_effect=slow_proc):
        await asyncio.gather(
            client.post(
                "/deploy",
                json={"project": "urlshortener", "commit": "a"},
                headers={"Authorization": f"Bearer {TOKEN}"},
            ),
            client.post(
                "/deploy",
                json={"project": "urlshortener", "commit": "b"},
                headers={"Authorization": f"Bearer {TOKEN}"},
            ),
        )

    # Serialized order: first deploy fully completes before second starts
    starts = [i for i, v in enumerate(order) if v == "start"]
    ends = [i for i, v in enumerate(order) if v == "end"]
    assert ends[0] < starts[1], f"deploys overlapped: {order}"


@pytest.mark.asyncio
async def test_different_projects_concurrent(client):
    """Two concurrent deploys for different projects must not deadlock."""
    proc = make_proc(returncode=0, stdout=b"ok")
    with patch("main.asyncio.create_subprocess_shell", return_value=proc):
        results = await asyncio.gather(
            client.post(
                "/deploy",
                json={"project": "socialapp", "commit": "a"},
                headers={"Authorization": f"Bearer {TOKEN}"},
            ),
            client.post(
                "/deploy",
                json={"project": "urlshortener", "commit": "b"},
                headers={"Authorization": f"Bearer {TOKEN}"},
            ),
        )
    assert all(r.status_code == 200 for r in results)
