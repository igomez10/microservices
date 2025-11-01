# Socialapp Admin Tool

A comprehensive Streamlit-based admin interface for managing the Socialapp API.

## Features

This admin tool provides a complete interface for all Socialapp API endpoints:

- **🔐 Authentication**
  - Create new users without authentication
  - Login with username/password to obtain OAuth2 bearer token
  - Token management and session handling

- **👥 Users Management**
  - List, create, update, and delete users
  - Change and reset passwords
  - View user comments, followers, and following lists
  - Manage user roles

- **💬 Comments & Feed**
  - View user feed
  - Get specific comments by ID
  - Create new comments

- **👣 Following Management**
  - Follow and unfollow users
  - View follower/following relationships

- **🎭 Roles Management**
  - CRUD operations for roles
  - Manage role scopes
  - List, add, and remove scopes from roles

- **🔑 Scopes Management**
  - CRUD operations for scopes
  - View scope details

- **🔗 URL Shortener**
  - Create shortened URLs
  - View URL metadata
  - Delete URLs
  - Test redirects

## Installation

1. Navigate to the streamlit directory:
```bash
cd /Users/ignacio/microservices/socialapp/streamlit
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

## Usage

1. Start the Streamlit app:
```bash
streamlit run app.py
```

2. Open your browser to the URL displayed (typically `http://localhost:8501`)

3. Configure the API base URL in the sidebar:
   - `http://localhost:8080` (default)
   - `http://localhost:8085`
   - `https://socialapp.gomezignacio.com`

4. Create a user or login with existing credentials to get started

## Authentication Flow

1. **First-time Setup**: Use the "Create User" section to create an initial user account
2. **Login**: Use the "Login" tab to authenticate and obtain a bearer token
3. **Authenticated Operations**: Once logged in, the token will be used automatically for all API calls that require authentication

## API Configuration

The app supports multiple API endpoints:
- **Local Development**: `http://localhost:8080` or `http://localhost:8085`
- **Production**: `https://socialapp.gomezignacio.com`

Select your preferred endpoint from the sidebar dropdown.

## Architecture

### Files

- `app.py` - Main Streamlit application with UI components
- `api_client.py` - API client wrapper for all Socialapp endpoints
- `requirements.txt` - Python dependencies

### Session State

The app uses Streamlit session state to manage:
- API client instance
- Authentication token
- Authentication status
- Current username

## Development

The application is structured with:
- **API Client Layer** (`api_client.py`): Handles all HTTP communication with the API
- **UI Layer** (`app.py`): Streamlit interface with section-based navigation
- **Session Management**: Token and authentication state persistence

## API Reference

This tool implements all endpoints defined in the Socialapp OpenAPI specification (`openapi.yaml`), including:
- User management endpoints
- Comment and feed endpoints
- Following/follower relationships
- Role-based access control (RBAC)
- Scope management
- URL shortening service

## Error Handling

The app provides comprehensive error handling:
- API errors are displayed with status codes and error messages
- Success messages confirm operations
- Response data is displayed in expandable JSON viewers

## Notes

- Some operations require authentication (bearer token)
- User creation does not require authentication
- All paginated endpoints default to 20 items per page
- Delete operations require confirmation checkboxes

## Support

For issues or questions about the Socialapp API, refer to the OpenAPI specification at:
`/Users/ignacio/microservices/socialapp/openapi.yaml`

