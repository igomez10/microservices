"""
API Client for Socialapp API
Provides wrapper functions for all endpoints defined in openapi.yaml
"""

import base64
from typing import Any, Dict, List, Optional

import requests


class SocialappAPIClient:
    """Client for interacting with the Socialapp API"""

    def __init__(self, base_url: str = "http://localhost:8080"):
        """Initialize the API client with a base URL"""
        self.base_url = base_url.rstrip("/")
        self.token: Optional[str] = None

    def set_token(self, token: str):
        """Set the OAuth2 bearer token for authenticated requests"""
        self.token = token

    def _get_headers(self, auth_required: bool = True) -> Dict[str, str]:
        """Get headers for API requests"""
        headers = {"Content-Type": "application/json"}
        if auth_required and self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        return headers

    def _handle_response(self, response: requests.Response) -> Dict[str, Any]:
        """Handle API response and return JSON or raise error"""
        try:
            if response.status_code in [200, 201, 204]:
                if response.text:
                    return {
                        "success": True,
                        "data": response.json(),
                        "status_code": response.status_code,
                    }
                return {
                    "success": True,
                    "data": None,
                    "status_code": response.status_code,
                }
            else:
                error_data = (
                    response.json()
                    if response.text
                    else {"message": "Unknown error"}
                )
                return {
                    "success": False,
                    "error": error_data,
                    "status_code": response.status_code,
                }
        except Exception as e:
            return {
                "success": False,
                "error": {"message": str(e)},
                "status_code": response.status_code,
            }

    # ==================== Authentication ====================

    def get_access_token(self, username: str, password: str) -> Dict[str, Any]:
        """Get OAuth2 access token using Basic Auth"""
        credentials = f"{username}:{password}"
        encoded_credentials = base64.b64encode(credentials.encode()).decode()
        headers = {
            "Authorization": f"Basic {encoded_credentials}",
            "Content-Type": "application/json",
        }
        try:
            response = requests.post(
                f"{self.base_url}/v1/oauth/token", headers=headers
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    # ==================== Users ====================

    def create_user(
        self,
        username: str,
        password: str,
        first_name: str,
        last_name: str,
        email: str,
    ) -> Dict[str, Any]:
        """Create a new user (no auth required)"""
        data = {
            "username": username,
            "password": password,
            "first_name": first_name,
            "last_name": last_name,
            "email": email,
        }
        try:
            response = requests.post(
                f"{self.base_url}/v1/users",
                json=data,
                headers={"Content-Type": "application/json"},
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def list_users(self, limit: int = 20, offset: int = 0) -> Dict[str, Any]:
        """List all users with pagination"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/users",
                params={"limit": limit, "offset": offset},
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_user(self, username: str) -> Dict[str, Any]:
        """Get a user by username"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/users/{username}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def update_user(
        self, username: str, user_data: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Update a user"""
        try:
            response = requests.put(
                f"{self.base_url}/v1/users/{username}",
                json=user_data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def delete_user(self, username: str) -> Dict[str, Any]:
        """Delete a user"""
        try:
            response = requests.delete(
                f"{self.base_url}/v1/users/{username}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def change_password(
        self, old_password: str, new_password: str
    ) -> Dict[str, Any]:
        """Change password for current user"""
        data = {"old_password": old_password, "new_password": new_password}
        try:
            response = requests.post(
                f"{self.base_url}/v1/password",
                json=data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def reset_password(self, email: str) -> Dict[str, Any]:
        """Reset password for a user"""
        data = {"email": email}
        try:
            response = requests.put(
                f"{self.base_url}/v1/password",
                json=data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_user_comments(
        self, username: str, limit: int = 20, offset: int = 0
    ) -> Dict[str, Any]:
        """Get all comments for a user"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/users/{username}/comments",
                params={"limit": limit, "offset": offset},
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_user_followers(self, username: str) -> Dict[str, Any]:
        """Get all followers for a user"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/users/{username}/followers",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_user_following(self, username: str) -> Dict[str, Any]:
        """Get all users that a user is following"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/users/{username}/following",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_user_roles(self, username: str) -> Dict[str, Any]:
        """Get all roles for a user"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/users/{username}/roles",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def update_user_roles(
        self, username: str, roles: List[str]
    ) -> Dict[str, Any]:
        """Update all roles for a user"""
        try:
            response = requests.put(
                f"{self.base_url}/v1/users/{username}/roles",
                json=roles,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    # ==================== Following ====================

    def follow_user(
        self, followed_username: str, follower_username: str
    ) -> Dict[str, Any]:
        """Follow a user"""
        try:
            response = requests.post(
                f"{self.base_url}/v1/users/{followed_username}/followers/{follower_username}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def unfollow_user(
        self, followed_username: str, follower_username: str
    ) -> Dict[str, Any]:
        """Unfollow a user"""
        try:
            response = requests.delete(
                f"{self.base_url}/v1/users/{followed_username}/followers/{follower_username}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    # ==================== Comments ====================

    def get_comment(self, comment_id: int) -> Dict[str, Any]:
        """Get a comment by ID"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/comments/{comment_id}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def create_comment(self, content: str, username: str) -> Dict[str, Any]:
        """Create a new comment"""
        data = {"content": content, "username": username}
        try:
            response = requests.post(
                f"{self.base_url}/v1/comments",
                json=data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_user_feed(self) -> Dict[str, Any]:
        """Get the feed for the authenticated user"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/feed", headers=self._get_headers()
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    # ==================== Roles ====================

    def list_roles(self, limit: int = 20, offset: int = 0) -> Dict[str, Any]:
        """List all roles"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/roles",
                params={"limit": limit, "offset": offset},
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def create_role(self, name: str, description: str = "") -> Dict[str, Any]:
        """Create a new role"""
        data = {"name": name, "description": description}
        try:
            response = requests.post(
                f"{self.base_url}/v1/roles",
                json=data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_role(self, role_id: int) -> Dict[str, Any]:
        """Get a role by ID"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/roles/{role_id}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def update_role(
        self, role_id: int, role_data: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Update a role"""
        try:
            response = requests.put(
                f"{self.base_url}/v1/roles/{role_id}",
                json=role_data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def delete_role(self, role_id: int) -> Dict[str, Any]:
        """Delete a role"""
        try:
            response = requests.delete(
                f"{self.base_url}/v1/roles/{role_id}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def list_role_scopes(
        self, role_id: int, limit: int = 20, offset: int = 0
    ) -> Dict[str, Any]:
        """Get all scopes for a role"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/roles/{role_id}/scopes",
                params={"limit": limit, "offset": offset},
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def add_scopes_to_role(
        self, role_id: int, scope_names: List[str]
    ) -> Dict[str, Any]:
        """Add scopes to a role"""
        try:
            response = requests.post(
                f"{self.base_url}/v1/roles/{role_id}/scopes",
                json=scope_names,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def remove_scope_from_role(
        self, role_id: int, scope_id: int
    ) -> Dict[str, Any]:
        """Remove a scope from a role"""
        try:
            response = requests.delete(
                f"{self.base_url}/v1/roles/{role_id}/scopes/{scope_id}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    # ==================== Scopes ====================

    def list_scopes(self, limit: int = 20, offset: int = 0) -> Dict[str, Any]:
        """List all scopes"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/scopes",
                params={"limit": limit, "offset": offset},
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def create_scope(self, name: str, description: str) -> Dict[str, Any]:
        """Create a new scope"""
        data = {"name": name, "description": description}
        try:
            response = requests.post(
                f"{self.base_url}/v1/scopes",
                json=data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_scope(self, scope_id: int) -> Dict[str, Any]:
        """Get a scope by ID"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/scopes/{scope_id}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def update_scope(
        self, scope_id: int, scope_data: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Update a scope"""
        try:
            response = requests.put(
                f"{self.base_url}/v1/scopes/{scope_id}",
                json=scope_data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def delete_scope(self, scope_id: int) -> Dict[str, Any]:
        """Delete a scope"""
        try:
            response = requests.delete(
                f"{self.base_url}/v1/scopes/{scope_id}",
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    # ==================== URLs ====================

    def create_url(self, url: str, alias: str) -> Dict[str, Any]:
        """Create a shortened URL"""
        data = {"url": url, "alias": alias}
        try:
            response = requests.post(
                f"{self.base_url}/v1/urls",
                json=data,
                headers=self._get_headers(),
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_url_data(self, alias: str) -> Dict[str, Any]:
        """Get URL metadata"""
        try:
            response = requests.get(
                f"{self.base_url}/v1/urls/{alias}/data",
                headers={"Content-Type": "application/json"},
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def delete_url(self, alias: str) -> Dict[str, Any]:
        """Delete a URL"""
        try:
            response = requests.delete(
                f"{self.base_url}/v1/urls/{alias}", headers=self._get_headers()
            )
            return self._handle_response(response)
        except Exception as e:
            return {"success": False, "error": {"message": str(e)}}

    def get_url_redirect(self, alias: str) -> str:
        """Get the redirect URL for an alias"""
        return f"{self.base_url}/v1/urls/{alias}"
