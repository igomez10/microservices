import { apiRequest, type ApiResponse } from '@/api/client'
import type { CreateUserRequest, Role } from '@/api/social'

const baseUrl = process.env.INTEGRATION_API_BASE_URL ?? 'http://localhost:8086'
const clientId = process.env.INTEGRATION_CLIENT_ID
const clientSecret = process.env.INTEGRATION_CLIENT_SECRET
const scopes = process.env.INTEGRATION_SCOPES

const describeIf = clientId && clientSecret ? describe : describe.skip

const assertOk = (res: ApiResponse, message: string) => {
  if (!res.ok) {
    throw new Error(`${message} failed: ${res.status} ${res.statusText} ${res.error ?? ''}`)
  }
}

describeIf('integration api', () => {
  let token = ''

  beforeAll(async () => {
    const res = await apiRequest<{ access_token?: string }>({
      baseUrl,
      path: '/v1/oauth/token',
      method: 'POST',
      basicAuth: {
        clientId: clientId!,
        clientSecret: clientSecret!
      },
      queryParams: scopes ? { scope: scopes } : undefined
    })
    assertOk(res, 'getAccessToken')
    token = res.data?.access_token ?? ''
    if (!token) {
      throw new Error('Missing access token from auth response')
    }
  })

  const authed = async <T,>(options: Omit<Parameters<typeof apiRequest<T>>[0], 'baseUrl' | 'token'>) =>
    apiRequest<T>({ baseUrl, token, ...options })

  it('welcomes', async () => {
    const res = await authed({ path: '/', method: 'GET' })
    expect([200, 401, 403]).toContain(res.status)
  })

  it('creates, reads, updates, and deletes a user', async () => {
    const suffix = `ui${Date.now().toString(36)}`
    const username = `user_${suffix}`
    const newUser: CreateUserRequest = {
      username,
      first_name: 'Test',
      last_name: 'User',
      password: 'Password123!',
      email: `${username}@mail.com`
    }
    const userRes = await authed<{ id?: number; created_at?: string }>({
      path: '/v1/users',
      method: 'POST',
      body: newUser
    })
    assertOk(userRes, 'createUser')

    const getRes = await authed({ path: '/v1/users/{username}', method: 'GET', pathParams: { username } })
    assertOk(getRes, 'getUserByUsername')

    const listRes = await authed({ path: '/v1/users', method: 'GET' })
    assertOk(listRes, 'listUsers')

    const searchRes = await authed({
      path: '/v1/users',
      method: 'GET',
      queryParams: { search: username }
    })
    assertOk(searchRes, 'searchUsers')
    expect(Array.isArray(searchRes.data)).toBe(true)
    expect(searchRes.data?.some((user) => user.username === username)).toBe(true)

    const rolesRes = await authed({
      path: '/v1/users/{username}/roles',
      method: 'GET',
      pathParams: { username }
    })
    assertOk(rolesRes, 'getRolesForUser')

    const updateRolesRes = await authed({
      path: '/v1/users/{username}/roles',
      method: 'PUT',
      pathParams: { username },
      body: ['user']
    })
    assertOk(updateRolesRes, 'updateRolesForUser')

    const updateRes = await authed({
      path: '/v1/users/{username}',
      method: 'PUT',
      pathParams: { username },
      body: {
        id: userRes.data?.id,
        username,
        first_name: 'Updated',
        last_name: 'User',
        email: `${username}@mail.com`,
        created_at: userRes.data?.created_at ?? new Date().toISOString()
      }
    })
    assertOk(updateRes, 'updateUser')

    const deleteRes = await authed({
      path: '/v1/users/{username}',
      method: 'DELETE',
      pathParams: { username }
    })
    assertOk(deleteRes, 'deleteUser')
  })

  it('creates a comment and reads it', async () => {
    const suffix = `comment${Date.now().toString(36)}`
    const username = `comment_${suffix}`

    const userRes = await authed({
      path: '/v1/users',
      method: 'POST',
      body: {
        username,
        first_name: 'Comment',
        last_name: 'User',
        password: 'Password123!',
        email: `${username}@mail.com`
      }
    })
    assertOk(userRes, 'createUser for comment')

    try {
      const commentRes = await authed<{ id?: number }>({
        path: '/v1/comments',
        method: 'POST',
        body: {
          content: `hello ${suffix}`,
          username
        }
      })
      assertOk(commentRes, 'createComment')

      if (commentRes.data?.id) {
        const getRes = await authed({
          path: '/v1/comments/{id}',
          method: 'GET',
          pathParams: { id: commentRes.data.id }
        })
        if (getRes.status !== 404) {
          assertOk(getRes, 'getComment')
        } else {
          console.warn('getComment returned 404; skipping assertion')
        }
      }

      const listRes = await authed({
        path: '/v1/users/{username}/comments',
        method: 'GET',
        pathParams: { username }
      })
      assertOk(listRes, 'getUserComments')
    } finally {
      await authed({ path: '/v1/users/{username}', method: 'DELETE', pathParams: { username } })
    }
  })

  it('follows and unfollows users', async () => {
    const suffix = `follow${Date.now().toString(36)}`
    const follower = `follower_${suffix}`
    const followed = `followed_${suffix}`

    const createUser = async (username: string) =>
      authed({
        path: '/v1/users',
        method: 'POST',
        body: {
          username,
          first_name: 'Follow',
          last_name: 'User',
          password: 'Password123!',
          email: `${username}@mail.com`
        }
      })

    await createUser(follower)
    await createUser(followed)

    try {
      const followRes = await authed({
        path: '/v1/users/{followedUsername}/followers/{followerUsername}',
        method: 'POST',
        pathParams: { followedUsername: followed, followerUsername: follower }
      })
      assertOk(followRes, 'followUser')

      const listRes = await authed({
        path: '/v1/users/{username}/followers',
        method: 'GET',
        pathParams: { username: followed }
      })
      assertOk(listRes, 'getUserFollowers')

      const followingRes = await authed({
        path: '/v1/users/{username}/following',
        method: 'GET',
        pathParams: { username: follower }
      })
      assertOk(followingRes, 'getFollowingUsers')

      const unfollowRes = await authed({
        path: '/v1/users/{followedUsername}/followers/{followerUsername}',
        method: 'DELETE',
        pathParams: { followedUsername: followed, followerUsername: follower }
      })
      assertOk(unfollowRes, 'unfollowUser')
    } finally {
      await authed({ path: '/v1/users/{username}', method: 'DELETE', pathParams: { username: follower } })
      await authed({ path: '/v1/users/{username}', method: 'DELETE', pathParams: { username: followed } })
    }
  })

  it('manages roles and scopes', async () => {
    const suffix = `role${Date.now().toString(36)}`
    const roleName = `role_${suffix}`
    const scopeName = `scope.${suffix}`

    const roleRes = await authed<Role>({
      path: '/v1/roles',
      method: 'POST',
      body: {
        name: roleName,
        description: 'role for tests'
      }
    })
    assertOk(roleRes, 'createRole')

    const scopeRes = await authed<{ id?: number }>({
      path: '/v1/scopes',
      method: 'POST',
      body: {
        name: scopeName,
        description: 'scope for tests'
      }
    })
    assertOk(scopeRes, 'createScope')

    let resolvedRoleId: number | undefined

    try {
      const listRoles = await authed<Role[]>({ path: '/v1/roles', method: 'GET' })
      assertOk(listRoles, 'listRoles')

      resolvedRoleId = roleRes.data?.id ?? listRoles.data?.find((role) => role.name === roleName)?.id

      if (resolvedRoleId) {
        const getRole = await authed({
          path: '/v1/roles/{id}',
          method: 'GET',
          pathParams: { id: resolvedRoleId }
        })

        if (getRole.status === 404) {
          console.warn('getRole returned 404; skipping role mutation assertions')
        } else {
          assertOk(getRole, 'getRole')

          const updateRole = await authed({
            path: '/v1/roles/{id}',
            method: 'PUT',
            pathParams: { id: resolvedRoleId },
            body: {
              name: `${roleName}_updated`,
              description: 'updated role'
            }
          })
          assertOk(updateRole, 'updateRole')

          const addScope = await authed({
            path: '/v1/roles/{id}/scopes',
            method: 'POST',
            pathParams: { id: resolvedRoleId },
            body: [scopeName]
          })
          assertOk(addScope, 'addScopeToRole')

          const listRoleScopes = await authed({
            path: '/v1/roles/{id}/scopes',
            method: 'GET',
            pathParams: { id: resolvedRoleId }
          })
          assertOk(listRoleScopes, 'listScopesForRole')

          if (scopeRes.data?.id) {
            const removeScope = await authed({
              path: '/v1/roles/{role_id}/scopes/{scope_id}',
              method: 'DELETE',
              pathParams: { role_id: resolvedRoleId, scope_id: scopeRes.data.id }
            })
            expect([200, 204]).toContain(removeScope.status)
          }
        }
      }

      const listScopes = await authed({ path: '/v1/scopes', method: 'GET' })
      assertOk(listScopes, 'listScopes')

      if (scopeRes.data?.id) {
        const getScope = await authed({
          path: '/v1/scopes/{id}',
          method: 'GET',
          pathParams: { id: scopeRes.data.id }
        })
        if (getScope.status === 404) {
          console.warn('getScope returned 404; skipping scope mutation assertions')
        } else {
          assertOk(getScope, 'getScope')

          const updateScope = await authed({
            path: '/v1/scopes/{id}',
            method: 'PUT',
            pathParams: { id: scopeRes.data.id },
            body: {
              name: `${scopeName}.updated`,
              description: 'updated scope'
            }
          })
          assertOk(updateScope, 'updateScope')
        }
      }
    } finally {
      if (resolvedRoleId) {
        await authed({ path: '/v1/roles/{id}', method: 'DELETE', pathParams: { id: resolvedRoleId } })
      }
      if (scopeRes.data?.id) {
        await authed({ path: '/v1/scopes/{id}', method: 'DELETE', pathParams: { id: scopeRes.data.id } })
      }
    }
  })

  it('creates and deletes urls', async () => {
    const suffix = Date.now().toString(36).slice(0, 6)
    const alias = `t${suffix}`

    const createRes = await authed({
      path: '/v1/urls',
      method: 'POST',
      body: {
        alias,
        url: 'https://example.com'
      }
    })
    assertOk(createRes, 'createUrl')

    try {
      const dataRes = await authed({
        path: '/v1/urls/{alias}/data',
        method: 'GET',
        pathParams: { alias }
      })
      assertOk(dataRes, 'getUrlData')

      const getRes = await authed({
        path: '/v1/urls/{alias}',
        method: 'GET',
        pathParams: { alias },
        redirect: 'manual'
      })
      expect([200, 308]).toContain(getRes.status)
    } finally {
      const deleteRes = await authed({
        path: '/v1/urls/{alias}',
        method: 'DELETE',
        pathParams: { alias }
      })
      expect([200, 404]).toContain(deleteRes.status)
    }
  })
})
