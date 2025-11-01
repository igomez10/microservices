"""
Socialapp Admin Tool - Streamlit Application
Full-featured admin interface for the Socialapp API
"""

import streamlit as st
from api_client import SocialappAPIClient

# Page configuration
st.set_page_config(
    page_title="Socialapp Admin",
    page_icon="🔧",
    layout="wide",
    initial_sidebar_state="expanded",
)

# Initialize session state
if "api_client" not in st.session_state:
    st.session_state.api_client = SocialappAPIClient()
if "authenticated" not in st.session_state:
    st.session_state.authenticated = False
if "token" not in st.session_state:
    st.session_state.token = None
if "username" not in st.session_state:
    st.session_state.username = None

# Sidebar configuration
st.sidebar.title("🔧 Socialapp Admin")

# API Base URL Configuration
st.sidebar.subheader("⚙️ Configuration")
base_url = st.sidebar.selectbox(
    "API Base URL",
    [
        "http://localhost:8080",
        "http://localhost:8085",
        "https://socialapp.gomezignacio.com",
    ],
    key="base_url",
)
st.session_state.api_client.base_url = base_url

# Authentication status
st.sidebar.subheader("🔐 Authentication")
if st.session_state.authenticated:
    st.sidebar.success(f"✅ Logged in as: {st.session_state.username}")
    if st.sidebar.button("Logout"):
        st.session_state.authenticated = False
        st.session_state.token = None
        st.session_state.username = None
        st.session_state.api_client.set_token(None)
        st.rerun()
else:
    st.sidebar.warning("⚠️ Not authenticated")

# Navigation
st.sidebar.subheader("📋 Navigation")
section = st.sidebar.radio(
    "Choose Section",
    [
        "🔐 Authentication",
        "👥 Users Management",
        "💬 Comments & Feed",
        "👣 Following Management",
        "🎭 Roles Management",
        "🔑 Scopes Management",
        "🔗 URL Shortener",
    ],
)


# Helper function to display API responses
def display_response(response, success_message="Operation successful"):
    """Display API response with proper formatting"""
    if response.get("success"):
        st.success(success_message)
        if response.get("data"):
            with st.expander("📄 Response Data", expanded=True):
                st.json(response["data"])
    else:
        st.error(
            f"❌ Error: {response.get('error', {}).get('message', 'Unknown error')}"
        )
        if response.get("status_code"):
            st.error(f"Status Code: {response['status_code']}")


# ==================== AUTHENTICATION SECTION ====================
if section == "🔐 Authentication":
    st.title("🔐 Authentication")

    tab1, tab2 = st.tabs(["Create User", "Login"])

    with tab1:
        st.subheader("Create New User")
        st.info("ℹ️ No authentication required for user creation")

        with st.form("create_user_form"):
            col1, col2 = st.columns(2)
            with col1:
                username = st.text_input("Username *", key="create_username")
                first_name = st.text_input(
                    "First Name *", key="create_first_name"
                )
                email = st.text_input("Email *", key="create_email")
            with col2:
                password = st.text_input(
                    "Password *", type="password", key="create_password"
                )
                last_name = st.text_input(
                    "Last Name *", key="create_last_name"
                )

            submitted = st.form_submit_button("Create User", type="primary")

            if submitted:
                if all([username, password, first_name, last_name, email]):
                    response = st.session_state.api_client.create_user(
                        username, password, first_name, last_name, email
                    )
                    display_response(response, "✅ User created successfully!")
                else:
                    st.error("❌ All fields are required!")

    with tab2:
        st.subheader("Login to Get Access Token")
        st.info("ℹ️ Enter your credentials to obtain an OAuth2 bearer token")

        with st.form("login_form"):
            login_username = st.text_input("Username", key="login_username")
            login_password = st.text_input(
                "Password", type="password", key="login_password"
            )

            submitted = st.form_submit_button("Login", type="primary")

            if submitted:
                if login_username and login_password:
                    response = st.session_state.api_client.get_access_token(
                        login_username, login_password
                    )
                    if response.get("success"):
                        token_data = response["data"]
                        st.session_state.token = token_data.get("access_token")
                        st.session_state.username = login_username
                        st.session_state.api_client.set_token(
                            st.session_state.token
                        )
                        st.session_state.authenticated = True
                        st.success("✅ Login successful!")

                        with st.expander("🎟️ Token Details", expanded=True):
                            st.json(token_data)
                        st.rerun()
                    else:
                        display_response(response)
                else:
                    st.error("❌ Username and password are required!")

# ==================== USERS MANAGEMENT SECTION ====================
elif section == "👥 Users Management":
    st.title("👥 Users Management")

    operation = st.selectbox(
        "Choose Operation",
        [
            "List Users",
            "Get User by Username",
            "Update User",
            "Delete User",
            "Change Password",
            "Reset Password",
            "Get User Comments",
            "Get User Followers",
            "Get User Following",
            "Manage User Roles",
        ],
    )

    if operation == "List Users":
        st.subheader("📋 List All Users")
        col1, col2 = st.columns(2)
        with col1:
            limit = st.number_input(
                "Limit", min_value=1, max_value=100, value=20
            )
        with col2:
            offset = st.number_input("Offset", min_value=0, value=0)

        if st.button("Fetch Users", type="primary"):
            response = st.session_state.api_client.list_users(limit, offset)
            display_response(response, "✅ Fetched users successfully")

    elif operation == "Get User by Username":
        st.subheader("🔍 Get User by Username")
        username = st.text_input("Username")

        if st.button("Get User", type="primary"):
            if username:
                response = st.session_state.api_client.get_user(username)
                display_response(
                    response, f"✅ User '{username}' retrieved successfully"
                )
            else:
                st.error("❌ Username is required!")

    elif operation == "Update User":
        st.subheader("✏️ Update User")
        username = st.text_input("Username to Update")

        with st.form("update_user_form"):
            col1, col2 = st.columns(2)
            with col1:
                first_name = st.text_input("First Name")
                last_name = st.text_input("Last Name")
            with col2:
                email = st.text_input("Email")

            submitted = st.form_submit_button("Update User", type="primary")

            if submitted:
                if username:
                    user_data = {}
                    if first_name:
                        user_data["first_name"] = first_name
                    if last_name:
                        user_data["last_name"] = last_name
                    if email:
                        user_data["email"] = email
                    if username:
                        user_data["username"] = username

                    if user_data:
                        response = st.session_state.api_client.update_user(
                            username, user_data
                        )
                        display_response(
                            response,
                            f"✅ User '{username}' updated successfully",
                        )
                    else:
                        st.warning("⚠️ No fields to update")
                else:
                    st.error("❌ Username is required!")

    elif operation == "Delete User":
        st.subheader("🗑️ Delete User")
        st.warning("⚠️ This action cannot be undone!")
        username = st.text_input("Username to Delete")

        if st.button("Delete User", type="primary"):
            if username:
                confirm = st.checkbox(
                    f"I confirm I want to delete user '{username}'"
                )
                if confirm:
                    response = st.session_state.api_client.delete_user(
                        username
                    )
                    display_response(
                        response, f"✅ User '{username}' deleted successfully"
                    )
            else:
                st.error("❌ Username is required!")

    elif operation == "Change Password":
        st.subheader("🔑 Change Password")
        st.info("ℹ️ Changes password for the currently authenticated user")

        with st.form("change_password_form"):
            old_password = st.text_input("Old Password", type="password")
            new_password = st.text_input("New Password", type="password")
            confirm_password = st.text_input(
                "Confirm New Password", type="password"
            )

            submitted = st.form_submit_button(
                "Change Password", type="primary"
            )

            if submitted:
                if old_password and new_password:
                    if new_password == confirm_password:
                        response = st.session_state.api_client.change_password(
                            old_password, new_password
                        )
                        display_response(
                            response, "✅ Password changed successfully"
                        )
                    else:
                        st.error("❌ New passwords do not match!")
                else:
                    st.error("❌ Both passwords are required!")

    elif operation == "Reset Password":
        st.subheader("🔄 Reset Password")
        email = st.text_input("Email")

        if st.button("Reset Password", type="primary"):
            if email:
                response = st.session_state.api_client.reset_password(email)
                display_response(response, "✅ Password reset initiated")
            else:
                st.error("❌ Email is required!")

    elif operation == "Get User Comments":
        st.subheader("💬 Get User Comments")
        username = st.text_input("Username")
        col1, col2 = st.columns(2)
        with col1:
            limit = st.number_input(
                "Limit", min_value=1, max_value=100, value=20
            )
        with col2:
            offset = st.number_input("Offset", min_value=0, value=0)

        if st.button("Get Comments", type="primary"):
            if username:
                response = st.session_state.api_client.get_user_comments(
                    username, limit, offset
                )
                display_response(
                    response,
                    f"✅ Comments for '{username}' retrieved successfully",
                )
            else:
                st.error("❌ Username is required!")

    elif operation == "Get User Followers":
        st.subheader("👥 Get User Followers")
        username = st.text_input("Username")

        if st.button("Get Followers", type="primary"):
            if username:
                response = st.session_state.api_client.get_user_followers(
                    username
                )
                display_response(
                    response,
                    f"✅ Followers for '{username}' retrieved successfully",
                )
            else:
                st.error("❌ Username is required!")

    elif operation == "Get User Following":
        st.subheader("👣 Get Users Being Followed")
        username = st.text_input("Username")

        if st.button("Get Following", type="primary"):
            if username:
                response = st.session_state.api_client.get_user_following(
                    username
                )
                display_response(
                    response,
                    f"✅ Following list for '{username}' retrieved successfully",
                )
            else:
                st.error("❌ Username is required!")

    elif operation == "Manage User Roles":
        st.subheader("🎭 Manage User Roles")

        tab1, tab2 = st.tabs(["Get Roles", "Update Roles"])

        with tab1:
            username = st.text_input("Username", key="get_roles_username")
            if st.button("Get User Roles", type="primary"):
                if username:
                    response = st.session_state.api_client.get_user_roles(
                        username
                    )
                    display_response(
                        response,
                        f"✅ Roles for '{username}' retrieved successfully",
                    )
                else:
                    st.error("❌ Username is required!")

        with tab2:
            username = st.text_input("Username", key="update_roles_username")
            roles_input = st.text_area(
                "Roles (one per line)",
                help="Enter role names, one per line. Example:\nadministrator\nuser\neditor",
            )

            if st.button("Update User Roles", type="primary"):
                if username and roles_input:
                    roles = [
                        role.strip()
                        for role in roles_input.split("\n")
                        if role.strip()
                    ]
                    response = st.session_state.api_client.update_user_roles(
                        username, roles
                    )
                    display_response(
                        response,
                        f"✅ Roles for '{username}' updated successfully",
                    )
                else:
                    st.error("❌ Username and at least one role are required!")

# ==================== COMMENTS & FEED SECTION ====================
elif section == "💬 Comments & Feed":
    st.title("💬 Comments & Feed Management")

    operation = st.selectbox(
        "Choose Operation",
        ["View User Feed", "Get Comment by ID", "Create Comment"],
    )

    if operation == "View User Feed":
        st.subheader("📰 View Your Feed")
        st.info("ℹ️ Displays the feed for the currently authenticated user")

        if st.button("Load Feed", type="primary"):
            response = st.session_state.api_client.get_user_feed()
            display_response(response, "✅ Feed loaded successfully")

    elif operation == "Get Comment by ID":
        st.subheader("🔍 Get Comment by ID")
        comment_id = st.number_input("Comment ID", min_value=1, value=1)

        if st.button("Get Comment", type="primary"):
            response = st.session_state.api_client.get_comment(comment_id)
            display_response(
                response, f"✅ Comment #{comment_id} retrieved successfully"
            )

    elif operation == "Create Comment":
        st.subheader("✍️ Create New Comment")

        with st.form("create_comment_form"):
            username = st.text_input(
                "Username", value=st.session_state.username or ""
            )
            content = st.text_area("Comment Content", height=150)

            submitted = st.form_submit_button("Create Comment", type="primary")

            if submitted:
                if username and content:
                    response = st.session_state.api_client.create_comment(
                        content, username
                    )
                    display_response(
                        response, "✅ Comment created successfully"
                    )
                else:
                    st.error("❌ Username and content are required!")

# ==================== FOLLOWING MANAGEMENT SECTION ====================
elif section == "👣 Following Management":
    st.title("👣 Following Management")

    operation = st.selectbox(
        "Choose Operation", ["Follow User", "Unfollow User"]
    )

    if operation == "Follow User":
        st.subheader("➕ Follow User")
        col1, col2 = st.columns(2)
        with col1:
            followed_username = st.text_input("User to Follow")
        with col2:
            follower_username = st.text_input(
                "Follower Username", value=st.session_state.username or ""
            )

        if st.button("Follow", type="primary"):
            if followed_username and follower_username:
                response = st.session_state.api_client.follow_user(
                    followed_username, follower_username
                )
                display_response(
                    response,
                    f"✅ '{follower_username}' now follows '{followed_username}'",
                )
            else:
                st.error("❌ Both usernames are required!")

    elif operation == "Unfollow User":
        st.subheader("➖ Unfollow User")
        col1, col2 = st.columns(2)
        with col1:
            followed_username = st.text_input("User to Unfollow")
        with col2:
            follower_username = st.text_input(
                "Follower Username", value=st.session_state.username or ""
            )

        if st.button("Unfollow", type="primary"):
            if followed_username and follower_username:
                response = st.session_state.api_client.unfollow_user(
                    followed_username, follower_username
                )
                display_response(
                    response,
                    f"✅ '{follower_username}' unfollowed '{followed_username}'",
                )
            else:
                st.error("❌ Both usernames are required!")

# ==================== ROLES MANAGEMENT SECTION ====================
elif section == "🎭 Roles Management":
    st.title("🎭 Roles Management")

    operation = st.selectbox(
        "Choose Operation",
        [
            "List Roles",
            "Create Role",
            "Get Role by ID",
            "Update Role",
            "Delete Role",
            "Manage Role Scopes",
        ],
    )

    if operation == "List Roles":
        st.subheader("📋 List All Roles")
        col1, col2 = st.columns(2)
        with col1:
            limit = st.number_input(
                "Limit", min_value=1, max_value=100, value=20
            )
        with col2:
            offset = st.number_input("Offset", min_value=0, value=0)

        if st.button("Fetch Roles", type="primary"):
            response = st.session_state.api_client.list_roles(limit, offset)
            display_response(response, "✅ Roles retrieved successfully")

    elif operation == "Create Role":
        st.subheader("➕ Create New Role")

        with st.form("create_role_form"):
            name = st.text_input("Role Name *")
            description = st.text_area("Description")

            submitted = st.form_submit_button("Create Role", type="primary")

            if submitted:
                if name:
                    response = st.session_state.api_client.create_role(
                        name, description
                    )
                    display_response(response, "✅ Role created successfully")
                else:
                    st.error("❌ Role name is required!")

    elif operation == "Get Role by ID":
        st.subheader("🔍 Get Role by ID")
        role_id = st.number_input("Role ID", min_value=1, value=1)

        if st.button("Get Role", type="primary"):
            response = st.session_state.api_client.get_role(role_id)
            display_response(
                response, f"✅ Role #{role_id} retrieved successfully"
            )

    elif operation == "Update Role":
        st.subheader("✏️ Update Role")
        role_id = st.number_input("Role ID", min_value=1, value=1)

        with st.form("update_role_form"):
            name = st.text_input("Role Name")
            description = st.text_area("Description")

            submitted = st.form_submit_button("Update Role", type="primary")

            if submitted:
                role_data = {}
                if name:
                    role_data["name"] = name
                if description:
                    role_data["description"] = description

                if role_data:
                    response = st.session_state.api_client.update_role(
                        role_id, role_data
                    )
                    display_response(
                        response, f"✅ Role #{role_id} updated successfully"
                    )
                else:
                    st.warning("⚠️ No fields to update")

    elif operation == "Delete Role":
        st.subheader("🗑️ Delete Role")
        st.warning("⚠️ This action cannot be undone!")
        role_id = st.number_input("Role ID", min_value=1, value=1)

        confirm = st.checkbox(f"I confirm I want to delete role #{role_id}")
        if st.button("Delete Role", type="primary") and confirm:
            response = st.session_state.api_client.delete_role(role_id)
            display_response(
                response, f"✅ Role #{role_id} deleted successfully"
            )

    elif operation == "Manage Role Scopes":
        st.subheader("🔑 Manage Role Scopes")

        tab1, tab2, tab3 = st.tabs(
            ["List Scopes", "Add Scopes", "Remove Scope"]
        )

        with tab1:
            role_id = st.number_input(
                "Role ID", min_value=1, value=1, key="list_scopes_role_id"
            )
            col1, col2 = st.columns(2)
            with col1:
                limit = st.number_input(
                    "Limit", min_value=1, max_value=100, value=20
                )
            with col2:
                offset = st.number_input("Offset", min_value=0, value=0)

            if st.button("List Role Scopes", type="primary"):
                response = st.session_state.api_client.list_role_scopes(
                    role_id, limit, offset
                )
                display_response(
                    response,
                    f"✅ Scopes for role #{role_id} retrieved successfully",
                )

        with tab2:
            role_id = st.number_input(
                "Role ID", min_value=1, value=1, key="add_scopes_role_id"
            )
            scopes_input = st.text_area(
                "Scope Names (one per line)",
                help="Example:\nsocialapp.roles.read\nsocialapp.roles.update",
            )

            if st.button("Add Scopes", type="primary"):
                if scopes_input:
                    scopes = [
                        scope.strip()
                        for scope in scopes_input.split("\n")
                        if scope.strip()
                    ]
                    response = st.session_state.api_client.add_scopes_to_role(
                        role_id, scopes
                    )
                    display_response(
                        response, f"✅ Scopes added to role #{role_id}"
                    )
                else:
                    st.error("❌ At least one scope is required!")

        with tab3:
            col1, col2 = st.columns(2)
            with col1:
                role_id = st.number_input(
                    "Role ID", min_value=1, value=1, key="remove_scope_role_id"
                )
            with col2:
                scope_id = st.number_input("Scope ID", min_value=1, value=1)

            if st.button("Remove Scope", type="primary"):
                response = st.session_state.api_client.remove_scope_from_role(
                    role_id, scope_id
                )
                display_response(
                    response,
                    f"✅ Scope #{scope_id} removed from role #{role_id}",
                )

# ==================== SCOPES MANAGEMENT SECTION ====================
elif section == "🔑 Scopes Management":
    st.title("🔑 Scopes Management")

    operation = st.selectbox(
        "Choose Operation",
        [
            "List Scopes",
            "Create Scope",
            "Get Scope by ID",
            "Update Scope",
            "Delete Scope",
        ],
    )

    if operation == "List Scopes":
        st.subheader("📋 List All Scopes")
        col1, col2 = st.columns(2)
        with col1:
            limit = st.number_input(
                "Limit", min_value=1, max_value=100, value=20
            )
        with col2:
            offset = st.number_input("Offset", min_value=0, value=0)

        if st.button("Fetch Scopes", type="primary"):
            response = st.session_state.api_client.list_scopes(limit, offset)
            display_response(response, "✅ Scopes retrieved successfully")

    elif operation == "Create Scope":
        st.subheader("➕ Create New Scope")

        with st.form("create_scope_form"):
            name = st.text_input(
                "Scope Name *", help="Example: socialapp.users.read"
            )
            description = st.text_area("Description *")

            submitted = st.form_submit_button("Create Scope", type="primary")

            if submitted:
                if name and description:
                    response = st.session_state.api_client.create_scope(
                        name, description
                    )
                    display_response(response, "✅ Scope created successfully")
                else:
                    st.error("❌ Both name and description are required!")

    elif operation == "Get Scope by ID":
        st.subheader("🔍 Get Scope by ID")
        scope_id = st.number_input("Scope ID", min_value=1, value=1)

        if st.button("Get Scope", type="primary"):
            response = st.session_state.api_client.get_scope(scope_id)
            display_response(
                response, f"✅ Scope #{scope_id} retrieved successfully"
            )

    elif operation == "Update Scope":
        st.subheader("✏️ Update Scope")
        scope_id = st.number_input("Scope ID", min_value=1, value=1)

        with st.form("update_scope_form"):
            name = st.text_input("Scope Name")
            description = st.text_area("Description")

            submitted = st.form_submit_button("Update Scope", type="primary")

            if submitted:
                scope_data = {}
                if name:
                    scope_data["name"] = name
                if description:
                    scope_data["description"] = description

                if scope_data:
                    response = st.session_state.api_client.update_scope(
                        scope_id, scope_data
                    )
                    display_response(
                        response, f"✅ Scope #{scope_id} updated successfully"
                    )
                else:
                    st.warning("⚠️ No fields to update")

    elif operation == "Delete Scope":
        st.subheader("🗑️ Delete Scope")
        st.warning("⚠️ This action cannot be undone!")
        scope_id = st.number_input("Scope ID", min_value=1, value=1)

        confirm = st.checkbox(f"I confirm I want to delete scope #{scope_id}")
        if st.button("Delete Scope", type="primary") and confirm:
            response = st.session_state.api_client.delete_scope(scope_id)
            display_response(
                response, f"✅ Scope #{scope_id} deleted successfully"
            )

# ==================== URL SHORTENER SECTION ====================
elif section == "🔗 URL Shortener":
    st.title("🔗 URL Shortener Management")

    operation = st.selectbox(
        "Choose Operation",
        [
            "Create Shortened URL",
            "Get URL Data",
            "Delete URL",
            "Test Redirect",
        ],
    )

    if operation == "Create Shortened URL":
        st.subheader("➕ Create Shortened URL")

        with st.form("create_url_form"):
            url = st.text_input(
                "Full URL *",
                help="Example: https://example.com/very/long/path",
            )
            alias = st.text_input("Alias *", help="Example: mylink")

            submitted = st.form_submit_button(
                "Create Short URL", type="primary"
            )

            if submitted:
                if url and alias:
                    response = st.session_state.api_client.create_url(
                        url, alias
                    )
                    if response.get("success"):
                        display_response(
                            response, "✅ Short URL created successfully"
                        )
                        short_url = f"{st.session_state.api_client.base_url}/v1/urls/{alias}"
                        st.info(f"🔗 Short URL: {short_url}")
                    else:
                        display_response(response)
                else:
                    st.error("❌ Both URL and alias are required!")

    elif operation == "Get URL Data":
        st.subheader("🔍 Get URL Metadata")
        alias = st.text_input("Alias")

        if st.button("Get URL Data", type="primary"):
            if alias:
                response = st.session_state.api_client.get_url_data(alias)
                display_response(
                    response,
                    f"✅ URL data for '{alias}' retrieved successfully",
                )
            else:
                st.error("❌ Alias is required!")

    elif operation == "Delete URL":
        st.subheader("🗑️ Delete URL")
        st.warning("⚠️ This action cannot be undone!")
        alias = st.text_input("Alias to Delete")

        confirm = st.checkbox(f"I confirm I want to delete alias '{alias}'")
        if st.button("Delete URL", type="primary") and confirm:
            if alias:
                response = st.session_state.api_client.delete_url(alias)
                display_response(
                    response, f"✅ URL '{alias}' deleted successfully"
                )
            else:
                st.error("❌ Alias is required!")

    elif operation == "Test Redirect":
        st.subheader("🔗 Test URL Redirect")
        alias = st.text_input("Alias")

        if alias:
            redirect_url = st.session_state.api_client.get_url_redirect(alias)
            st.info(f"Redirect URL: {redirect_url}")
            st.markdown(f"[Click here to test redirect]({redirect_url})")

# Footer
st.sidebar.markdown("---")
st.sidebar.markdown("**Socialapp Admin v1.0**")
st.sidebar.markdown("Built with Streamlit")
