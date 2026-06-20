import instance from '@/utils/axios'

export interface MessageTemplate {
  id: string
  name: string
  title: string
  content: string
  active: boolean
  recipient_type: number
  variables: string[]
  created_at: number
  updated_at: number
}

export interface UpdateMessageTemplateRequest {
  title: string
  content: string
}

export const GetMessageTemplates = () => {
  return instance({
    url: '/message-templates',
    method: 'get'
  })
}

export const UpdateMessageTemplate = (id: string, request: UpdateMessageTemplateRequest) => {
  return instance({
    url: `/message-templates/${id}`,
    method: 'put',
    data: request
  })
}

export const UpdateMessageTemplateActive = (id: string, active: boolean) => {
  return instance({
    url: `/message-templates/${id}/active`,
    method: 'put',
    data: { active }
  })
}
