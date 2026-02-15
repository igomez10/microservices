// @vitest-environment node
import fs from 'node:fs'
import path from 'node:path'
import yaml from 'yaml'
import { operations } from '@/api/operations'

type OpenApiSpec = {
  paths?: Record<string, Record<string, { operationId?: string }>>
}

describe('openapi coverage', () => {
  it('tracks every operation from openapi.yaml', () => {
    const specPath = path.resolve(__dirname, '../../../openapi.yaml')
    const raw = fs.readFileSync(specPath, 'utf8')
    const spec = yaml.parse(raw) as OpenApiSpec
    const opIds: string[] = []

    Object.values(spec.paths ?? {}).forEach((methods) => {
      Object.entries(methods).forEach(([method, operation]) => {
        if (['get', 'post', 'put', 'delete', 'patch', 'head', 'options'].includes(method)) {
          if (operation.operationId) {
            opIds.push(operation.operationId)
          }
        }
      })
    })

    const missing = opIds.filter((id) => !operations.includes(id as (typeof operations)[number]))
    const extra = operations.filter((id) => !opIds.includes(id))

    expect(missing, `Missing operations: ${missing.join(', ')}`).toHaveLength(0)
    expect(extra, `Extra operations: ${extra.join(', ')}`).toHaveLength(0)
  })
})
