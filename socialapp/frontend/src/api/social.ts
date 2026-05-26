import { apiRequest, type ApiResponse } from './client'
import type { components } from './types'

export type AccessToken = components['schemas']['AccessToken']
export type User = components['schemas']['User']
export type CreateUserRequest = components['schemas']['CreateUserRequest']
export type CreateUserResponse = components['schemas']['CreateUserResponse']
export type CreateCommentRequest = components['schemas']['CreateCommentRequest']
export type Comment = components['schemas']['Comment']
export type ChangePasswordRequest = components['schemas']['ChangePasswordRequest']
export type ResetPasswordRequest = components['schemas']['ResetPasswordRequest']
export type Role = components['schemas']['Role']
export type Scope = components['schemas']['Scope']
export type URLRecord = components['schemas']['URL']
export type LikeRequest = components['schemas']['LikeRequest']

export type ApiContext = {
  baseUrl: string
  token?: string
  clientId?: string
  clientSecret?: string
}

export type SocialApi = ReturnType<typeof createSocialApi>

export function createSocialApi(ctx: ApiContext) {
  const authed = <T,>(options: Omit<Parameters<typeof apiRequest<T>>[0], 'baseUrl' | 'token'>) =>
    apiRequest<T>({ baseUrl: ctx.baseUrl, token: ctx.token, ...options })

  const basic = <T,>(options: Omit<Parameters<typeof apiRequest<T>>[0], 'baseUrl' | 'basicAuth'>) =>
    apiRequest<T>({
      baseUrl: ctx.baseUrl,
      basicAuth: ctx.clientId && ctx.clientSecret ? { clientId: ctx.clientId, clientSecret: ctx.clientSecret } : undefined,
      ...options
    })

  return {
    welcome: (): Promise<ApiResponse<string>> => authed({ path: '/', method: 'GET' }),
    getAccessToken: (scopes?: string): Promise<ApiResponse<AccessToken>> =>
      basic({
        path: '/v1/oauth/token',
        method: 'POST',
        queryParams: scopes ? { scope: scopes } : undefined
      }),
    listUsers: (params?: { limit?: number; offset?: number; search?: string }): Promise<ApiResponse<User[]>> =>
      authed({ path: '/v1/users', method: 'GET', queryParams: params }),
    createUser: (body: CreateUserRequest): Promise<ApiResponse<CreateUserResponse>> =>
      authed({ path: '/v1/users', method: 'POST', body }),
    getUserByUsername: (username: string): Promise<ApiResponse<User>> =>
      authed({ path: '/v1/users/{username}', method: 'GET', pathParams: { username } }),
    updateUser: (username: string, body: User): Promise<ApiResponse<User>> =>
      authed({ path: '/v1/users/{username}', method: 'PUT', pathParams: { username }, body }),
    deleteUser: (username: string): Promise<ApiResponse<User>> =>
      authed({ path: '/v1/users/{username}', method: 'DELETE', pathParams: { username } }),
    changePassword: (body: ChangePasswordRequest): Promise<ApiResponse<User>> =>
      authed({ path: '/v1/password', method: 'POST', body }),
    resetPassword: (body: ResetPasswordRequest): Promise<ApiResponse<User>> =>
      authed({ path: '/v1/password', method: 'PUT', body }),
    getUserFeed: (): Promise<ApiResponse<Comment[]>> => authed({ path: '/v1/feed', method: 'GET' }),
    getUserComments: (username: string): Promise<ApiResponse<Comment[]>> =>
      authed({ path: '/v1/users/{username}/comments', method: 'GET', pathParams: { username } }),
    getUserFollowers: (username: string): Promise<ApiResponse<User[]>> =>
      authed({ path: '/v1/users/{username}/followers', method: 'GET', pathParams: { username } }),
    getFollowingUsers: (username: string): Promise<ApiResponse<User[]>> =>
      authed({ path: '/v1/users/{username}/following', method: 'GET', pathParams: { username } }),
    followUser: (followedUsername: string, followerUsername: string): Promise<ApiResponse<unknown>> =>
      authed({
        path: '/v1/users/{followedUsername}/followers/{followerUsername}',
        method: 'POST',
        pathParams: { followedUsername, followerUsername }
      }),
    unfollowUser: (followedUsername: string, followerUsername: string): Promise<ApiResponse<unknown>> =>
      authed({
        path: '/v1/users/{followedUsername}/followers/{followerUsername}',
        method: 'DELETE',
        pathParams: { followedUsername, followerUsername }
      }),
    getRolesForUser: (username: string): Promise<ApiResponse<Role[]>> =>
      authed({ path: '/v1/users/{username}/roles', method: 'GET', pathParams: { username } }),
    updateRolesForUser: (username: string, roles: string[]): Promise<ApiResponse<Role[]>> =>
      authed({ path: '/v1/users/{username}/roles', method: 'PUT', pathParams: { username }, body: roles }),
    getComment: (id: number): Promise<ApiResponse<Comment>> =>
      authed({ path: '/v1/comments/{id}', method: 'GET', pathParams: { id } }),
    searchComments: (params?: {
      username?: string
      start_time?: string
      end_time?: string
    }): Promise<ApiResponse<Comment[]>> => authed({ path: '/v1/comments', method: 'GET', queryParams: params }),
    createComment: (body: CreateCommentRequest): Promise<ApiResponse<Comment>> =>
      authed({ path: '/v1/comments', method: 'POST', body }),
    likeComment: (body: LikeRequest): Promise<ApiResponse<Comment>> =>
      authed({ path: '/like', method: 'POST', body }),
    unlikeComment: (body: LikeRequest): Promise<ApiResponse<Comment>> =>
      authed({ path: '/like', method: 'DELETE', body }),
    listRoles: (params?: { limit?: number; offset?: number }): Promise<ApiResponse<Role[]>> =>
      authed({ path: '/v1/roles', method: 'GET', queryParams: params }),
    createRole: (body: Role): Promise<ApiResponse<Role>> => authed({ path: '/v1/roles', method: 'POST', body }),
    getRole: (id: number): Promise<ApiResponse<Role>> =>
      authed({ path: '/v1/roles/{id}', method: 'GET', pathParams: { id } }),
    updateRole: (id: number, body: Role): Promise<ApiResponse<Role>> =>
      authed({ path: '/v1/roles/{id}', method: 'PUT', pathParams: { id }, body }),
    deleteRole: (id: number): Promise<ApiResponse<unknown>> =>
      authed({ path: '/v1/roles/{id}', method: 'DELETE', pathParams: { id } }),
    listScopesForRole: (id: number): Promise<ApiResponse<Scope[]>> =>
      authed({ path: '/v1/roles/{id}/scopes', method: 'GET', pathParams: { id } }),
    addScopeToRole: (id: number, scopes: string[]): Promise<ApiResponse<unknown>> =>
      authed({ path: '/v1/roles/{id}/scopes', method: 'POST', pathParams: { id }, body: scopes }),
    removeScopeFromRole: (roleId: number, scopeId: number): Promise<ApiResponse<unknown>> =>
      authed({
        path: '/v1/roles/{role_id}/scopes/{scope_id}',
        method: 'DELETE',
        pathParams: { role_id: roleId, scope_id: scopeId }
      }),
    listScopes: (params?: { limit?: number; offset?: number }): Promise<ApiResponse<Scope[]>> =>
      authed({ path: '/v1/scopes', method: 'GET', queryParams: params }),
    createScope: (body: Scope): Promise<ApiResponse<Scope>> =>
      authed({ path: '/v1/scopes', method: 'POST', body }),
    getScope: (id: number): Promise<ApiResponse<Scope>> =>
      authed({ path: '/v1/scopes/{id}', method: 'GET', pathParams: { id } }),
    updateScope: (id: number, body: Scope): Promise<ApiResponse<Scope>> =>
      authed({ path: '/v1/scopes/{id}', method: 'PUT', pathParams: { id }, body }),
    deleteScope: (id: number): Promise<ApiResponse<unknown>> =>
      authed({ path: '/v1/scopes/{id}', method: 'DELETE', pathParams: { id } }),
    getUrl: (alias: string): Promise<ApiResponse<unknown>> =>
      authed({ path: '/v1/urls/{alias}', method: 'GET', pathParams: { alias }, redirect: 'manual' }),
    getUrlData: (alias: string): Promise<ApiResponse<URLRecord>> =>
      authed({ path: '/v1/urls/{alias}/data', method: 'GET', pathParams: { alias } }),
    createUrl: (body: URLRecord): Promise<ApiResponse<URLRecord>> =>
      authed({ path: '/v1/urls', method: 'POST', body }),
    deleteUrl: (alias: string): Promise<ApiResponse<unknown>> =>
      authed({ path: '/v1/urls/{alias}', method: 'DELETE', pathParams: { alias } })
  }
}
