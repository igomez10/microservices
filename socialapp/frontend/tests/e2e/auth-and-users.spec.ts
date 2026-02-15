import { expect, test } from '@playwright/test'

const clientId = process.env.E2E_CLIENT_ID
const clientSecret = process.env.E2E_CLIENT_SECRET
const baseUrl = process.env.E2E_BASE_URL ?? 'http://localhost:4173'
const apiBaseUrl = process.env.E2E_API_BASE_URL ?? `${baseUrl}/api`
const scopes =
  process.env.E2E_SCOPES ??
  'socialapp.users.list socialapp.users.create socialapp.users.read socialapp.users.update socialapp.users.delete socialapp.users.roles.list socialapp.users.roles.update socialapp.comments.create socialapp.comments.read socialapp.feed.read socialapp.follower.create socialapp.follower.read socialapp.follower.delete socialapp.following.list socialapp.roles.list socialapp.roles.create socialapp.roles.read socialapp.roles.update socialapp.roles.delete socialapp.scopes.list socialapp.scopes.create socialapp.scopes.read socialapp.scopes.update socialapp.scopes.delete socialapp.roles.scopes.list socialapp.roles.scopes.create socialapp.roles.scopes.delete shortly.url.create shortly.url.delete'

const shouldRun = clientId && clientSecret

test.describe.configure({ mode: 'serial' })

test.skip(!shouldRun, 'E2E_CLIENT_ID and E2E_CLIENT_SECRET are required for UI tests')

test('get access token', async ({ page }) => {
  await page.goto('/')
  await page.getByLabel('API Base URL').fill(apiBaseUrl)
  await page.getByRole('button', { name: 'Sign in' }).click()

  await page.getByRole('button', { name: 'Settings' }).click()
  await page.getByLabel('Client ID').fill(clientId!)
  await page.getByLabel('Client Secret').fill(clientSecret!)
  await page.getByLabel('Scopes (space-separated)').fill(scopes)
  await page.getByRole('button', { name: 'Get access token' }).click()

  await expect(page.getByText('Token set')).toBeVisible()
})

test('list users and create/delete a user', async ({ page }) => {
  await page.goto('/')
  await page.getByLabel('API Base URL').fill(apiBaseUrl)
  await page.getByRole('button', { name: 'Sign in' }).click()

  await page.getByRole('button', { name: 'Settings' }).click()
  await page.getByLabel('Client ID').fill(clientId!)
  await page.getByLabel('Client Secret').fill(clientSecret!)
  await page.getByLabel('Scopes (space-separated)').fill(scopes)
  await page.getByRole('button', { name: 'Get access token' }).click()
  await expect(page.getByText('Token set')).toBeVisible()

  await page.getByRole('button', { name: 'Users' }).click()
  await page.getByRole('button', { name: 'Refresh' }).click()

  const suffix = Math.random().toString(36).slice(2, 8)
  const username = `ui_${suffix}`

  const registerCard = page.getByRole('heading', { name: 'Register' }).locator('..').locator('..')
  await registerCard.getByLabel('Username').fill(username)
  await registerCard.getByLabel('First name').fill('UI')
  await registerCard.getByLabel('Last name').fill('Test')
  await registerCard.getByLabel('Email').fill(`${username}@mail.com`)
  await registerCard.getByLabel('Password').fill('Password123!')
  await registerCard.getByRole('button', { name: 'Create user' }).click()

  const profileCard = page.getByRole('heading', { name: 'Profile' }).locator('..').locator('..')
  await expect(profileCard).toBeVisible()
  await profileCard.getByRole('button', { name: 'Delete' }).click()
})
