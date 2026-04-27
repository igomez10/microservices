import { useEffect, useMemo, useRef, useState, type Dispatch, type SetStateAction } from 'react'
import AuthPanel from '@/components/AuthPanel'
import { apiRequest } from '@/api/client'
import { clearCookie, readCookie, readJwtPayload, writeCookie } from '@/lib/storage'
import {
  createSocialApi,
  type AccessToken,
  type Comment,
  type CreateUserRequest,
  type Role,
  type Scope,
  type URLRecord,
  type User
} from '@/api/social'

const DEFAULT_API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ??
  (typeof window !== 'undefined' ? '/api' : 'http://127.0.0.1:5173/api')
const DEFAULT_CLIENT_ID = import.meta.env.VITE_CLIENT_ID ?? ''
const DEFAULT_CLIENT_SECRET = import.meta.env.VITE_CLIENT_SECRET ?? ''
const ALL_OAUTH_SCOPES = [
  'shortly.url.create',
  'shortly.url.delete',
  'shortly.url.update',
  'socialapp.comments.create',
  'socialapp.comments.delete',
  'socialapp.comments.list',
  'socialapp.comments.read',
  'socialapp.comments.update',
  'socialapp.feed.read',
  'socialapp.follower.create',
  'socialapp.follower.delete',
  'socialapp.follower.read',
  'socialapp.followers.list',
  'socialapp.following.list',
  'socialapp.roles.create',
  'socialapp.roles.delete',
  'socialapp.roles.list',
  'socialapp.roles.read',
  'socialapp.roles.scopes.create',
  'socialapp.roles.scopes.delete',
  'socialapp.roles.scopes.list',
  'socialapp.roles.update',
  'socialapp.scopes.create',
  'socialapp.scopes.delete',
  'socialapp.scopes.list',
  'socialapp.scopes.read',
  'socialapp.scopes.update',
  'socialapp.users.create',
  'socialapp.users.delete',
  'socialapp.users.list',
  'socialapp.users.read',
  'socialapp.users.roles.create',
  'socialapp.users.roles.delete',
  'socialapp.users.roles.list',
  'socialapp.users.roles.update',
  'socialapp.users.update'
].join(' ')
const DEFAULT_SCOPES = import.meta.env.VITE_OAUTH_SCOPES ?? ALL_OAUTH_SCOPES

const TOKEN_COOKIE = 'socialapp.token'
const USERNAME_COOKIE = 'socialapp.username'
const USERS_PER_PAGE = 20
const USERS_PAGE_COUNT = 10

function usernameFromToken(token?: string) {
  if (!token) {
    return ''
  }

  const payload = readJwtPayload(token)
  const candidates = [
    payload?.preferred_username,
    payload?.username,
    payload?.user_name,
    payload?.sub
  ]

  const username = candidates.find((value) => typeof value === 'string' && value.trim().length > 0)
  return typeof username === 'string' ? username : ''
}

type View = 'feed' | 'discovery' | 'users' | 'profile' | 'admin' | 'settings'
type AdminTab = 'roles' | 'scopes' | 'urls'
type SettingsTab = 'auth' | 'password' | 'danger'
type ProfileTab = 'posts' | 'followers' | 'following' | 'roles'

function viewFromPath(pathname: string): View {
  switch (pathname) {
    case '/discovery':
      return 'discovery'
    case '/users':
      return 'users'
    case '/profile':
      return 'profile'
    case '/settings':
      return 'settings'
    case '/admin':
      return 'admin'
    case '/':
    case '/feed':
    default:
      return 'feed'
  }
}

function pathFromView(view: View) {
  switch (view) {
    case 'discovery':
      return '/discovery'
    case 'users':
      return '/users'
    case 'profile':
      return '/profile'
    case 'settings':
      return '/settings'
    case 'admin':
      return '/admin'
    case 'feed':
    default:
      return '/feed'
  }
}

type StatusState = {
  status: string
  error: string
  loading: boolean
}

const initialStatus: StatusState = { status: '', error: '', loading: false }

const AVATAR_COLORS = [
  '#1a1a2e',
  '#16213e',
  '#0f3460',
  '#533483',
  '#4a0e4e',
  '#2c3333',
  '#395b64',
  '#2c3639'
]

function avatarColor(username: string) {
  let hash = 0
  for (let index = 0; index < username.length; index += 1) {
    hash = username.charCodeAt(index) + ((hash << 5) - hash)
  }
  return AVATAR_COLORS[Math.abs(hash) % AVATAR_COLORS.length]
}

function timeAgo(dateStr?: string) {
  if (!dateStr) {
    return 'just now'
  }
  const now = new Date()
  const then = new Date(dateStr)
  if (Number.isNaN(then.getTime())) {
    return 'just now'
  }
  const diff = Math.floor((now.getTime() - then.getTime()) / 1000)
  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  if (diff < 604800) return `${Math.floor(diff / 86400)}d`
  if (then.getFullYear() !== now.getFullYear()) {
    return then.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  }
  return then.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function formatDate(dateStr?: string) {
  if (!dateStr) {
    return 'Unknown'
  }
  const value = new Date(dateStr)
  if (Number.isNaN(value.getTime())) {
    return 'Unknown'
  }
  return value.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
}

function getInitials(user: { username: string; first_name?: string; last_name?: string }) {
  const first = user.first_name?.trim().charAt(0) ?? user.username.charAt(0)
  const last = user.last_name?.trim().charAt(0) ?? ''
  return `${first}${last}`.toUpperCase()
}

type IconProps = {
  d: string
  size?: number
  stroke?: string
  fill?: string
}

function Icon({ d, size = 20, stroke = 'currentColor', fill = 'none' }: IconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill={fill}
      stroke={stroke}
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d={d} />
    </svg>
  )
}

const Icons = {
  home: (size = 20) => <Icon size={size} d="M3 9.5L12 3l9 6.5V20a1 1 0 01-1 1H4a1 1 0 01-1-1V9.5z" />,
  search: (size = 20) => <Icon size={size} d="M21 21l-4.35-4.35M11 19a8 8 0 100-16 8 8 0 000 16z" />,
  user: (size = 20) => <Icon size={size} d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2M12 11a4 4 0 100-8 4 4 0 000 8z" />,
  users: (size = 20) =>
    <Icon
      size={size}
      d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75M9 11a4 4 0 100-8 4 4 0 000 8z"
    />,
  shield: (size = 20) => <Icon size={size} d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />,
  heart: (size = 20, liked = false) => (
    <Icon
      size={size}
      d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z"
      fill={liked ? '#c45d5d' : 'none'}
      stroke={liked ? '#c45d5d' : 'currentColor'}
    />
  ),
  logout: (size = 20) => <Icon size={size} d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4M16 17l5-5-5-5M21 12H9" />,
  settings: (size = 20) => (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <circle cx="12" cy="12" r="3" />
      <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 01-2.83 2.83l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z" />
    </svg>
  )
}

function Avatar({
  user,
  size = 'md'
}: {
  user: { username: string; first_name?: string; last_name?: string }
  size?: 'sm' | 'md' | 'lg' | 'xl'
}) {
  return (
    <div className={`avatar avatar-${size}`} style={{ background: avatarColor(user.username) }}>
      {getInitials(user)}
    </div>
  )
}

export default function App() {
  const [view, setView] = useState<View>(() =>
    typeof window !== 'undefined' ? viewFromPath(window.location.pathname) : 'feed'
  )
  const [loggedIn, setLoggedIn] = useState(() => Boolean(readCookie(TOKEN_COOKIE)))
  const [isSignup, setIsSignup] = useState(false)
  const [adminTab, setAdminTab] = useState<AdminTab>('roles')
  const [settingsTab, setSettingsTab] = useState<SettingsTab>('auth')
  const [profileTab, setProfileTab] = useState<ProfileTab>('posts')

  const [baseUrl] = useState(DEFAULT_API_BASE_URL)
  const [token, setToken] = useState<string | undefined>(() => {
    const stored = readCookie(TOKEN_COOKIE)
    return stored && stored.length > 0 ? stored : undefined
  })
  const [currentUsername, setCurrentUsername] = useState(() => {
    const cookieUsername = readCookie(USERNAME_COOKIE)
    if (cookieUsername) {
      return cookieUsername
    }
    const tokenValue = readCookie(TOKEN_COOKIE)
    return usernameFromToken(tokenValue)
  })
  const [profileUsername, setProfileUsername] = useState('')
  const [clientId, setClientId] = useState(DEFAULT_CLIENT_ID)
  const [clientSecret, setClientSecret] = useState(DEFAULT_CLIENT_SECRET)
  const [scopes, setScopes] = useState(DEFAULT_SCOPES)
  const [loginForm, setLoginForm] = useState({ username: '', password: '' })
  const [loginStatus, setLoginStatus] = useState<StatusState>(initialStatus)

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
  const [isPostingComment, setIsPostingComment] = useState(false)
  const [postedCommentNotice, setPostedCommentNotice] = useState<Comment | null>(null)
  const [discoveryStatus, setDiscoveryStatus] = useState<StatusState>(initialStatus)
  const [discoveryItems, setDiscoveryItems] = useState<Comment[]>([])
  const [likedCommentIds, setLikedCommentIds] = useState<Set<number>>(() => new Set())
  const [likeStatus, setLikeStatus] = useState<StatusState>(initialStatus)

  const [commentLookupStatus, setCommentLookupStatus] = useState<StatusState>(initialStatus)
  const [commentLookupId, setCommentLookupId] = useState('')
  const [commentLookupResult, setCommentLookupResult] = useState<Comment | null>(null)

  const [usersStatus, setUsersStatus] = useState<StatusState>(initialStatus)
  const [users, setUsers] = useState<User[]>([])
  const [userSearch, setUserSearch] = useState('')
  const [usersPage, setUsersPage] = useState(1)
  const [usersLoaded, setUsersLoaded] = useState(false)
  const usersRequestIdRef = useRef(0)
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

  const navigateToView = (nextView: View) => {
    setView(nextView)
    if (typeof window === 'undefined') {
      return
    }
    const nextPath = pathFromView(nextView)
    if (window.location.pathname !== nextPath) {
      window.history.pushState({}, '', nextPath)
    }
  }

  const handleTokenChange = (value: string | undefined) => {
    setToken(value)
    if (value) {
      writeCookie(TOKEN_COOKIE, value)
      const inferredUsername = currentUsername || usernameFromToken(value)
      if (inferredUsername) {
        writeCookie(USERNAME_COOKIE, inferredUsername)
        setCurrentUsername(inferredUsername)
        if (!profileUsername) {
          setProfileUsername(inferredUsername)
        }
      }
      setLoggedIn(true)
    } else {
      clearCookie(TOKEN_COOKIE)
      clearCookie(USERNAME_COOKIE)
      setCurrentUsername('')
      setProfileUsername('')
      setLoggedIn(false)
    }
  }

  const loginWithBasicAuth = async (credentials?: { username?: string; password?: string }) => {
    const username = credentials?.username?.trim() ?? loginForm.username.trim()
    const password = credentials?.password?.trim() ?? loginForm.password.trim()

    if (!username || !password) {
      setStatus(setLoginStatus, { status: '', error: 'Enter a username and password', loading: false })
      return
    }

    setStatus(setLoginStatus, { ...initialStatus, loading: true })
    const res = await apiRequest<AccessToken>({
      baseUrl,
      path: '/v1/oauth/token',
      method: 'POST',
      basicAuth: {
        clientId: username,
        clientSecret: password
      },
      queryParams: scopes ? { scope: scopes } : undefined
    })

    setStatus(setLoginStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })

    if (!res.ok || !res.data?.access_token) {
      return
    }

    handleTokenChange(res.data.access_token)
    writeCookie(USERNAME_COOKIE, username)
    setCurrentUsername(username)
    setProfileUsername(username)
    setLoginForm((prev) => ({ ...prev, password: '' }))
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
    if (!currentUsername) {
      setStatus(setFeedStatus, { status: '', error: 'Missing current username', loading: false })
      return
    }
    setIsPostingComment(true)
    setPostedCommentNotice(null)
    const res = await api.createComment({
      username: currentUsername,
      content: commentForm.content
    })
    setStatus(setFeedStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      setCommentForm({ username: '', content: '' })
      setPostedCommentNotice(res.data ?? { username: currentUsername, content: commentForm.content })
      await loadFeed()
    }
    setIsPostingComment(false)
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

  const loadUsers = async (page = usersPage, search = userSearch) => {
    const requestID = usersRequestIdRef.current + 1
    usersRequestIdRef.current = requestID
    setUsersLoaded(false)
    setStatus(setUsersStatus, { ...initialStatus, loading: true })
    const trimmedSearch = search.trim()
    const res = await api.listUsers({
      limit: USERS_PER_PAGE,
      offset: (page - 1) * USERS_PER_PAGE,
      search: trimmedSearch || undefined
    })

    if (requestID !== usersRequestIdRef.current) {
      return
    }

    setUsers(res.data ?? [])
    setUsersPage(page)
    setUsersLoaded(true)
    setStatus(setUsersStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const loadDiscovery = async () => {
    setStatus(setDiscoveryStatus, { ...initialStatus, loading: true })
    const res = await api.searchComments()
    const sortedComments = [...(res.data ?? [])].sort((left, right) => {
      const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0
      const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0
      return rightTime - leftTime
    })

    setDiscoveryItems(sortedComments)
    setStatus(setDiscoveryStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const replaceCommentById = (items: Comment[], updatedComment: Comment) =>
    items.map((item) => (item.id === updatedComment.id ? updatedComment : item))

  const toggleCommentLike = async (comment: Comment) => {
    if (!comment.id) {
      return
    }

    const isLiked = likedCommentIds.has(comment.id)
    setStatus(setLikeStatus, { ...initialStatus, loading: true })

    const res = isLiked
      ? await api.unlikeComment({ comment_id: comment.id })
      : await api.likeComment({ comment_id: comment.id })

    let updatedComment = res.data as Comment | undefined
    if (!updatedComment && res.ok) {
      const commentRes = await api.getComment(comment.id)
      updatedComment = commentRes.data
    }

    if (updatedComment) {
      setFeedItems((prev) => replaceCommentById(prev, updatedComment as Comment))
      setDiscoveryItems((prev) => replaceCommentById(prev, updatedComment as Comment))
      setProfileComments((prev) => replaceCommentById(prev, updatedComment as Comment))
      setLikedCommentIds((prev) => {
        const next = new Set(prev)
        if (isLiked) {
          next.delete(comment.id as number)
        } else {
          next.add(comment.id as number)
        }
        return next
      })
    }

    setStatus(setLikeStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
  }

  const renderCommentLikeButton = (comment: Comment) => {
    const isLiked = comment.id ? likedCommentIds.has(comment.id) : false
    return (
      <button
        type="button"
        className={`comment-action comment-like-button ${isLiked ? 'liked' : ''}`}
        onClick={() => {
          void toggleCommentLike(comment)
        }}
        disabled={likeStatus.loading || !comment.id}
      >
        {Icons.heart(14, isLiked)} {comment.like_count ?? 0}
      </button>
    )
  }

  const registerUser = async () => {
    const username = newUser.username.trim()
    const password = newUser.password.trim()
    const email = newUser.email.trim()
    const firstName = newUser.first_name.trim()
    const lastName = newUser.last_name.trim()

    if (!username || !password || !email || !firstName || !lastName) {
      setStatus(setLoginStatus, { status: '', error: 'Complete all registration fields', loading: false })
      return
    }

    setStatus(setLoginStatus, { ...initialStatus, loading: true })
    const res = await api.createUser(newUser)
    setStatus(setLoginStatus, {
      status: `${res.status} ${res.statusText}`,
      error: res.error ?? '',
      loading: false
    })
    if (res.ok) {
      const createdUsername = username
      setNewUser({ username: '', first_name: '', last_name: '', password: '', email: '' })
      setLoginForm({ username: createdUsername, password: '' })
      setIsSignup(false)
      await loginWithBasicAuth({ username: createdUsername, password })
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
      navigateToView('users')
    }
  }

  const updateUserRoles = async () => {
    const username = userRoleUsername || selectedUser?.username
    if (!username) {
      setStatus(setUserRoleStatus, { status: '', error: 'Provide username', loading: false })
      return
    }
    const parsedRoles = rolesInput
      .split(',')
      .map((role) => role.trim())
      .filter(Boolean)
    setStatus(setUserRoleStatus, { ...initialStatus, loading: true })
    const res = await api.updateRolesForUser(username, parsedRoles)
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

  const sidebarUser = selectedUser ?? users.find((user) => user.username === currentUsername) ?? users[0] ?? null
  const canPost = currentUsername.trim().length > 0 && commentForm.content.trim().length > 0

  useEffect(() => {
    if (!loggedIn) {
      return
    }
    if (view === 'feed' && feedItems.length === 0 && !feedStatus.loading) {
      void loadFeed()
    }
  }, [loggedIn, view, feedItems.length, feedStatus.loading])

  useEffect(() => {
    if (!loggedIn) {
      return
    }
    if (view === 'discovery' && discoveryItems.length === 0 && !discoveryStatus.loading) {
      void loadDiscovery()
    }
  }, [discoveryItems.length, discoveryStatus.loading, loggedIn, view])

  useEffect(() => {
    if (!loggedIn) {
      return
    }
    if ((view === 'users' || view === 'profile') && !usersLoaded && !usersStatus.loading) {
      void loadUsers()
    }
  }, [loggedIn, view, usersLoaded, usersStatus.loading])

  useEffect(() => {
    if (!loggedIn || view !== 'users') {
      return
    }

    setUsersLoaded(false)
    const timer = window.setTimeout(() => {
      void loadUsers(1, userSearch)
    }, 200)

    return () => {
      window.clearTimeout(timer)
    }
  }, [api, loggedIn, userSearch, view])

  useEffect(() => {
    if (currentUsername || !token) {
      return
    }
    const inferredUsername = usernameFromToken(token)
    if (inferredUsername) {
      writeCookie(USERNAME_COOKIE, inferredUsername)
      setCurrentUsername(inferredUsername)
      if (!profileUsername) {
        setProfileUsername(inferredUsername)
      }
    }
  }, [currentUsername, profileUsername, token])

  useEffect(() => {
    if (!loggedIn || view !== 'profile' || !profileUsername) {
      return
    }
    void loadProfile(profileUsername)
  }, [loggedIn, view, profileUsername])

  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    const desiredPath = pathFromView(view)
    if (window.location.pathname !== desiredPath) {
      window.history.replaceState({}, '', desiredPath)
    }
  }, [view])

  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    const syncViewFromLocation = () => {
      setView(viewFromPath(window.location.pathname))
    }

    window.addEventListener('popstate', syncViewFromLocation)
    return () => {
      window.removeEventListener('popstate', syncViewFromLocation)
    }
  }, [])

  const renderProfileWorkspace = () => {
    if (!selectedUser || !profileForm) {
      return (
        <div className="empty-state">
          <div className="empty-state-text">Loading your profile.</div>
        </div>
      )
    }

    return (
        <div className="card animate-in">
        <div className="card-body">
          <div className="flex-between">
            <h3>Profile</h3>
          </div>
          {profileStatus.error ? <p className="error-text">{profileStatus.error}</p> : null}

          <div className="profile-header compact">
            <Avatar user={selectedUser} size="lg" />
            <div className="profile-info">
              <h2 className="profile-name">
                {selectedUser.first_name} {selectedUser.last_name}
              </h2>
              <div className="profile-handle">@{selectedUser.username}</div>
              <p className="profile-bio">Joined {formatDate(selectedUser.created_at)}</p>
              <div className="profile-stats">
                <span className="profile-stat">
                  <strong>{following.length}</strong>following
                </span>
                <span className="profile-stat">
                  <strong>{followers.length}</strong>followers
                </span>
                <span className="profile-stat">
                  <strong>{profileComments.length}</strong>posts
                </span>
              </div>
            </div>
          </div>

          <div className="form-grid">
            <div className="input-group">
              <span className="input-label">Username</span>
              <div className="input input-mono">{profileForm.username}</div>
            </div>
            <div className="input-group">
              <span className="input-label">Email</span>
              <div className="input">{profileForm.email}</div>
            </div>
            <div className="input-group">
              <span className="input-label">First name</span>
              <div className="input">{profileForm.first_name}</div>
            </div>
            <div className="input-group">
              <span className="input-label">Last name</span>
              <div className="input">{profileForm.last_name}</div>
            </div>
          </div>

          <div className="tabs mt-24">
            <button
              className={`tab ${profileTab === 'posts' ? 'active' : ''}`}
              onClick={() => setProfileTab('posts')}
            >
              Posts
            </button>
            <button
              className={`tab ${profileTab === 'followers' ? 'active' : ''}`}
              onClick={() => setProfileTab('followers')}
            >
              Followers
            </button>
            <button
              className={`tab ${profileTab === 'following' ? 'active' : ''}`}
              onClick={() => setProfileTab('following')}
            >
              Following
            </button>
            <button
              className={`tab ${profileTab === 'roles' ? 'active' : ''}`}
              onClick={() => setProfileTab('roles')}
            >
              Roles
            </button>
          </div>

          {profileTab === 'posts' && (
            <div className="animate-in">
              {profileComments.length > 0 ? (
                profileComments.map((comment) => (
                  <div key={comment.id ?? comment.content} className="comment-item compact">
                    <div className="comment-header">
                      <Avatar user={selectedUser} size="md" />
                      <div className="comment-meta">
                        <span className="comment-author">
                          {selectedUser.first_name} {selectedUser.last_name}
                        </span>
                        <span className="comment-handle">@{comment.username}</span>
                        <span className="comment-time">{timeAgo(comment.created_at)}</span>
                      </div>
                    </div>
                    <div className="comment-body">{comment.content}</div>
                    <div className="comment-actions">
                      {renderCommentLikeButton(comment)}
                    </div>
                  </div>
                ))
              ) : (
                <div className="empty-state">
                  <div className="empty-state-text">No posts yet</div>
                </div>
              )}
            </div>
          )}

          {profileTab === 'followers' && (
            <div className="animate-in" style={{ paddingTop: '16px' }}>
              {followers.length > 0 ? (
                followers.map((user) => (
                  <div key={user.username} className="user-row">
                    <Avatar user={user} size="md" />
                    <div className="user-row-info">
                      <div className="user-row-name">
                        {user.first_name} {user.last_name}
                      </div>
                      <div className="user-row-handle">@{user.username}</div>
                    </div>
                  </div>
                ))
              ) : (
                <div className="empty-state">
                  <div className="empty-state-text">No followers</div>
                </div>
              )}
            </div>
          )}

          {profileTab === 'following' && (
            <div className="animate-in" style={{ paddingTop: '16px' }}>
              {following.length > 0 ? (
                following.map((user) => (
                  <div key={user.username} className="user-row">
                    <Avatar user={user} size="md" />
                    <div className="user-row-info">
                      <div className="user-row-name">
                        {user.first_name} {user.last_name}
                      </div>
                      <div className="user-row-handle">@{user.username}</div>
                    </div>
                  </div>
                ))
              ) : (
                <div className="empty-state">
                  <div className="empty-state-text">No followed users</div>
                </div>
              )}
            </div>
          )}

          {profileTab === 'roles' && (
            <div className="animate-in" style={{ paddingTop: '16px' }}>
              <div className="tag-list mt-16">
                {rolesForUser.length > 0
                  ? rolesForUser.map((role) => (
                      <span key={role.name} className="badge badge-blue">
                        {role.name}
                      </span>
                    ))
                  : null}
                {rolesForUser.length === 0 ? <span className="small-muted">No roles assigned.</span> : null}
              </div>
            </div>
          )}
        </div>
      </div>
    )
  }

  if (!loggedIn) {
    return (
      <div className="login-page">
        <div className="login-container animate-in">
          <div className="login-brand">
            <h1>Socialapp</h1>
            <p>A space for thoughtful conversation</p>
          </div>
          <div className="login-card">
            <h2 className="login-title">{isSignup ? 'Create account' : 'Welcome back'}</h2>
            <form
              onSubmit={(event) => {
                event.preventDefault()
                if (isSignup) {
                  void registerUser()
                  return
                }
                void loginWithBasicAuth()
              }}
            >
              {isSignup ? (
                <>
                  <div className="form-grid">
                    <label className="input-group">
                      <span className="input-label">First name</span>
                      <input
                        className="input"
                        value={newUser.first_name}
                        onChange={(event) => setNewUser((prev) => ({ ...prev, first_name: event.target.value }))}
                        autoComplete="given-name"
                        placeholder="John"
                      />
                    </label>
                    <label className="input-group">
                      <span className="input-label">Last name</span>
                      <input
                        className="input"
                        value={newUser.last_name}
                        onChange={(event) => setNewUser((prev) => ({ ...prev, last_name: event.target.value }))}
                        autoComplete="family-name"
                        placeholder="Doe"
                      />
                    </label>
                  </div>

                  <label className="input-group">
                    <span className="input-label">Email</span>
                    <input
                      className="input"
                      type="email"
                      value={newUser.email}
                      onChange={(event) => setNewUser((prev) => ({ ...prev, email: event.target.value }))}
                      autoComplete="email"
                      placeholder="john@example.com"
                    />
                  </label>
                </>
              ) : null}

              <label className="input-group">
                <span className="input-label">Username</span>
                <input
                  className="input"
                  value={isSignup ? newUser.username : loginForm.username}
                  onChange={(event) => {
                    const { value } = event.target
                    if (isSignup) {
                      setNewUser((prev) => ({ ...prev, username: value }))
                      return
                    }
                    setLoginForm((prev) => ({ ...prev, username: value }))
                  }}
                  autoComplete={isSignup ? 'new-username' : 'username'}
                  placeholder="some2"
                />
              </label>

              <label className="input-group">
                <span className="input-label">Password</span>
                <input
                  className="input"
                  type="password"
                  value={isSignup ? newUser.password : loginForm.password}
                  onChange={(event) => {
                    const { value } = event.target
                    if (isSignup) {
                      setNewUser((prev) => ({ ...prev, password: value }))
                      return
                    }
                    setLoginForm((prev) => ({ ...prev, password: value }))
                  }}
                  autoComplete={isSignup ? 'new-password' : 'current-password'}
                  placeholder="********"
                />
              </label>

              <button
                type="submit"
                className="btn btn-primary login-submit"
                disabled={loginStatus.loading}
              >
                {loginStatus.loading ? 'Signing in...' : isSignup ? 'Create account' : 'Sign in'}
              </button>
            </form>

            <div className="login-footer">
              {isSignup ? (
                <>
                  Already have an account?{' '}
                  <button type="button" className="link-button" onClick={() => setIsSignup(false)}>
                    Sign in
                  </button>
                </>
              ) : (
                <>
                  Do not have an account?{' '}
                  <button type="button" className="link-button" onClick={() => setIsSignup(true)}>
                    Create one
                  </button>
                </>
              )}
            </div>

            {loginStatus.error ? <p className="error-text">Error: {loginStatus.error}</p> : null}
            <p className="small-muted">{token ? 'Stored token detected.' : 'No access token stored yet.'}</p>
          </div>
        </div>
      </div>
    )
  }

  const navItems: Array<{ id: View; label: string; icon: (size?: number) => JSX.Element }> = [
    { id: 'feed', label: 'Feed', icon: Icons.home },
    { id: 'discovery', label: 'Discovery', icon: Icons.search },
    { id: 'users', label: 'Users', icon: Icons.users },
    { id: 'profile', label: 'Profile', icon: Icons.user },
    { id: 'settings', label: 'Settings', icon: Icons.settings }
  ]

  return (
    <div className="app-root">
      <aside className="sidebar">
        <button
          type="button"
          className="sidebar-brand"
          onClick={() => {
            navigateToView('feed')
          }}
        >
          <h1>Socialapp</h1>
        </button>

        <nav className="sidebar-nav">
          <div className="sidebar-section-label">Navigation</div>
          {navItems.map((item) => (
            <button
              key={item.id}
              className={`nav-item ${view === item.id ? 'active' : ''}`}
              onClick={() => {
                navigateToView(item.id)
                if (item.id === 'profile' && currentUsername) {
                  setProfileUsername(currentUsername)
                }
                if (item.id === 'admin') {
                  void loadRoles()
                }
                if (item.id === 'users') {
                  void loadUsers()
                }
              }}
            >
              <span className="nav-icon">{item.icon(18)}</span>
              {item.label}
            </button>
          ))}
        </nav>

        <div className="sidebar-user">
          {sidebarUser ? <Avatar user={sidebarUser} size="sm" /> : <div className="avatar avatar-sm">?</div>}
          <div className="sidebar-user-info">
            <div className="sidebar-user-name">
              {sidebarUser ? `${sidebarUser.first_name} ${sidebarUser.last_name}` : 'Guest user'}
            </div>
            <div className="sidebar-user-handle">{sidebarUser ? `@${sidebarUser.username}` : '@guest'}</div>
          </div>
          <button className="btn btn-ghost btn-icon" style={{ padding: '4px' }} onClick={() => handleTokenChange(undefined)}>
            {Icons.logout(16)}
          </button>
        </div>
      </aside>

      <main className="main-content">
        {view === 'feed' && (
          <div>
            <div className="page-header">
              <div className="flex-between">
                <div>
                  <h2 className="page-title">Feed</h2>
                  <p className="page-subtitle">Posts from people in your network</p>
                </div>
              </div>
            </div>

            <div className="compose-box animate-in">
              <Avatar
                user={{
                  username: currentUsername || 'you',
                  first_name: currentUsername || 'Y',
                  last_name: ''
                }}
                size="md"
              />
              <div className="compose-input">
                <textarea
                  className="compose-textarea"
                  rows={3}
                  placeholder="What's on your mind?"
                  value={commentForm.content}
                  onChange={(event) => setCommentForm((prev) => ({ ...prev, content: event.target.value }))}
                />
                <div className="compose-footer">
                  <button
                    className="btn btn-primary btn-sm"
                    onClick={() => {
                      void postComment()
                    }}
                    disabled={!canPost || isPostingComment}
                  >
                    {isPostingComment ? 'Posting...' : 'Post'}
                  </button>
                </div>
                {feedStatus.error ? <p className="error-text">{feedStatus.error}</p> : null}
                {likeStatus.error ? <p className="error-text">{likeStatus.error}</p> : null}
                {postedCommentNotice ? (
                  <pre className="highlight">
                    Posted comment by @{postedCommentNotice.username}
                    {'\n'}
                    {postedCommentNotice.content}
                  </pre>
                ) : null}
              </div>
            </div>

            <div className="card card-flat">
              {feedItems.map((item, index) => {
                const knownUser = users.find((candidate) => candidate.username === item.username)
                const avatarUser =
                  knownUser ??
                  ({
                    username: item.username,
                    first_name: item.username,
                    last_name: ''
                  } as const)

                return (
                  <div key={`${item.username}-${item.id ?? item.content}-${index}`} className="comment-item">
                    <div className="comment-header">
                      <Avatar user={avatarUser} size="md" />
                      <div className="comment-meta">
                        <span className="comment-author">
                          {knownUser ? `${knownUser.first_name} ${knownUser.last_name}` : item.username}
                        </span>
                        <span className="comment-handle">@{item.username}</span>
                        <span className="comment-time">{timeAgo(item.created_at)}</span>
                      </div>
                    </div>
                    <div className="comment-body">{item.content}</div>
                    <div className="comment-actions">
                      {renderCommentLikeButton(item)}
                    </div>
                  </div>
                )
              })}
              {feedItems.length === 0 ? (
                <div className="empty-state">
                  <div className="empty-state-text">No feed items yet.</div>
                </div>
              ) : null}
            </div>

          </div>
        )}

        {view === 'discovery' && (
          <div>
            <div className="page-header">
              <div className="flex-between">
                <div>
                  <h2 className="page-title">Discovery</h2>
                  <p className="page-subtitle">Latest public posts happening across the app</p>
                </div>
              </div>
            </div>

            <div className="card card-flat">
              {discoveryStatus.error ? <p className="error-text">{discoveryStatus.error}</p> : null}
              {likeStatus.error ? <p className="error-text">{likeStatus.error}</p> : null}
              {discoveryItems.map((item, index) => {
                const knownUser = users.find((candidate) => candidate.username === item.username)
                const avatarUser =
                  knownUser ??
                  ({
                    username: item.username,
                    first_name: item.username,
                    last_name: ''
                  } as const)

                return (
                  <div key={`${item.username}-${item.id ?? item.content}-${index}`} className="comment-item">
                    <div className="comment-header">
                      <Avatar user={avatarUser} size="md" />
                      <div className="comment-meta">
                        <span className="comment-author">
                          {knownUser ? `${knownUser.first_name} ${knownUser.last_name}` : item.username}
                        </span>
                        <span className="comment-handle">@{item.username}</span>
                        <span className="comment-time">{timeAgo(item.created_at)}</span>
                      </div>
                    </div>
                    <div className="comment-body">{item.content}</div>
                    <div className="comment-actions">
                      {renderCommentLikeButton(item)}
                    </div>
                  </div>
                )
              })}
              {discoveryItems.length === 0 ? (
                <div className="empty-state">
                  <div className="empty-state-text">No public posts discovered yet.</div>
                </div>
              ) : null}
            </div>
          </div>
        )}

        {view === 'users' && (
          <div>
            <div className="page-header">
              <div className="flex-between">
                <div>
                  <h2 className="page-title">Users</h2>
                  <p className="page-subtitle">Discover and manage users in the network</p>
                </div>
              </div>
            </div>

            <div className="page-body">
              <div className="search-bar animate-in">
                {Icons.search(18)}
                <input
                  placeholder="Search users..."
                  value={userSearch}
                  onChange={(event) => setUserSearch(event.target.value)}
                />
              </div>

              <div className="content-grid">
                <div className="card animate-in delay-1">
                  <div className="card-body">
                    <div className="flex-between">
                      <h3>Directory</h3>
                      <button
                        className="btn btn-secondary btn-sm"
                        onClick={() => {
                          void loadUsers(usersPage, userSearch)
                        }}
                        disabled={usersStatus.loading}
                      >
                        {usersStatus.loading ? 'Loading...' : 'Refresh'}
                      </button>
                    </div>
                    {usersStatus.error ? <p className="error-text">{usersStatus.error}</p> : null}

                    {users.length > 0 ? (
                      users.map((user) => (
                        <button
                          key={user.username}
                          className="user-row selectable"
                          onClick={() => {
                            setProfileUsername(user.username)
                            navigateToView('profile')
                          }}
                        >
                          <Avatar user={user} size="md" />
                          <span className="user-row-info">
                            <span className="user-row-name">
                              {user.first_name} {user.last_name}
                            </span>
                            <span className="user-row-handle">@{user.username}</span>
                            <span className="user-row-bio">{user.email}</span>
                          </span>
                        </button>
                      ))
                    ) : (
                      <div className="empty-state">
                        <div className="empty-state-text">No users found</div>
                      </div>
                    )}
                    <div className="flex gap-8 mt-16" style={{ flexWrap: 'wrap', justifyContent: 'center' }}>
                      {Array.from({ length: USERS_PAGE_COUNT }, (_, index) => index + 1).map((page) => (
                        <button
                          key={page}
                          className={`btn btn-sm ${usersPage === page ? 'btn-primary' : 'btn-secondary'}`}
                          onClick={() => {
                            void loadUsers(page, userSearch)
                          }}
                          disabled={usersStatus.loading}
                        >
                          {page}
                        </button>
                      ))}
                    </div>
                  </div>
                </div>

              </div>

            </div>
          </div>
        )}

        {view === 'profile' && (
          <div>
            <div className="page-header">
              <div className="flex-between">
                <div>
                  <h2 className="page-title">Profile</h2>
                  <p className="page-subtitle">Your account details, followers, and roles</p>
                </div>
              </div>
            </div>

            <div className="page-body">
              {renderProfileWorkspace()}
            </div>
          </div>
        )}

        {view === 'settings' && (
          <div>
            <div className="page-header">
              <h2 className="page-title">Settings</h2>
              <p className="page-subtitle">Manage API access and account credentials</p>
            </div>

            <div className="page-body">
              <div className="tabs animate-in" style={{ marginBottom: '32px' }}>
                <button
                  className={`tab ${settingsTab === 'auth' ? 'active' : ''}`}
                  onClick={() => setSettingsTab('auth')}
                >
                  Auth
                </button>
                <button
                  className={`tab ${settingsTab === 'password' ? 'active' : ''}`}
                  onClick={() => setSettingsTab('password')}
                >
                  Password
                </button>
                <button
                  className={`tab ${settingsTab === 'danger' ? 'active' : ''}`}
                  onClick={() => setSettingsTab('danger')}
                >
                  Danger zone
                </button>
              </div>

              {settingsTab === 'auth' && (
                <div className="animate-in delay-1">
                  <div className="card">
                    <div className="card-body">
                      <p className="small-muted">API endpoint configured via frontend environment.</p>
                      <div className="tag-list mt-16">
                        <span className="badge badge-default">{token ? 'Token set' : 'No token'}</span>
                      </div>
                    </div>
                  </div>

                  <div className="card mt-24">
                    <div className="card-body">
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
                  </div>
                </div>
              )}

              {settingsTab === 'password' && (
                <div className="animate-in delay-1 content-grid">
                  <div className="card">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Change password</h3>
                        <button
                          className="btn btn-primary btn-sm"
                          onClick={() => {
                            void changePassword()
                          }}
                        >
                          Update
                        </button>
                      </div>
                      <label className="input-group">
                        <span className="input-label">Old password</span>
                        <input
                          className="input"
                          type="password"
                          value={passwordForm.old_password}
                          onChange={(event) =>
                            setPasswordForm((prev) => ({ ...prev, old_password: event.target.value }))
                          }
                        />
                      </label>
                      <label className="input-group">
                        <span className="input-label">New password</span>
                        <input
                          className="input"
                          type="password"
                          value={passwordForm.new_password}
                          onChange={(event) =>
                            setPasswordForm((prev) => ({ ...prev, new_password: event.target.value }))
                          }
                        />
                      </label>
                    </div>
                  </div>

                  <div className="card">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Reset password</h3>
                        <button
                          className="btn btn-secondary btn-sm"
                          onClick={() => {
                            void resetPassword()
                          }}
                        >
                          Send reset
                        </button>
                      </div>
                      <label className="input-group">
                        <span className="input-label">Email</span>
                        <input
                          className="input"
                          type="email"
                          value={resetForm.email}
                          onChange={(event) => setResetForm({ email: event.target.value })}
                        />
                      </label>
                    </div>
                  </div>

                  {passwordStatus.error ? (
                    <div className="card">
                      <div className="card-body">
                        <h3>Password errors</h3>
                        <p className="error-text">{passwordStatus.error}</p>
                      </div>
                    </div>
                  ) : null}
                </div>
              )}

              {settingsTab === 'danger' && (
                <div className="animate-in delay-1" style={{ maxWidth: '520px' }}>
                  <div className="card danger-card">
                    <div className="card-body">
                      <h3>Delete account</h3>
                      <p className="small-muted">
                        Deleting a user account is permanent and removes profile data, posts, and relationships.
                      </p>
                      <button
                        className="btn btn-danger"
                        onClick={() => {
                          void deleteProfile()
                        }}
                        disabled={!selectedUser}
                      >
                        Delete my account
                      </button>
                      {!selectedUser ? (
                        <p className="small-muted">Load a profile first from Users or Profile page.</p>
                      ) : null}
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {view === 'admin' && (
          <div>
            <div className="page-header">
              <h2 className="page-title">Admin</h2>
              <p className="page-subtitle">Roles, scopes, and URL management</p>
            </div>

            <div className="page-body">
              <div className="tabs animate-in" style={{ marginBottom: '32px' }}>
                <button className={`tab ${adminTab === 'roles' ? 'active' : ''}`} onClick={() => setAdminTab('roles')}>
                  Roles
                </button>
                <button className={`tab ${adminTab === 'scopes' ? 'active' : ''}`} onClick={() => setAdminTab('scopes')}>
                  Scopes
                </button>
                <button className={`tab ${adminTab === 'urls' ? 'active' : ''}`} onClick={() => setAdminTab('urls')}>
                  Short URLs
                </button>
              </div>

              {adminTab === 'roles' && (
                <div className="animate-in delay-1">
                  <div className="card">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Roles</h3>
                        <div className="flex gap-8">
                          <button
                            className="btn btn-secondary btn-sm"
                            onClick={() => {
                              void loadRoles()
                            }}
                          >
                            Refresh
                          </button>
                          <button
                            className="btn btn-secondary btn-sm"
                            onClick={() => {
                              void fetchRole()
                            }}
                          >
                            Fetch
                          </button>
                          <button
                            className="btn btn-primary btn-sm"
                            onClick={() => {
                              void createRole()
                            }}
                          >
                            Create
                          </button>
                          <button
                            className="btn btn-secondary btn-sm"
                            onClick={() => {
                              void updateRole()
                            }}
                          >
                            Update
                          </button>
                          <button
                            className="btn btn-danger btn-sm"
                            onClick={() => {
                              void deleteRole()
                            }}
                          >
                            Delete
                          </button>
                        </div>
                      </div>

                      {rolesStatus.error ? <p className="error-text">{rolesStatus.error}</p> : null}

                      <div className="form-grid">
                        <label className="input-group">
                          <span className="input-label">Role ID (for update/delete)</span>
                          <input
                            className="input input-mono"
                            value={roleId}
                            onChange={(event) => setRoleId(event.target.value)}
                          />
                        </label>
                        <label className="input-group">
                          <span className="input-label">Name</span>
                          <input
                            className="input input-mono"
                            value={roleForm.name}
                            onChange={(event) =>
                              setRoleForm((prev) => ({ ...prev, name: event.target.value }))
                            }
                          />
                        </label>
                        <label className="input-group">
                          <span className="input-label">Description</span>
                          <input
                            className="input"
                            value={roleForm.description ?? ''}
                            onChange={(event) =>
                              setRoleForm((prev) => ({ ...prev, description: event.target.value }))
                            }
                          />
                        </label>
                      </div>

                      <div className="table-wrap mt-24">
                        <table>
                          <thead>
                            <tr>
                              <th>Name</th>
                              <th>Description</th>
                              <th>ID</th>
                            </tr>
                          </thead>
                          <tbody>
                            {roles.map((role) => (
                              <tr key={role.id ?? role.name}>
                                <td>
                                  <span className="badge badge-blue">{role.name}</span>
                                </td>
                                <td>{role.description}</td>
                                <td className="text-mono">{role.id ?? '-'}</td>
                              </tr>
                            ))}
                            {roles.length === 0 ? (
                              <tr>
                                <td colSpan={3}>No roles loaded.</td>
                              </tr>
                            ) : null}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </div>

                  <div className="card mt-24">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>User role assignment</h3>
                        <div className="flex gap-8">
                          <button
                            className="btn btn-secondary btn-sm"
                            onClick={() => {
                              void fetchUserRoles()
                            }}
                          >
                            Fetch roles
                          </button>
                          <button
                            className="btn btn-primary btn-sm"
                            onClick={() => {
                              void updateUserRoles()
                            }}
                          >
                            Update roles
                          </button>
                        </div>
                      </div>

                      {userRoleStatus.error ? <p className="error-text">{userRoleStatus.error}</p> : null}

                      <div className="form-grid">
                        <label className="input-group">
                          <span className="input-label">Username</span>
                          <input
                            className="input input-mono"
                            value={userRoleUsername}
                            onChange={(event) => setUserRoleUsername(event.target.value)}
                          />
                        </label>
                        <label className="input-group">
                          <span className="input-label">Role names (comma-separated)</span>
                          <input
                            className="input input-mono"
                            value={rolesInput}
                            onChange={(event) => setRolesInput(event.target.value)}
                          />
                        </label>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {adminTab === 'scopes' && (
                <div className="animate-in delay-1">
                  <div className="card">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Scopes</h3>
                        <div className="flex gap-8">
                          <button
                            className="btn btn-secondary btn-sm"
                            onClick={() => {
                              void loadScopes()
                            }}
                          >
                            Refresh
                          </button>
                          <button
                            className="btn btn-secondary btn-sm"
                            onClick={() => {
                              void fetchScope()
                            }}
                          >
                            Fetch
                          </button>
                          <button
                            className="btn btn-primary btn-sm"
                            onClick={() => {
                              void createScope()
                            }}
                          >
                            Create
                          </button>
                          <button
                            className="btn btn-secondary btn-sm"
                            onClick={() => {
                              void updateScope()
                            }}
                          >
                            Update
                          </button>
                          <button
                            className="btn btn-danger btn-sm"
                            onClick={() => {
                              void deleteScope()
                            }}
                          >
                            Delete
                          </button>
                        </div>
                      </div>

                      {scopesStatus.error ? <p className="error-text">{scopesStatus.error}</p> : null}

                      <div className="form-grid">
                        <label className="input-group">
                          <span className="input-label">Scope ID (for update/delete)</span>
                          <input
                            className="input input-mono"
                            value={scopeId}
                            onChange={(event) => setScopeId(event.target.value)}
                          />
                        </label>
                        <label className="input-group">
                          <span className="input-label">Name</span>
                          <input
                            className="input input-mono"
                            value={scopeForm.name}
                            onChange={(event) =>
                              setScopeForm((prev) => ({ ...prev, name: event.target.value }))
                            }
                          />
                        </label>
                        <label className="input-group">
                          <span className="input-label">Description</span>
                          <input
                            className="input"
                            value={scopeForm.description}
                            onChange={(event) =>
                              setScopeForm((prev) => ({ ...prev, description: event.target.value }))
                            }
                          />
                        </label>
                      </div>

                      <div className="table-wrap mt-24">
                        <table>
                          <thead>
                            <tr>
                              <th>Scope</th>
                              <th>Description</th>
                              <th>ID</th>
                            </tr>
                          </thead>
                          <tbody>
                            {scopesList.map((scope) => (
                              <tr key={scope.id ?? scope.name}>
                                <td className="text-mono">{scope.name}</td>
                                <td>{scope.description}</td>
                                <td className="text-mono">{scope.id ?? '-'}</td>
                              </tr>
                            ))}
                            {scopesList.length === 0 ? (
                              <tr>
                                <td colSpan={3}>No scopes loaded.</td>
                              </tr>
                            ) : null}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </div>

                  <div className="card mt-24">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Role scopes</h3>
                        <div className="flex gap-8">
                          <button
                            className="btn btn-secondary btn-sm"
                            onClick={() => {
                              void loadRoleScopes()
                            }}
                          >
                            Load
                          </button>
                          <button
                            className="btn btn-primary btn-sm"
                            onClick={() => {
                              void addScopesToRole()
                            }}
                          >
                            Add scopes
                          </button>
                          <button
                            className="btn btn-danger btn-sm"
                            onClick={() => {
                              void removeScopeFromRole()
                            }}
                          >
                            Remove scope
                          </button>
                        </div>
                      </div>

                      {roleScopesStatus.error ? <p className="error-text">{roleScopesStatus.error}</p> : null}

                      <div className="form-grid">
                        <label className="input-group">
                          <span className="input-label">Role ID</span>
                          <input
                            className="input input-mono"
                            value={roleScopeId}
                            onChange={(event) => setRoleScopeId(event.target.value)}
                          />
                        </label>
                        <label className="input-group">
                          <span className="input-label">Scope names (comma-separated)</span>
                          <input
                            className="input input-mono"
                            value={roleScopeInput}
                            onChange={(event) => setRoleScopeInput(event.target.value)}
                          />
                        </label>
                        <label className="input-group">
                          <span className="input-label">Remove scope ID</span>
                          <input
                            className="input input-mono"
                            value={removeScopeId}
                            onChange={(event) => setRemoveScopeId(event.target.value)}
                          />
                        </label>
                      </div>

                      <div className="tag-list mt-16">
                        {roleScopes.map((scope) => (
                          <span key={scope.id ?? scope.name} className="badge badge-green">
                            {scope.name}
                          </span>
                        ))}
                        {roleScopes.length === 0 ? <span className="small-muted">No role scopes loaded.</span> : null}
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {adminTab === 'urls' && (
                <div className="animate-in delay-1 content-grid">
                  <div className="card">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Create URL</h3>
                        <button
                          className="btn btn-primary btn-sm"
                          onClick={() => {
                            void createUrl()
                          }}
                        >
                          Create
                        </button>
                      </div>
                      <label className="input-group">
                        <span className="input-label">Alias</span>
                        <input
                          className="input input-mono"
                          value={urlForm.alias}
                          onChange={(event) =>
                            setUrlForm((prev) => ({ ...prev, alias: event.target.value }))
                          }
                        />
                      </label>
                      <label className="input-group">
                        <span className="input-label">Target URL</span>
                        <input
                          className="input"
                          value={urlForm.url}
                          onChange={(event) => setUrlForm((prev) => ({ ...prev, url: event.target.value }))}
                        />
                      </label>
                    </div>
                  </div>

                  <div className="card">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Resolve</h3>
                        <button
                          className="btn btn-secondary btn-sm"
                          onClick={() => {
                            void resolveUrl()
                          }}
                        >
                          Resolve
                        </button>
                      </div>
                      <label className="input-group">
                        <span className="input-label">Alias</span>
                        <input
                          className="input input-mono"
                          value={urlLookupAlias}
                          onChange={(event) => setUrlLookupAlias(event.target.value)}
                        />
                      </label>
                      <p className="small-muted">{urlResolve ? `Location: ${urlResolve}` : 'No redirect yet.'}</p>
                    </div>
                  </div>

                  <div className="card">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Metadata</h3>
                        <button
                          className="btn btn-secondary btn-sm"
                          onClick={() => {
                            void fetchUrlData()
                          }}
                        >
                          Fetch
                        </button>
                      </div>
                      <label className="input-group">
                        <span className="input-label">Alias</span>
                        <input
                          className="input input-mono"
                          value={urlLookupAlias}
                          onChange={(event) => setUrlLookupAlias(event.target.value)}
                        />
                      </label>
                      {urlData ? (
                        <pre className="highlight">
                          {urlData.alias}
                          {' -> '}
                          {urlData.url}
                        </pre>
                      ) : (
                        <p className="small-muted">No metadata loaded.</p>
                      )}
                    </div>
                  </div>

                  <div className="card">
                    <div className="card-body">
                      <div className="flex-between">
                        <h3>Delete</h3>
                        <button
                          className="btn btn-danger btn-sm"
                          onClick={() => {
                            void deleteUrl()
                          }}
                        >
                          Delete
                        </button>
                      </div>
                      <label className="input-group">
                        <span className="input-label">Alias</span>
                        <input
                          className="input input-mono"
                          value={urlDeleteAlias}
                          onChange={(event) => setUrlDeleteAlias(event.target.value)}
                        />
                      </label>
                    </div>
                  </div>

                  {(urlStatus.error || urlData) ? (
                    <div className="card">
                      <div className="card-body">
                        {urlStatus.error ? <p className="error-text">{urlStatus.error}</p> : null}
                        {urlData ? (
                          <p className="small-muted">Alias {urlData.alias} points to {urlData.url}</p>
                        ) : null}
                      </div>
                    </div>
                  ) : null}
                </div>
              )}
            </div>
          </div>
        )}
      </main>
    </div>
  )
}
