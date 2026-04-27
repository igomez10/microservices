import { expect, test, type Page } from '@playwright/test'

const username = process.env.E2E_USERNAME
const password = process.env.E2E_PASSWORD

async function signIn(page: Page) {
  await page.goto('/')
  await page.getByLabel('Username').fill(username)
  await page.getByLabel('Password').fill(password)
  await page.getByRole('button', { name: 'Sign in' }).click()
  await expect(page.getByRole('button', { name: 'Users' })).toBeVisible()
}

test.describe.configure({ mode: 'serial' })

test.skip(!username || !password, 'E2E_USERNAME and E2E_PASSWORD are required for frontend auth tests')

test('signs in with basic auth and stores the token cookie', async ({ page, context }) => {
  await signIn(page)

  const cookies = await context.cookies()
  const tokenCookie = cookies.find((cookie) => cookie.name === 'socialapp.token')
  expect(tokenCookie?.value).toBeTruthy()
})

test('shows paginated users with page navigation 1 through 10', async ({ page }) => {
  await signIn(page)
  await page.getByRole('button', { name: 'Users' }).click()

  for (let pageNumber = 1; pageNumber <= 10; pageNumber += 1) {
    await expect(page.getByRole('button', { name: String(pageNumber), exact: true })).toBeVisible()
  }

  const pageTwoButton = page.getByRole('button', { name: '2', exact: true })
  await pageTwoButton.click()
  await expect(pageTwoButton).toHaveClass(/btn-primary/)
})

test('profile nav shows my own profile by default', async ({ page }) => {
  await signIn(page)
  await page.getByRole('button', { name: 'Profile' }).click()

  await expect(page.getByRole('heading', { name: 'Profile' })).toBeVisible()
  await expect(page.locator('.profile-handle')).toHaveText(`@${username}`)
})

test('clicking a user from users opens that user profile', async ({ page }) => {
  await signIn(page)
  await page.getByRole('button', { name: 'Users' }).click()
  await expect(page.locator('.user-row.selectable').first()).toBeVisible()

  const userRows = page.locator('.user-row.selectable')
  const count = await userRows.count()
  const firstHandle = (await userRows.nth(0).locator('.user-row-handle').textContent())?.trim() ?? ''
  const secondHandle = count > 1 ? (await userRows.nth(1).locator('.user-row-handle').textContent())?.trim() ?? '' : ''
  const targetHandle = firstHandle === `@${username}` && secondHandle ? secondHandle : firstHandle
  const targetRow = targetHandle === secondHandle ? userRows.nth(1) : userRows.nth(0)

  await targetRow.click()

  await expect(page.getByRole('heading', { name: 'Profile' })).toBeVisible()
  await expect(page.locator('.profile-handle')).toHaveText(targetHandle)
})

test('posting from feed clears the textarea and shows a posted comment notice', async ({ page }) => {
  await signIn(page)
  await page.getByRole('button', { name: 'Feed' }).click()

  const comment = `Playwright post ${Date.now()}`
  const textarea = page.getByPlaceholder("What's on your mind?")
  await textarea.fill(comment)
  await page.getByRole('button', { name: 'Post' }).click()

  await expect(textarea).toHaveValue('')
  await expect(page.getByText(`Posted comment by @${username}`)).toBeVisible()
  await expect(page.getByText(comment, { exact: false })).toBeVisible()
})

test('liking a post sends successful like and unlike requests', async ({ page }) => {
  await signIn(page)
  await page.getByRole('button', { name: 'Feed' }).click()

  const comment = `Playwright like test ${Date.now()}`
  const textarea = page.getByPlaceholder("What's on your mind?")
  await textarea.fill(comment)
  await page.getByRole('button', { name: 'Post' }).click()
  await expect(page.getByText(`Posted comment by @${username}`)).toBeVisible()

  await page.getByRole('button', { name: 'Profile' }).click()
  await page.getByRole('button', { name: 'Posts' }).click()

  const targetPost = page.locator('.comment-item').filter({ hasText: comment }).first()
  await expect(targetPost).toBeVisible()

  const likeButton = targetPost.locator('.comment-like-button')
  const likeResponsePromise = page.waitForResponse((response) => {
    return response.url().includes('/like') && response.request().method() === 'POST'
  })
  await likeButton.click()
  const likeResponse = await likeResponsePromise
  expect(likeResponse.status()).toBe(200)

  const unlikeResponsePromise = page.waitForResponse((response) => {
    return response.url().includes('/like') && response.request().method() === 'DELETE'
  })
  await likeButton.click()
  const unlikeResponse = await unlikeResponsePromise
  expect(unlikeResponse.status()).toBe(200)
})
