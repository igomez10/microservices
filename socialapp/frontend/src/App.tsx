import { useMemo, useState, type Dispatch, type SetStateAction } from 'react'
import AuthPanel from '@/components/AuthPanel'
import { clearStored, readStored, writeStored } from '@/lib/storage'
import {
  createSocialApi,
  type Comment,
  type CreateUserRequest,
  type Role,
  type Scope,
  type URLRecord,
  type User
} from '@/api/social'

const DEFAULT_API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ??
  (typeof window !== 'undefined' ? `${window.location.origin}/api` : 'http://localhost:8086')
const DEFAULT_CLIENT_ID = import.meta.env.VITE_CLIENT_ID ?? ''
const DEFAULT_CLIENT_SECRET = import.meta.env.VITE_CLIENT_SECRET ?? ''
const DEFAULT_SCOPES = import.meta.env.VITE_OAUTH_SCOPES ?? ''

const BASE_URL_KEY = 'socialapp.baseUrl'
const TOKEN_KEY = 'socialapp.token'

type View = 'home' | 'users' | 'admin' | 'urls' | 'settings'

type StatusState = {
  status: string
  error: string
  loading: boolean
}

const initialStatus: StatusState = { status: '', error: '', loading: false }

export default function App() {
  const [view, setView] = useState<View>('home')
  const [baseUrl, setBaseUrl] = useState(() => readStored(BASE_URL_KEY, DEFAULT_API_BASE_URL))
  const [token, setToken] = useState<string | undefined>(() => {
    const stored = readStored<string | undefined>(TOKEN_KEY, undefined)
    return stored && stored.length > 0 ? stored : undefined
  })
  const [clientId, setClientId] = useState(DEFAULT_CLIENT_ID)
  const [clientSecret, setClientSecret] = useState(DEFAULT_CLIENT_SECRET)
  const [scopes, setScopes] = useState(DEFAULT_SCOPES)

  const api = useMemo(
    () =>
      createSocialApi({
        baseUrl,
        token,
        clientId,
        clientSecret
      }),
    [baseUrl, token, clientId, clientSecret]
  )

  const [serverStatus, setServerStatus] = useState<StatusState>(initialStatus)
  const [serverMessage, setServerMessage] = useState('')

  const [feedStatus, setFeedStatus] = useState<StatusState>(initialStatus)
  const [feedItems, setFeedItems] = useState<Comment[]>([])
  const [commentForm, setCommentForm] = useState({ username: '', content: '' })

  const [commentLookupStatus, setCommentLookupStatus] = useState<StatusState>(initialStatus)
  const [commentLookupId, setCommentLookupId] = useState('')
  const [commentLookupResult, setCommentLookupResult] = useState<Comment | null>(null)

  const [usersStatus, setUsersStatus] = useState<StatusState>(initialStatus)
  const [users, setUsers] = useState<User[]>([])
  const [userSearch, setUserSearch] = useState('')
  const [newUser, setNewUser] = useState<CreateUserRequest>({
    username: '',
    first_name: '',
    last_name: '',
    password: '',
    email: ''
  })

  const [profileStatus, setProfileStatus] = useState<StatusState>(initialStatus)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [profileForm, setProfileForm] = useState<User | null>(null)
  const [profileComments, setProfileComments] = useState<Comment[]>([])
  const [followers, setFollowers] = useState<User[]>([])
  const [following, setFollowing] = useState<User[]>([])
  const [rolesForUser, setRolesForUser] = useState<Role[]>([])
  const [rolesInput, setRolesInput] = useState('')
  const [followActor, setFollowActor] = useState('')

  const [rolesStatus, setRolesStatus] = useState<StatusState>(initialStatus)
  const [roles, setRoles] = useState<Role[]>([])
  const [roleForm, setRoleForm] = useState<Role>({ name: '', description: '' })
  const [roleId, setRoleId] = useState('')

  const [scopesStatus, setScopesStatus] = useState<StatusState>(initialStatus)
  const [scopesList, setScopesList] = useState<Scope[]>([])
  const [scopeForm, setScopeForm] = useState<Scope>({ name: '', description: '' })
  const [scopeId, setScopeId] = useState('')

  const [roleScopesStatus, setRoleScopesStatus] = useState<StatusState>(initialStatus)
  const [roleScopes, setRoleScopes] = useState<Scope[]>([])
  const [roleScopeId, setRoleScopeId] = useState('')
  const [roleScopeInput, setRoleScopeInput] = useState('')
  const [removeScopeId, setRemoveScopeId] = useState('')

  const [userRoleStatus, setUserRoleStatus] = useState<StatusState>(initialStatus)
  const [userRoleUsername, setUserRoleUsername] = useState('')

  const [passwordStatus, setPasswordStatus] = useState<StatusState>(initialStatus)
  const [passwordForm, setPasswordForm] = useState({ old_password: '', new_password: '' })
  const [resetForm, setResetForm] = useState({ email: '' })

  const [urlStatus, setUrlStatus] = useState<StatusState>(initialStatus)
  const [urlForm, setUrlForm] = useState<URLRecord>({ alias: '', url: '' })
  const [urlLookupAlias, setUrlLookupAlias] = useState('')
  const [urlData, setUrlData] = useState<URLRecord | null>(null)
  const [urlResolve, setUrlResolve] = useState<string>('')
  const [urlDeleteAlias, setUrlDeleteAlias] = useState('')

  const setStatus = (setter: Dispatch<SetStateAction<StatusState>>, status: StatusState) =>
    setter(status)

  const handleTokenChange = (value: string | undefined) => {
    setToken(value)
    if (value) {
      writeStored(TOKEN_KEY, value)
    } else {
      clearStored(TOKEN_KEY)
    }
  }

  const handleBaseUrlChange = (value: string) => {
    setBaseUrl(value)
    writeStored(BASE_URL_KEY, value)
  }

  const pingServer = async () => {
    setStatus(setServerStatus, { ...initialStatus, loading: true })
    const res = await api.welcome()
    setServerMessage(res.rawText ?? (typeof res.data === 'string' ? res.data : ''))
    setStatus(setServerStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const loadFeed = async () => {
    setStatus(setFeedStatus, { ...initialStatus, loading: true })
    const res = await api.getUserFeed()
    setFeedItems(res.data ?? [])
    setStatus(setFeedStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const postComment = async () => {
    setStatus(setFeedStatus, { ...feedStatus, loading: true })
    const res = await api.createComment({
      username: commentForm.username,
      content: commentForm.content
    })
    setStatus(setFeedStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      setCommentForm({ username: '', content: '' })
      await loadFeed()
    }
  }

  const lookupComment = async () => {
    const id = Number(commentLookupId)
    if (!Number.isFinite(id)) {
      setStatus(setCommentLookupStatus, { status: '', error: 'Comment ID must be a number', loading: false })
      return
    }
    setStatus(setCommentLookupStatus, { ...initialStatus, loading: true })
    const res = await api.getComment(id)
    setCommentLookupResult(res.data ?? null)
    setStatus(setCommentLookupStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const loadUsers = async () => {
    setStatus(setUsersStatus, { ...initialStatus, loading: true })
    const res = await api.listUsers()
    setUsers(res.data ?? [])
    setStatus(setUsersStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const registerUser = async () => {
    setStatus(setUsersStatus, { ...usersStatus, loading: true })
    const res = await api.createUser(newUser)
    setStatus(setUsersStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      const createdUsername = newUser.username
      setNewUser({ username: '', first_name: '', last_name: '', password: '', email: '' })
      await loadUsers()
      if (createdUsername) {
        await loadProfile(createdUsername)
      }
    }
  }

  const loadProfile = async (username: string) => {
    setStatus(setProfileStatus, { ...initialStatus, loading: true })
    const [userRes, commentsRes, followersRes, followingRes, rolesRes] = await Promise.all([
      api.getUserByUsername(username),
      api.getUserComments(username),
      api.getUserFollowers(username),
      api.getFollowingUsers(username),
      api.getRolesForUser(username)
    ])

    const user = userRes.data ?? null
    setSelectedUser(user)
    setProfileForm(user)
    setProfileComments(commentsRes.data ?? [])
    setFollowers(followersRes.data ?? [])
    setFollowing(followingRes.data ?? [])
    setRolesForUser(rolesRes.data ?? [])
    setRolesInput((rolesRes.data ?? []).map((role) => role.name).join(', '))
    setUserRoleUsername(username)

    const error =
      userRes.error || commentsRes.error || followersRes.error || followingRes.error || rolesRes.error || ''

    setStatus(setProfileStatus, {
      status: `${userRes.status} ${userRes.statusText}`,
      error,
      loading: false
    })
  }

  const updateProfile = async () => {
    if (!profileForm || !profileForm.username) {
      return
    }
    setStatus(setProfileStatus, { ...profileStatus, loading: true })
    const res = await api.updateUser(profileForm.username, profileForm)
    setStatus(setProfileStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      await loadProfile(profileForm.username)
    }
  }

  const deleteProfile = async () => {
    if (!selectedUser) {
      return
    }
    setStatus(setProfileStatus, { ...profileStatus, loading: true })
    const res = await api.deleteUser(selectedUser.username)
    setStatus(setProfileStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      setSelectedUser(null)
      setProfileForm(null)
      await loadUsers()
    }
  }

  const followSelectedUser = async (action: 'follow' | 'unfollow') => {
    if (!selectedUser || !followActor) {
      setStatus(setProfileStatus, { status: '', error: 'Provide follower username', loading: false })
      return
    }
    setStatus(setProfileStatus, { ...profileStatus, loading: true })
    const res =
      action === 'follow'
        ? await api.followUser(selectedUser.username, followActor)
        : await api.unfollowUser(selectedUser.username, followActor)
    setStatus(setProfileStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      await loadProfile(selectedUser.username)
    }
  }

  const updateUserRoles = async () => {
    const username = userRoleUsername || selectedUser?.username
    if (!username) {
      setStatus(setUserRoleStatus, { status: '', error: 'Provide username', loading: false })
      return
    }
    const roles = rolesInput
      .split(',')
      .map((role) => role.trim())
      .filter(Boolean)
    setStatus(setUserRoleStatus, { ...initialStatus, loading: true })
    const res = await api.updateRolesForUser(username, roles)
    setStatus(setUserRoleStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok && selectedUser?.username === username) {
      await loadProfile(username)
    }
  }

  const loadRoles = async () => {
    setStatus(setRolesStatus, { ...initialStatus, loading: true })
    const res = await api.listRoles()
    setRoles(res.data ?? [])
    setStatus(setRolesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const createRole = async () => {
    setStatus(setRolesStatus, { ...rolesStatus, loading: true })
    const res = await api.createRole(roleForm)
    setStatus(setRolesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      setRoleForm({ name: '', description: '' })
      await loadRoles()
    }
  }

  const fetchRole = async () => {
    const id = Number(roleId)
    if (!Number.isFinite(id)) {
      setStatus(setRolesStatus, { status: '', error: 'Role ID must be a number', loading: false })
      return
    }
    setStatus(setRolesStatus, { ...rolesStatus, loading: true })
    const res = await api.getRole(id)
    if (res.data) {
      setRoleForm({ name: res.data.name, description: res.data.description })
    }
    setStatus(setRolesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const updateRole = async () => {
    const id = Number(roleId)
    if (!Number.isFinite(id)) {
      setStatus(setRolesStatus, { status: '', error: 'Role ID must be a number', loading: false })
      return
    }
    setStatus(setRolesStatus, { ...rolesStatus, loading: true })
    const res = await api.updateRole(id, roleForm)
    setStatus(setRolesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      await loadRoles()
    }
  }

  const deleteRole = async () => {
    const id = Number(roleId)
    if (!Number.isFinite(id)) {
      setStatus(setRolesStatus, { status: '', error: 'Role ID must be a number', loading: false })
      return
    }
    setStatus(setRolesStatus, { ...rolesStatus, loading: true })
    const res = await api.deleteRole(id)
    setStatus(setRolesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      await loadRoles()
    }
  }

  const loadScopes = async () => {
    setStatus(setScopesStatus, { ...initialStatus, loading: true })
    const res = await api.listScopes()
    setScopesList(res.data ?? [])
    setStatus(setScopesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const createScope = async () => {
    setStatus(setScopesStatus, { ...scopesStatus, loading: true })
    const res = await api.createScope(scopeForm)
    setStatus(setScopesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      setScopeForm({ name: '', description: '' })
      await loadScopes()
    }
  }

  const fetchScope = async () => {
    const id = Number(scopeId)
    if (!Number.isFinite(id)) {
      setStatus(setScopesStatus, { status: '', error: 'Scope ID must be a number', loading: false })
      return
    }
    setStatus(setScopesStatus, { ...scopesStatus, loading: true })
    const res = await api.getScope(id)
    if (res.data) {
      setScopeForm({ name: res.data.name, description: res.data.description })
    }
    setStatus(setScopesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const updateScope = async () => {
    const id = Number(scopeId)
    if (!Number.isFinite(id)) {
      setStatus(setScopesStatus, { status: '', error: 'Scope ID must be a number', loading: false })
      return
    }
    setStatus(setScopesStatus, { ...scopesStatus, loading: true })
    const res = await api.updateScope(id, scopeForm)
    setStatus(setScopesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      await loadScopes()
    }
  }

  const deleteScope = async () => {
    const id = Number(scopeId)
    if (!Number.isFinite(id)) {
      setStatus(setScopesStatus, { status: '', error: 'Scope ID must be a number', loading: false })
      return
    }
    setStatus(setScopesStatus, { ...scopesStatus, loading: true })
    const res = await api.deleteScope(id)
    setStatus(setScopesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      await loadScopes()
    }
  }

  const loadRoleScopes = async () => {
    const id = Number(roleScopeId)
    if (!Number.isFinite(id)) {
      setStatus(setRoleScopesStatus, { status: '', error: 'Role ID must be a number', loading: false })
      return
    }
    setStatus(setRoleScopesStatus, { ...initialStatus, loading: true })
    const res = await api.listScopesForRole(id)
    setRoleScopes(res.data ?? [])
    setStatus(setRoleScopesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const addScopesToRole = async () => {
    const id = Number(roleScopeId)
    if (!Number.isFinite(id)) {
      setStatus(setRoleScopesStatus, { status: '', error: 'Role ID must be a number', loading: false })
      return
    }
    const names = roleScopeInput
      .split(',')
      .map((scope) => scope.trim())
      .filter(Boolean)
    setStatus(setRoleScopesStatus, { ...roleScopesStatus, loading: true })
    const res = await api.addScopeToRole(id, names)
    setStatus(setRoleScopesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      await loadRoleScopes()
    }
  }

  const removeScopeFromRole = async () => {
    const roleIdValue = Number(roleScopeId)
    const scopeIdValue = Number(removeScopeId)
    if (!Number.isFinite(roleIdValue) || !Number.isFinite(scopeIdValue)) {
      setStatus(setRoleScopesStatus, { status: '', error: 'Role and scope ID must be numbers', loading: false })
      return
    }
    setStatus(setRoleScopesStatus, { ...roleScopesStatus, loading: true })
    const res = await api.removeScopeFromRole(roleIdValue, scopeIdValue)
    setStatus(setRoleScopesStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      await loadRoleScopes()
    }
  }

  const fetchUserRoles = async () => {
    if (!userRoleUsername) {
      setStatus(setUserRoleStatus, { status: '', error: 'Provide username', loading: false })
      return
    }
    setStatus(setUserRoleStatus, { ...initialStatus, loading: true })
    const res = await api.getRolesForUser(userRoleUsername)
    if (res.data) {
      setRolesForUser(res.data)
      setRolesInput(res.data.map((role) => role.name).join(', '))
    }
    setStatus(setUserRoleStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const changePassword = async () => {
    setStatus(setPasswordStatus, { ...initialStatus, loading: true })
    const res = await api.changePassword(passwordForm)
    setStatus(setPasswordStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      setPasswordForm({ old_password: '', new_password: '' })
    }
  }

  const resetPassword = async () => {
    setStatus(setPasswordStatus, { ...initialStatus, loading: true })
    const res = await api.resetPassword(resetForm)
    setStatus(setPasswordStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      setResetForm({ email: '' })
    }
  }

  const createUrl = async () => {
    setStatus(setUrlStatus, { ...initialStatus, loading: true })
    const res = await api.createUrl(urlForm)
    setUrlData(res.data ?? null)
    setStatus(setUrlStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const fetchUrlData = async () => {
    if (!urlLookupAlias) {
      setStatus(setUrlStatus, { status: '', error: 'Provide alias', loading: false })
      return
    }
    setStatus(setUrlStatus, { ...initialStatus, loading: true })
    const res = await api.getUrlData(urlLookupAlias)
    setUrlData(res.data ?? null)
    setStatus(setUrlStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const resolveUrl = async () => {
    if (!urlLookupAlias) {
      setStatus(setUrlStatus, { status: '', error: 'Provide alias', loading: false })
      return
    }
    setStatus(setUrlStatus, { ...initialStatus, loading: true })
    const res = await api.getUrl(urlLookupAlias)
    setUrlResolve(res.headers?.location ?? '')
    setStatus(setUrlStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const deleteUrl = async () => {
    if (!urlDeleteAlias) {
      setStatus(setUrlStatus, { status: '', error: 'Provide alias', loading: false })
      return
    }
    setStatus(setUrlStatus, { ...initialStatus, loading: true })
    const res = await api.deleteUrl(urlDeleteAlias)
    setStatus(setUrlStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      setUrlDeleteAlias('')
    }
  }

  const filteredUsers = users.filter((user) =>
    user.username.toLowerCase().includes(userSearch.toLowerCase())
  )

  return (
    <div className="shell">
      <header className="topbar">
        <div>
          <div className="brand">Socialapp</div>
          <div className="muted">A social workspace for your API.</div>
        </div>
        <div className="top-controls">
          <label className="field">
            API Base URL
            <input value={baseUrl} onChange={(event) => handleBaseUrlChange(event.target.value)} />
          </label>
          <div className="pill">{token ? 'Session active' : 'No session'}</div>
        </div>
      </header>

      <nav className="nav">
        <button className={view === 'home' ? 'active' : ''} onClick={() => setView('home')}>
          Home
        </button>
        <button className={view === 'users' ? 'active' : ''} onClick={() => setView('users')}>
          Users
        </button>
        <button className={view === 'admin' ? 'active' : ''} onClick={() => setView('admin')}>
          Admin
        </button>
        <button className={view === 'urls' ? 'active' : ''} onClick={() => setView('urls')}>
          URLs
        </button>
        <button className={view === 'settings' ? 'active' : ''} onClick={() => setView('settings')}>
          Settings
        </button>
      </nav>

      <main className="content">
        {view === 'home' && (
          <section className="view">
            <div className="grid">
              <div className="card">
                <div className="card-header">
                  <h2>Server</h2>
                  <button className="button" onClick={pingServer} disabled={serverStatus.loading}>
                    {serverStatus.loading ? 'Pinging…' : 'Ping server'}
                  </button>
                </div>
                <div className="muted">{serverStatus.status || 'Ping the API root endpoint.'}</div>
                {serverStatus.error ? <div className="error">{serverStatus.error}</div> : null}
                {serverMessage ? <div className="highlight">{serverMessage}</div> : null}
              </div>

              <div className="card">
                <div className="card-header">
                  <h2>Compose</h2>
                  <button className="button" onClick={postComment} disabled={feedStatus.loading}>
                    Post
                  </button>
                </div>
                <label className="field">
                  Username
                  <input
                    value={commentForm.username}
                    onChange={(event) =>
                      setCommentForm((prev) => ({ ...prev, username: event.target.value }))
                    }
                  />
                </label>
                <label className="field">
                  Comment
                  <textarea
                    value={commentForm.content}
                    onChange={(event) =>
                      setCommentForm((prev) => ({ ...prev, content: event.target.value }))
                    }
                  />
                </label>
                {feedStatus.error ? <div className="error">{feedStatus.error}</div> : null}
              </div>

              <div className="card span-two">
                <div className="card-header">
                  <h2>Feed</h2>
                  <button className="button secondary" onClick={loadFeed} disabled={feedStatus.loading}>
                    {feedStatus.loading ? 'Loading…' : 'Refresh feed'}
                  </button>
                </div>
                <div className="muted">{feedStatus.status || 'Latest posts from the network.'}</div>
                {feedStatus.error ? <div className="error">{feedStatus.error}</div> : null}
                <div className="list">
                  {feedItems.map((item) => (
                    <div key={`${item.username}-${item.id}`} className="list-item">
                      <div className="avatar">{item.username.slice(0, 1).toUpperCase()}</div>
                      <div>
                        <div className="list-title">{item.username}</div>
                        <div className="muted">{item.content}</div>
                      </div>
                    </div>
                  ))}
                  {feedItems.length === 0 ? <div className="muted">No feed items yet.</div> : null}
                </div>
              </div>

              <div className="card">
                <div className="card-header">
                  <h2>Find Comment</h2>
                  <button className="button secondary" onClick={lookupComment}>
                    Fetch
                  </button>
                </div>
                <label className="field">
                  Comment ID
                  <input value={commentLookupId} onChange={(event) => setCommentLookupId(event.target.value)} />
                </label>
                <div className="muted">{commentLookupStatus.status}</div>
                {commentLookupStatus.error ? <div className="error">{commentLookupStatus.error}</div> : null}
                {commentLookupResult ? (
                  <div className="highlight">{commentLookupResult.content}</div>
                ) : (
                  <div className="muted">No comment loaded.</div>
                )}
              </div>
            </div>
          </section>
        )}

        {view === 'users' && (
          <section className="view">
            <div className="grid">
              <div className="card span-two">
                <div className="card-header">
                  <h2>Directory</h2>
                  <div className="row">
                    <input
                      placeholder="Search usernames"
                      value={userSearch}
                      onChange={(event) => setUserSearch(event.target.value)}
                    />
                    <button className="button secondary" onClick={loadUsers}>
                      Refresh
                    </button>
                  </div>
                </div>
                <div className="muted">{usersStatus.status || 'Browse registered users.'}</div>
                {usersStatus.error ? <div className="error">{usersStatus.error}</div> : null}
                <div className="list">
                  {filteredUsers.map((user) => (
                    <button
                      key={user.username}
                      className={`list-item selectable ${selectedUser?.username === user.username ? 'active' : ''}`}
                      onClick={() => loadProfile(user.username)}
                    >
                      <div className="avatar">{user.username.slice(0, 1).toUpperCase()}</div>
                      <div>
                        <div className="list-title">{user.username}</div>
                        <div className="muted">{user.email}</div>
                      </div>
                    </button>
                  ))}
                  {filteredUsers.length === 0 ? <div className="muted">No users found.</div> : null}
                </div>
              </div>

              <div className="card">
                <div className="card-header">
                  <h2>Register</h2>
                  <button className="button" onClick={registerUser}>
                    Create user
                  </button>
                </div>
                <label className="field">
                  Username
                  <input
                    value={newUser.username}
                    onChange={(event) => setNewUser((prev) => ({ ...prev, username: event.target.value }))}
                  />
                </label>
                <label className="field">
                  First name
                  <input
                    value={newUser.first_name}
                    onChange={(event) => setNewUser((prev) => ({ ...prev, first_name: event.target.value }))}
                  />
                </label>
                <label className="field">
                  Last name
                  <input
                    value={newUser.last_name}
                    onChange={(event) => setNewUser((prev) => ({ ...prev, last_name: event.target.value }))}
                  />
                </label>
                <label className="field">
                  Email
                  <input
                    value={newUser.email}
                    onChange={(event) => setNewUser((prev) => ({ ...prev, email: event.target.value }))}
                  />
                </label>
                <label className="field">
                  Password
                  <input
                    type="password"
                    value={newUser.password}
                    onChange={(event) => setNewUser((prev) => ({ ...prev, password: event.target.value }))}
                  />
                </label>
              </div>
            </div>

            {selectedUser && profileForm ? (
              <div className="grid">
                <div className="card span-two">
                  <div className="card-header">
                    <h2>Profile</h2>
                    <div className="row">
                      <button className="button secondary" onClick={() => loadProfile(selectedUser.username)}>
                        Refresh
                      </button>
                      <button className="button" onClick={updateProfile}>
                        Update
                      </button>
                      <button className="button danger" onClick={deleteProfile}>
                        Delete
                      </button>
                    </div>
                  </div>
                  <div className="muted">{profileStatus.status}</div>
                  {profileStatus.error ? <div className="error">{profileStatus.error}</div> : null}
                  <div className="grid">
                    <label className="field">
                      Username
                      <input
                        value={profileForm.username}
                        onChange={(event) =>
                          setProfileForm((prev) => (prev ? { ...prev, username: event.target.value } : prev))
                        }
                      />
                    </label>
                    <label className="field">
                      Email
                      <input
                        value={profileForm.email}
                        onChange={(event) =>
                          setProfileForm((prev) => (prev ? { ...prev, email: event.target.value } : prev))
                        }
                      />
                    </label>
                    <label className="field">
                      First name
                      <input
                        value={profileForm.first_name}
                        onChange={(event) =>
                          setProfileForm((prev) => (prev ? { ...prev, first_name: event.target.value } : prev))
                        }
                      />
                    </label>
                    <label className="field">
                      Last name
                      <input
                        value={profileForm.last_name}
                        onChange={(event) =>
                          setProfileForm((prev) => (prev ? { ...prev, last_name: event.target.value } : prev))
                        }
                      />
                    </label>
                  </div>
                </div>

                <div className="card">
                  <div className="card-header">
                    <h2>Followers</h2>
                  </div>
                  <div className="list">
                    {followers.map((user) => (
                      <div key={user.username} className="list-item">
                        <div className="avatar">{user.username.slice(0, 1).toUpperCase()}</div>
                        <div>
                          <div className="list-title">{user.username}</div>
                          <div className="muted">{user.email}</div>
                        </div>
                      </div>
                    ))}
                    {followers.length === 0 ? <div className="muted">No followers.</div> : null}
                  </div>
                </div>

                <div className="card">
                  <div className="card-header">
                    <h2>Following</h2>
                  </div>
                  <div className="list">
                    {following.map((user) => (
                      <div key={user.username} className="list-item">
                        <div className="avatar">{user.username.slice(0, 1).toUpperCase()}</div>
                        <div>
                          <div className="list-title">{user.username}</div>
                          <div className="muted">{user.email}</div>
                        </div>
                      </div>
                    ))}
                    {following.length === 0 ? <div className="muted">Not following anyone.</div> : null}
                  </div>
                </div>

                <div className="card span-two">
                  <div className="card-header">
                    <h2>Follow actions</h2>
                    <div className="row">
                      <button className="button" onClick={() => followSelectedUser('follow')}>
                        Follow
                      </button>
                      <button className="button secondary" onClick={() => followSelectedUser('unfollow')}>
                        Unfollow
                      </button>
                    </div>
                  </div>
                  <label className="field">
                    Follower username
                    <input value={followActor} onChange={(event) => setFollowActor(event.target.value)} />
                  </label>
                </div>

                <div className="card span-two">
                  <div className="card-header">
                    <h2>User comments</h2>
                  </div>
                  <div className="list">
                    {profileComments.map((comment) => (
                      <div key={comment.id ?? comment.content} className="list-item">
                        <div className="avatar">{comment.username.slice(0, 1).toUpperCase()}</div>
                        <div>
                          <div className="list-title">{comment.username}</div>
                          <div className="muted">{comment.content}</div>
                        </div>
                      </div>
                    ))}
                    {profileComments.length === 0 ? <div className="muted">No comments yet.</div> : null}
                  </div>
                </div>

                <div className="card">
                  <div className="card-header">
                    <h2>Roles</h2>
                    <button className="button secondary" onClick={updateUserRoles}>
                      Update roles
                    </button>
                  </div>
                  <label className="field">
                    Role names (comma-separated)
                    <input value={rolesInput} onChange={(event) => setRolesInput(event.target.value)} />
                  </label>
                  <div className="list">
                    {rolesForUser.map((role) => (
                      <div key={role.name} className="list-item">
                        <div className="avatar">R</div>
                        <div>
                          <div className="list-title">{role.name}</div>
                          <div className="muted">{role.description}</div>
                        </div>
                      </div>
                    ))}
                    {rolesForUser.length === 0 ? <div className="muted">No roles assigned.</div> : null}
                  </div>
                </div>
              </div>
            ) : (
              <div className="empty">Select a user to view their profile.</div>
            )}
          </section>
        )}

        {view === 'admin' && (
          <section className="view">
            <div className="grid">
              <div className="card span-two">
                <div className="card-header">
                  <h2>Roles</h2>
                  <div className="row">
                    <button className="button secondary" onClick={loadRoles}>
                      Refresh
                    </button>
                    <button className="button secondary" onClick={fetchRole}>
                      Fetch
                    </button>
                    <button className="button" onClick={createRole}>
                      Create
                    </button>
                    <button className="button secondary" onClick={updateRole}>
                      Update
                    </button>
                    <button className="button danger" onClick={deleteRole}>
                      Delete
                    </button>
                  </div>
                </div>
                <div className="muted">{rolesStatus.status}</div>
                {rolesStatus.error ? <div className="error">{rolesStatus.error}</div> : null}
                <div className="grid">
                  <label className="field">
                    Role ID (for update/delete)
                    <input value={roleId} onChange={(event) => setRoleId(event.target.value)} />
                  </label>
                  <label className="field">
                    Name
                    <input value={roleForm.name} onChange={(event) => setRoleForm((prev) => ({ ...prev, name: event.target.value }))} />
                  </label>
                  <label className="field">
                    Description
                    <input
                      value={roleForm.description ?? ''}
                      onChange={(event) => setRoleForm((prev) => ({ ...prev, description: event.target.value }))}
                    />
                  </label>
                </div>
                <div className="list">
                  {roles.map((role) => (
                    <div key={role.id ?? role.name} className="list-item">
                      <div className="avatar">R</div>
                      <div>
                        <div className="list-title">
                          {role.name} {role.id ? `#${role.id}` : ''}
                        </div>
                        <div className="muted">{role.description}</div>
                      </div>
                    </div>
                  ))}
                  {roles.length === 0 ? <div className="muted">No roles loaded.</div> : null}
                </div>
              </div>

              <div className="card span-two">
                <div className="card-header">
                  <h2>Scopes</h2>
                  <div className="row">
                    <button className="button secondary" onClick={loadScopes}>
                      Refresh
                    </button>
                    <button className="button secondary" onClick={fetchScope}>
                      Fetch
                    </button>
                    <button className="button" onClick={createScope}>
                      Create
                    </button>
                    <button className="button secondary" onClick={updateScope}>
                      Update
                    </button>
                    <button className="button danger" onClick={deleteScope}>
                      Delete
                    </button>
                  </div>
                </div>
                <div className="muted">{scopesStatus.status}</div>
                {scopesStatus.error ? <div className="error">{scopesStatus.error}</div> : null}
                <div className="grid">
                  <label className="field">
                    Scope ID (for update/delete)
                    <input value={scopeId} onChange={(event) => setScopeId(event.target.value)} />
                  </label>
                  <label className="field">
                    Name
                    <input value={scopeForm.name} onChange={(event) => setScopeForm((prev) => ({ ...prev, name: event.target.value }))} />
                  </label>
                  <label className="field">
                    Description
                    <input
                      value={scopeForm.description}
                      onChange={(event) => setScopeForm((prev) => ({ ...prev, description: event.target.value }))}
                    />
                  </label>
                </div>
                <div className="list">
                  {scopesList.map((scope) => (
                    <div key={scope.id ?? scope.name} className="list-item">
                      <div className="avatar">S</div>
                      <div>
                        <div className="list-title">
                          {scope.name} {scope.id ? `#${scope.id}` : ''}
                        </div>
                        <div className="muted">{scope.description}</div>
                      </div>
                    </div>
                  ))}
                  {scopesList.length === 0 ? <div className="muted">No scopes loaded.</div> : null}
                </div>
              </div>

              <div className="card span-two">
                <div className="card-header">
                  <h2>Role scopes</h2>
                  <div className="row">
                    <button className="button secondary" onClick={loadRoleScopes}>
                      Load
                    </button>
                    <button className="button" onClick={addScopesToRole}>
                      Add scopes
                    </button>
                    <button className="button danger" onClick={removeScopeFromRole}>
                      Remove scope
                    </button>
                  </div>
                </div>
                <div className="muted">{roleScopesStatus.status}</div>
                {roleScopesStatus.error ? <div className="error">{roleScopesStatus.error}</div> : null}
                <div className="grid">
                  <label className="field">
                    Role ID
                    <input value={roleScopeId} onChange={(event) => setRoleScopeId(event.target.value)} />
                  </label>
                  <label className="field">
                    Scope names (comma-separated)
                    <input value={roleScopeInput} onChange={(event) => setRoleScopeInput(event.target.value)} />
                  </label>
                  <label className="field">
                    Remove scope ID
                    <input value={removeScopeId} onChange={(event) => setRemoveScopeId(event.target.value)} />
                  </label>
                </div>
                <div className="list">
                  {roleScopes.map((scope) => (
                    <div key={scope.id ?? scope.name} className="list-item">
                      <div className="avatar">S</div>
                      <div>
                        <div className="list-title">
                          {scope.name} {scope.id ? `#${scope.id}` : ''}
                        </div>
                        <div className="muted">{scope.description}</div>
                      </div>
                    </div>
                  ))}
                  {roleScopes.length === 0 ? <div className="muted">No role scopes loaded.</div> : null}
                </div>
              </div>

              <div className="card span-two">
                <div className="card-header">
                  <h2>User role assignment</h2>
                  <div className="row">
                    <button className="button secondary" onClick={fetchUserRoles}>
                      Fetch roles
                    </button>
                    <button className="button" onClick={updateUserRoles}>
                      Update roles
                    </button>
                  </div>
                </div>
                <div className="muted">{userRoleStatus.status}</div>
                {userRoleStatus.error ? <div className="error">{userRoleStatus.error}</div> : null}
                <label className="field">
                  Username
                  <input value={userRoleUsername} onChange={(event) => setUserRoleUsername(event.target.value)} />
                </label>
                <label className="field">
                  Role names (comma-separated)
                  <input value={rolesInput} onChange={(event) => setRolesInput(event.target.value)} />
                </label>
              </div>
            </div>
          </section>
        )}

        {view === 'urls' && (
          <section className="view">
            <div className="grid">
              <div className="card">
                <div className="card-header">
                  <h2>Create URL</h2>
                  <button className="button" onClick={createUrl}>
                    Create
                  </button>
                </div>
                <label className="field">
                  Alias
                  <input value={urlForm.alias} onChange={(event) => setUrlForm((prev) => ({ ...prev, alias: event.target.value }))} />
                </label>
                <label className="field">
                  Target URL
                  <input value={urlForm.url} onChange={(event) => setUrlForm((prev) => ({ ...prev, url: event.target.value }))} />
                </label>
              </div>

              <div className="card">
                <div className="card-header">
                  <h2>Resolve</h2>
                  <button className="button secondary" onClick={resolveUrl}>
                    Resolve
                  </button>
                </div>
                <label className="field">
                  Alias
                  <input value={urlLookupAlias} onChange={(event) => setUrlLookupAlias(event.target.value)} />
                </label>
                <div className="muted">{urlResolve ? `Location: ${urlResolve}` : 'No redirect yet.'}</div>
              </div>

              <div className="card">
                <div className="card-header">
                  <h2>Metadata</h2>
                  <button className="button secondary" onClick={fetchUrlData}>
                    Fetch
                  </button>
                </div>
                <label className="field">
                  Alias
                  <input value={urlLookupAlias} onChange={(event) => setUrlLookupAlias(event.target.value)} />
                </label>
                {urlData ? (
                  <div className="highlight">{urlData.alias} → {urlData.url}</div>
                ) : (
                  <div className="muted">No metadata loaded.</div>
                )}
              </div>

              <div className="card">
                <div className="card-header">
                  <h2>Delete</h2>
                  <button className="button danger" onClick={deleteUrl}>
                    Delete
                  </button>
                </div>
                <label className="field">
                  Alias
                  <input value={urlDeleteAlias} onChange={(event) => setUrlDeleteAlias(event.target.value)} />
                </label>
              </div>

              <div className="card span-two">
                <div className="card-header">
                  <h2>URL status</h2>
                </div>
                <div className="muted">{urlStatus.status}</div>
                {urlStatus.error ? <div className="error">{urlStatus.error}</div> : null}
                {urlData ? (
                  <div className="highlight">Alias {urlData.alias} created for {urlData.url}</div>
                ) : (
                  <div className="muted">No URL action yet.</div>
                )}
              </div>
            </div>
          </section>
        )}

        {view === 'settings' && (
          <section className="view">
            <div className="grid">
              <div className="card span-two">
                <AuthPanel
                  baseUrl={baseUrl}
                  token={token}
                  onToken={handleTokenChange}
                  clientId={clientId}
                  clientSecret={clientSecret}
                  scopes={scopes}
                  onClientId={setClientId}
                  onClientSecret={setClientSecret}
                  onScopes={setScopes}
                />
              </div>

              <div className="card">
                <div className="card-header">
                  <h2>Change password</h2>
                  <button className="button" onClick={changePassword}>
                    Update
                  </button>
                </div>
                <label className="field">
                  Old password
                  <input
                    type="password"
                    value={passwordForm.old_password}
                    onChange={(event) => setPasswordForm((prev) => ({ ...prev, old_password: event.target.value }))}
                  />
                </label>
                <label className="field">
                  New password
                  <input
                    type="password"
                    value={passwordForm.new_password}
                    onChange={(event) => setPasswordForm((prev) => ({ ...prev, new_password: event.target.value }))}
                  />
                </label>
              </div>

              <div className="card">
                <div className="card-header">
                  <h2>Reset password</h2>
                  <button className="button secondary" onClick={resetPassword}>
                    Send reset
                  </button>
                </div>
                <label className="field">
                  Email
                  <input value={resetForm.email} onChange={(event) => setResetForm({ email: event.target.value })} />
                </label>
              </div>

              <div className="card span-two">
                <div className="card-header">
                  <h2>Password status</h2>
                </div>
                <div className="muted">{passwordStatus.status}</div>
                {passwordStatus.error ? <div className="error">{passwordStatus.error}</div> : null}
              </div>
            </div>
          </section>
        )}
      </main>
    </div>
  )
}
